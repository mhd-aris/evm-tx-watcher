# EVM Transaction Watcher

![Go](https://img.shields.io/badge/Go-1.24+-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)
![Redis](https://img.shields.io/badge/Redis-7+-red.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)

A production-ready Go service for monitoring EVM-compatible blockchain addresses and sending webhook notifications when transactions occur. Supports multiple chains with ETH transfers and ERC-20 token tracking.

## 🚀 Features

- **Multi-Chain Support**: Ethereum, Base, and Arbitrum Sepolia testnets
- **Real-time Monitoring**: WebSocket subscriptions with configurable confirmations  
- **Comprehensive Tracking**: ETH transfers and ERC-20 token transfers
- **Webhook Notifications**: HMAC-signed webhooks with exponential backoff retry
- **High Performance**: Redis caching and PostgreSQL with optimized queries
- **Production Ready**: Clean architecture, comprehensive logging, error handling

## 🏗️ Architecture

The system consists of three main components:

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│   API       │    │   Worker     │    │  Webhook    │
│   Server    │    │   Process    │    │  Notifier   │
│   :8080     │    │              │    │             │
└─────────────┘    └──────────────┘    └─────────────┘
       │                   │                   │
       │            ┌──────▼──────┐           │
       │            │ Block       │           │
       │            │ Watchers    │           │
       │            │ (3 chains)  │           │
       │            └─────────────┘           │
       │                   │                   │
       ▼                   ▼                   ▼
┌─────────────────────────────────────────────────────┐
│                PostgreSQL Database                  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐   │
│  │ addresses   │ │transactions │ │webhook_     │   │
│  │ webhooks    │ │token_       │ │deliveries   │   │
│  │             │ │transfers    │ │             │   │
│  └─────────────┘ └─────────────┘ └─────────────┘   │
└─────────────────────────────────────────────────────┘
                          │
                          ▼
                 ┌─────────────┐
                 │    Redis    │
                 │   Cache &   │
                 │    Queue    │
                 └─────────────┘
```

## 🛠️ Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL 15+
- Redis 7+
- RPC endpoints for target networks

### Installation

1. **Clone and setup**
   ```bash
   git clone <repository>
   cd evm-tx-watcher
   make setup-dev
   ```

2. **Configure environment**
   ```bash
   cp env.example .env
   # Edit .env with your database credentials and RPC URLs
   ```

3. **Setup database**
   ```bash
   # Create database
   createdb evm_tx_watcher
   
   # Run migrations
   make migrate-up
   ```

4. **Build and run**
   ```bash
   # Build binaries
   make build
   
   # Start API server (terminal 1)
   make run-api
   
   # Start worker (terminal 2)
   make run-worker
   ```

### Configuration

The application uses environment variables for configuration. See `env.example` for a complete list:

**Required settings:**
```bash
# Database
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=evm_tx_watcher

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# RPC URLs (get from Alchemy, Infura, etc.)
RPC_ETHEREUM_SEPOLIA=https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
RPC_BASE_SEPOLIA=https://base-sepolia.g.alchemy.com/v2/YOUR_API_KEY  
RPC_ARBITRUM_SEPOLIA=https://arb-sepolia.g.alchemy.com/v2/YOUR_API_KEY
```

## 📋 API Usage

### Register Address for Monitoring

```bash
curl -X POST http://localhost:8080/api/v1/addresses \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x742d35Cc622C21F87d25Cf1177BB1B28F7E30aDF",
    "chain_id": 11155111,
    "webhook_url": "https://webhook.site/your-uuid",
    "secret": "your-webhook-secret",
    "label": "My Test Wallet"
  }'
```

### Get Registered Addresses

```bash
curl http://localhost:8080/api/v1/addresses
```

### Webhook Payload

Your webhook will receive transaction notifications with this structure:

```json
{
  "transaction_hash": "0x...",
  "block_number": 12345,
  "chain_id": 11155111,
  "from": "0x...",
  "to": "0x...",
  "value": "1000000000000000000",
  "gas_used": 21000,
  "gas_price": "20000000000",
  "status": 1,
  "timestamp": "2024-01-01T00:00:00Z",
  "token_transfers": [
    {
      "token_address": "0x...",
      "from": "0x...",
      "to": "0x...",
      "value": "1000000000000000000",
      "token_symbol": "USDC",
      "token_decimals": 6
    }
  ]
}
```

## 🔧 Development

### Available Commands

```bash
# Build binaries
make build

# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# Generate docs
make docs

# Clean build artifacts
make clean

# Setup development environment
make setup-dev

# Database migrations
make migrate-up
make migrate-down

# Quick development setup
make dev
```

### Project Structure

```
├── cmd/
│   ├── api/          # API server entry point
│   └── worker/       # Worker process entry point
├── internal/
│   ├── app/          # Application initialization
│   ├── blockchain/   # Blockchain clients and watchers
│   ├── cache/        # Redis client and operations
│   ├── config/       # Configuration management
│   ├── db/           # Database connection
│   ├── domain/       # Domain models
│   ├── dto/          # Data transfer objects
│   ├── http/         # HTTP handlers and routing
│   ├── processor/    # Transaction processing logic
│   ├── repository/   # Data access layer
│   ├── service/      # Business logic
│   ├── util/         # Utilities and helpers
│   ├── validator/    # Input validation
│   └── webhook/      # Webhook notification system
├── migrations/       # Database migrations
└── docs/            # API documentation
```

## 🌐 Supported Networks

| Network | Chain ID | Status |
|---------|----------|--------|
| Ethereum Sepolia | 11155111 | ✅ Active |
| Base Sepolia | 84532 | ✅ Active |
| Arbitrum Sepolia | 421614 | ✅ Active |

## 📊 Performance

- **Throughput**: 1000+ addresses across multiple chains
- **Latency**: ~30 seconds for confirmed transactions (5 block confirmations)
- **Reliability**: 99%+ webhook delivery with exponential backoff retry
- **Scalability**: Horizontal scaling ready with Redis/PostgreSQL

## 🚀 Deployment

### Docker (Coming Soon)
```bash
docker-compose up -d
```

### Manual Deployment
1. Build for production: `make prod-build`
2. Deploy binaries with your preferred method
3. Setup PostgreSQL and Redis
4. Configure environment variables
5. Run migrations: `make migrate-up`
6. Start services: `./bin/api` and `./bin/worker`

## 🧪 Testing

Test the system with your own addresses:

1. Get testnet ETH from faucets:
   - [Ethereum Sepolia Faucet](https://sepoliafaucet.com/)
   - [Base Sepolia Faucet](https://bridge.base.org/deposit)
   - [Arbitrum Sepolia Faucet](https://bridge.arbitrum.io/)

2. Register your address via API
3. Send test transactions
4. Check webhook deliveries

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙋‍♂️ Support

- Check the [Development Roadmap](DEVELOPMENT_ROADMAP.md) for future plans
- Open an issue for bugs or feature requests
- Contribute to make it better!

## 🔗 RPC Providers

Get free API keys from:
- [Alchemy](https://www.alchemy.com/) (Recommended)
- [Infura](https://infura.io/)
- [QuickNode](https://www.quicknode.com/)
- [Ankr](https://www.ankr.com/)

Free public endpoints are also available but with rate limits.

---

**Built with ❤️ for the Ethereum community**

*This project demonstrates production-ready Web3 backend development with Go, showcasing clean architecture, multi-chain support, and reliable webhook delivery systems.*

