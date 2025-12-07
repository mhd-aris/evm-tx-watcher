# Technical Development Guide

## 🏗️ **Architecture Deep Dive**

### **System Design Principles**

#### **1. Clean Architecture Implementation**
```
┌─────────────────────────────────────────────────────────────┐
│                     Presentation Layer                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   REST API      │  │   CLI Worker    │  │   Swagger   │ │
│  │   (Echo)        │  │   (Cobra)       │  │   Docs      │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                    Application Layer                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │  Address Svc    │  │  Webhook Svc    │  │ Processor   │ │
│  │  (Business)     │  │  (Delivery)     │  │ (Filtering) │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                     Domain Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Address       │  │  Transaction    │  │  Webhook    │ │
│  │   (Entity)      │  │   (Entity)      │  │  (Entity)   │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                  Infrastructure Layer                       │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │  PostgreSQL     │  │     Redis       │  │ Blockchain  │ │
│  │  (Repository)   │  │    (Cache)      │  │  (Client)   │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

#### **2. Dependency Flow**
- **Inward Dependencies**: Infrastructure → Domain ← Application ← Presentation
- **No Circular Dependencies**: Clean separation with interfaces
- **Testable Design**: Easy mocking and unit testing

### **Key Design Patterns**

#### **Repository Pattern**
```go
// Domain Interface (internal/repository/)
type AddressRepository interface {
    Create(ctx context.Context, tx *sqlx.Tx, address *domain.Address) (domain.Address, error)
    FindByAddress(ctx context.Context, address string) (*domain.Address, error)
    GetWatchedAddresses(ctx context.Context) ([]*domain.WatchedAddress, error)
}

// Infrastructure Implementation
type addressRepository struct {
    db *sqlx.DB
}
```

#### **Unit of Work Pattern**
```go
// Transactional consistency across multiple repositories
func (s *addressService) Register(ctx context.Context, req *dto.RegisterAddressRequest) (*dto.AddressResponse, *errors.AppError) {
    return s.unitOfWork.WithTransaction(ctx, func(tx *sqlx.Tx) error {
        // Create address
        addr, err := s.addressRepo.Create(ctx, tx, newAddress)
        if err != nil {
            return err
        }
        
        // Create webhook
        _, err = s.webhookRepo.Create(ctx, tx, newWebhook)
        return err
    })
}
```

#### **Observer Pattern (Blockchain Events)**
```go
// Event-driven architecture for blockchain monitoring
type BlockWatcher struct {
    client        *Client
    confirmations int64
}

func (w *Watcher) Start(ctx context.Context, out chan<- *types.Block) error {
    // Subscribe to blockchain events
    headers := make(chan *types.Header)
    sub, err := w.client.SubscribeNewHeads(ctx, headers)
    
    // Process confirmed blocks
    for header := range headers {
        if isConfirmed(header) {
            block := w.getBlock(header)
            out <- block  // Notify observers
        }
    }
}
```

---

## 🔧 **Component Analysis**

### **1. Configuration Management (`internal/config/`)**

#### **Features**
- Environment variable binding with Viper
- Validation and default values
- Hardcoded network configurations with override capability
- Type-safe configuration structs

#### **Key Implementation**
```go
type Config struct {
    AppPort   string            `mapstructure:"APP_PORT"`
    LogLevel  string            `mapstructure:"LOG_LEVEL"`
    DB        DatabaseConfig    `mapstructure:",squash"`
    Redis     RedisConfig       `mapstructure:",squash"`
    Networks  map[string]NetworkConfig `mapstructure:"-"`
}

// Hardcoded testnet configurations
func getTestnetNetworks() map[string]NetworkConfig {
    return map[string]NetworkConfig{
        "ethereum-sepolia": {Name: "ethereum-sepolia", ChainID: 11155111},
        "base-sepolia":     {Name: "base-sepolia", ChainID: 84532},
        "arbitrum-sepolia": {Name: "arbitrum-sepolia", ChainID: 421614},
    }
}
```

### **2. Blockchain Integration (`internal/blockchain/`)**

#### **Client Architecture**
```go
type Client struct {
    NetworkConfig config.NetworkConfig
    ethClient     *ethclient.Client
    logger        *util.Logger
}

// Chain ID verification on connection
func New(networkConfig config.NetworkConfig, logger *util.Logger) (*Client, error) {
    ethClient, err := ethclient.Dial(networkConfig.RPC)
    if err != nil {
        return nil, fmt.Errorf("failed to dial %s: %w", networkConfig.RPC, err)
    }

    // Verify chain ID matches configuration
    chainID, err := ethClient.ChainID(ctx)
    if chainID.Int64() != networkConfig.ChainID {
        return nil, fmt.Errorf("chain ID mismatch: expected %d, got %d", 
            networkConfig.ChainID, chainID.Int64())
    }
    
    return &Client{NetworkConfig: networkConfig, ethClient: ethClient, logger: logger}, nil
}
```

#### **Block Watcher with Confirmations**
```go
type Watcher struct {
    client        *Client
    confirmations int64  // Default: 5 blocks
    networkConfig config.NetworkConfig
}

func (w *Watcher) Start(ctx context.Context, out chan<- *types.Block) error {
    // Track pending blocks with confirmation requirements
    pending := make(map[uint64]*types.Header)
    
    for header := range headers {
        pending[header.Number.Uint64()] = header
        currentHead := header.Number.Uint64()
        
        // Process confirmed blocks
        for blockNum, blockHeader := range pending {
            if currentHead >= blockNum+uint64(w.confirmations) {
                block := w.getConfirmedBlock(blockHeader)
                out <- block
                delete(pending, blockNum)
            }
        }
    }
}
```

### **3. Data Layer (`internal/repository/`)**

#### **Database Schema Design**
```sql
-- Core entities with proper relationships
CREATE TABLE addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    address VARCHAR(42) NOT NULL,
    chain_id INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT unique_address_per_chain UNIQUE (address, chain_id)
);

CREATE TABLE webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    address_id UUID NOT NULL,
    url VARCHAR(2048) NOT NULL,
    secret VARCHAR(255) NOT NULL,
    CONSTRAINT fk_webhooks_address_id FOREIGN KEY (address_id) REFERENCES addresses(id)
);

-- Transaction tracking with comprehensive data
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hash VARCHAR(66) NOT NULL UNIQUE,
    block_number BIGINT NOT NULL,
    chain_id BIGINT NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42),
    value NUMERIC(78, 0) NOT NULL DEFAULT 0,
    status INTEGER NOT NULL DEFAULT 1
);

-- ERC-20 token transfer tracking
CREATE TABLE token_transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL,
    token_address VARCHAR(42) NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42) NOT NULL,
    value NUMERIC(78, 0) NOT NULL,
    CONSTRAINT fk_token_transfers_transaction_id FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);
```

#### **Repository Implementation Pattern**
```go
type addressRepository struct {
    db *sqlx.DB
}

func (r *addressRepository) Create(ctx context.Context, tx *sqlx.Tx, address *domain.Address) (domain.Address, error) {
    query := `INSERT INTO addresses (id, address, chain_id, is_active, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`
    
    _, err := tx.ExecContext(ctx, query,
        address.ID, address.Address, address.ChainID, 
        address.IsActive, address.CreatedAt, address.UpdatedAt)
    
    if err != nil {
        return domain.Address{}, fmt.Errorf("failed to insert address: %w", err)
    }
    
    return *address, nil
}

// Optimized query for watched addresses
func (r *addressRepository) GetWatchedAddresses(ctx context.Context) ([]*domain.WatchedAddress, error) {
    query := `
        SELECT DISTINCT a.address, a.chain_id, a.is_active, w.id as webhook_id, w.url as webhook_url
        FROM addresses a
        JOIN webhooks w ON a.id = w.address_id
        WHERE a.is_active = true`
    
    var watchedAddresses []*domain.WatchedAddress
    err := r.db.SelectContext(ctx, &watchedAddresses, query)
    return watchedAddresses, err
}
```

### **4. Caching Layer (`internal/cache/`)**

#### **Redis Integration**
```go
type RedisClient struct {
    client *redis.Client
}

// Address caching for fast lookup
func (r *RedisClient) CacheWatchedAddresses(ctx context.Context, addresses []domain.WatchedAddress) error {
    data, err := json.Marshal(addresses)
    if err != nil {
        return fmt.Errorf("failed to marshal watched addresses: %w", err)
    }
    
    return r.client.Set(ctx, WatchedAddressesKey, data, 5*time.Minute).Err()
}

// Webhook delivery queue
func (r *RedisClient) QueueWebhookDelivery(ctx context.Context, delivery *domain.WebhookDelivery) error {
    data, err := json.Marshal(delivery)
    if err != nil {
        return fmt.Errorf("failed to marshal webhook delivery: %w", err)
    }
    
    return r.client.LPush(ctx, WebhookQueueKey, data).Err()
}

// Block processing idempotency
func (r *RedisClient) SetProcessedBlock(ctx context.Context, network string, blockNumber int64) error {
    key := fmt.Sprintf(ProcessedBlockKey, network, blockNumber)
    return r.client.Set(ctx, key, "1", 24*time.Hour).Err()
}
```

### **5. Processing Pipeline (`internal/processor/`)**

#### **Current Implementation (Basic)**
```go
type Processor struct {
    log *util.Logger
}

func (p *Processor) HandleBlock(ctx context.Context, blk *types.Block) error {
    p.log.Infof("[Processor] block=%d hash=%s txs=%d",
        blk.NumberU64(), blk.Hash().Hex(), len(blk.Transactions()))

    // Process transactions with logging
    for i, tx := range blk.Transactions() {
        if i >= 3 { // Limit logging to prevent spam
            p.log.Infof("[Processor] ... and %d more transactions", len(blk.Transactions())-3)
            break
        }
        
        p.log.Infof("[Processor] tx[%d]: hash=%s value=%s",
            i, tx.Hash().Hex(), tx.Value().String())
    }
    
    return nil
}
```

#### **Advanced Implementation (Designed but not active)**
```go
// Full processor with filtering and webhook queuing
type Processor struct {
    logger           *util.Logger
    redis            *cache.RedisClient
    unitOfWork       repository.UnitOfWork
    addressRepo      repository.AddressRepository
    transactionRepo  repository.TransactionRepository
    
    watchedAddresses map[string]map[string]*domain.WatchedAddress
}

func (p *Processor) HandleBlockEvent(ctx context.Context, blockEvent *watcher.BlockEvent) error {
    // 1. Update watched addresses cache
    // 2. Filter relevant transactions
    // 3. Save to database
    // 4. Queue webhook deliveries
}
```

---

## 🚀 **Development Workflow**

### **Build System (Makefile)**
```makefile
# Professional development commands
build:          # Build API and Worker binaries
clean:          # Clean build artifacts  
test:           # Run tests (when implemented)
run-api:        # Start API server
run-worker:     # Start worker process
setup-dev:      # Setup development environment
migrate-up:     # Run database migrations up
migrate-down:   # Run database migrations down
fmt:            # Format Go code
lint:           # Run linter
docs:           # Generate swagger documentation
```

### **Environment Configuration**
```bash
# Application settings
APP_PORT=8080
LOG_LEVEL=info
LOG_FORMAT=text

# Database (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=evm_tx_watcher

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0

# Blockchain RPC URLs
RPC_ETHEREUM_SEPOLIA=https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
RPC_BASE_SEPOLIA=https://base-sepolia.g.alchemy.com/v2/YOUR_API_KEY
RPC_ARBITRUM_SEPOLIA=https://arb-sepolia.g.alchemy.com/v2/YOUR_API_KEY
```

### **Development Commands**
```bash
# Setup development environment
make setup-dev

# Build and test
make build
make test

# Run services
make run-api     # Terminal 1: API server
make run-worker  # Terminal 2: Blockchain worker

# Database operations
make migrate-up
make migrate-down

# Code quality
make fmt
make lint
```

---

## 🔍 **Code Quality Standards**

### **Error Handling Pattern**
```go
// Custom error types with context
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

// Consistent error wrapping
func (s *addressService) Register(ctx context.Context, req *dto.RegisterAddressRequest) (*dto.AddressResponse, *errors.AppError) {
    existingAddress, err := s.addressRepo.FindByAddress(ctx, req.Address)
    if err != nil {
        return nil, errors.Wrap(errors.ErrCodeDatabase, "failed to check existing address", err)
    }
    
    if existingAddress != nil {
        return nil, errors.AlreadyExists("Address already exists")
    }
    
    // Continue processing...
}
```

### **Logging Standards**
```go
// Structured logging with context
logger.WithFields(logrus.Fields{
    "network":     networkConfig.Name,
    "chain_id":    networkConfig.ChainID,
    "block_number": blockNumber,
}).Info("Processing confirmed block")

// Error logging with stack trace
logger.WithError(err).Errorf("Failed to process block %d for %s", blockNumber, network)
```

### **Testing Strategy (To Be Implemented)**
```go
// Unit test example
func TestAddressService_Register(t *testing.T) {
    // Setup mocks
    mockRepo := &mocks.AddressRepository{}
    mockUnitOfWork := &mocks.UnitOfWork{}
    
    service := NewAddressService(mockUnitOfWork, mockRepo, nil)
    
    // Test cases
    t.Run("successful registration", func(t *testing.T) {
        // Test implementation
    })
    
    t.Run("duplicate address", func(t *testing.T) {
        // Test implementation
    })
}

// Integration test example
func TestWorker_BlockProcessing(t *testing.T) {
    // Setup test environment
    testConfig := &config.Config{
        Networks: getTestNetworks(),
    }
    
    // Test blockchain integration
    worker := app.NewWorker(testConfig, logger)
    
    // Verify block processing
}
```

---

## 📊 **Performance Considerations**

### **Database Optimization**
```sql
-- Indexes for fast queries
CREATE INDEX idx_addresses_address ON addresses(address);
CREATE INDEX idx_addresses_chain_id ON addresses(chain_id);
CREATE INDEX idx_transactions_hash ON transactions(hash);
CREATE INDEX idx_transactions_block_number ON transactions(block_number);
CREATE INDEX idx_transactions_from_address ON transactions(from_address);
CREATE INDEX idx_transactions_to_address ON transactions(to_address);

-- Composite indexes for complex queries
CREATE INDEX idx_transactions_chain_block ON transactions(chain_id, block_number);
```

### **Caching Strategy**
```go
// Multi-level caching
// 1. In-memory cache for hot data
type Processor struct {
    watchedAddresses map[string]map[string]*domain.WatchedAddress // [chainID][address]
    lastCacheUpdate  time.Time
}

// 2. Redis cache for shared data
func (p *Processor) updateWatchedAddressesCache(ctx context.Context) error {
    // Update cache every 5 minutes
    if time.Since(p.lastCacheUpdate) < 5*time.Minute {
        return nil
    }
    
    // Try Redis first, fallback to database
    addresses, err := p.redis.GetWatchedAddresses(ctx)
    if len(addresses) == 0 {
        addresses, err = p.addressRepo.GetWatchedAddresses(ctx)
        p.redis.CacheWatchedAddresses(ctx, addresses)
    }
    
    // Update in-memory structure for O(1) lookup
    p.buildAddressLookupMap(addresses)
}
```

### **Concurrency Patterns**
```go
// Worker coordination with sync.WaitGroup
func RunWorker(ctx context.Context, cfg *config.Config, logger *util.Logger) error {
    var wg sync.WaitGroup
    blockChan := make(chan *types.Block, 50)
    
    // Start watchers for each network
    for _, networkConfig := range cfg.Networks {
        wg.Add(1)
        go func(network config.NetworkConfig) {
            defer wg.Done()
            watcher.Start(ctx, blockChan)
        }(networkConfig)
    }
    
    // Graceful shutdown
    go func() {
        wg.Wait()
        close(blockChan)
    }()
}
```

---

## 🎯 **Next Development Steps**

### **Immediate (Database Integration)**
1. **PostgreSQL Setup**: Connection and migration execution
2. **Repository Testing**: Verify CRUD operations
3. **API Integration**: Connect endpoints to database

### **Short Term (Transaction Processing)**
1. **Address Filtering**: Implement watched address filtering
2. **ERC-20 Detection**: Parse Transfer events from logs
3. **Database Persistence**: Save filtered transactions

### **Medium Term (Webhook System)**
1. **HTTP Client**: Implement webhook delivery
2. **HMAC Signatures**: Security implementation
3. **Retry Logic**: Exponential backoff with Redis queue

### **Long Term (Production Features)**
1. **Monitoring**: Metrics and health checks
2. **Testing**: Comprehensive test suite
3. **Documentation**: Code comments and guides
4. **Performance**: Optimization and benchmarking

This technical guide provides the foundation for understanding and extending the EVM Transaction Watcher codebase.
