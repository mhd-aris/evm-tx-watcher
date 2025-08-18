# EVM Transaction Watcher

A simple REST API service for monitoring Ethereum wallet addresses and getting notified when transactions occur.

## 🎯 What This Project Does

This service allows you to:
- Register Ethereum wallet addresses to monitor
- Set up webhook notifications for transactions
- Get real-time alerts when monitored addresses receive or send transactions

Perfect for DeFi applications, portfolio trackers, or any service that needs to react to blockchain events.

## 🚀 Getting Started

### Prerequisites
- Go 1.24.6 or higher
- PostgreSQL database
- Basic understanding of Ethereum addresses

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/mhd-aris/evm-tx-watcher.git
   cd evm-tx-watcher
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Run the application**
   ```bash
   make run
   # or
   go run cmd/api/main.go
   ```

## 📡 API Endpoints

### Health Check
```http
GET /health
```

### Register Address to Monitor
```http
POST /api/v1/addresses
Content-Type: application/json

{
  "address": "0x742d35Cc6634C0532925a3b8D0C9c23a5d04C02a",
  "chain_id": 1,
  "webhook_url": "https://your-app.com/webhook",
  "secret": "your-webhook-secret",
  "label": "My Wallet"
}
```

### Get All Monitored Addresses
```http
GET /api/v1/addresses
```

## 🛠 Tech Stack

- **Language**: Go
- **Web Framework**: Echo v4
- **Database**: PostgreSQL with SQLx
- **Validation**: go-playground/validator
- **Config**: Viper
- **Logging**: Logrus
- **Documentation**: Swagger

## 📚 Project Structure

```
├── cmd/api/           # Application entry point
├── internal/
│   ├── config/        # Configuration management
│   ├── domain/        # Business entities
│   ├── dto/           # Data transfer objects
│   ├── http/          # HTTP handlers & routing
│   ├── repository/    # Data access layer
│   ├── service/       # Business logic
│   └── validator/     # Custom validators
├── docs/              # API documentation
└── migrations/        # Database migrations
```

## 🌟 Features

- ✅ REST API for address registration
- ✅ Ethereum address validation
- ✅ Webhook URL validation
- ✅ Swagger API documentation
- ✅ Clean architecture pattern
- ✅ Structured logging
- ⏳ Real-time transaction monitoring (coming soon)
- ⏳ Multiple blockchain support (coming soon)

## 🧪 Development

```bash
# Run with hot reload
make run

# Build application
make build

# Run linter
make lint
```

## 📝 Environment Variables

```env
APP_PORT=8080
LOG_LEVEL=info
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=evm_watcher
```

## 🤝 Contributing

This is a portfolio project, but feedback and suggestions are welcome! Feel free to:
- Open issues for bugs or feature requests
- Submit pull requests for improvements
- Star the repository if you find it useful

## 📈 Roadmap

- [ ] Real-time blockchain monitoring
- [ ] Support for multiple EVM networks (Ethereum, Arbitrum, Base, etc.)
- [ ] Transaction filtering and conditions
- [ ] Rate limiting and authentication
- [ ] Metrics and monitoring dashboard
- [ ] Docker containerization

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **Live Demo**: Coming soon
- **API Documentation**: `/swagger` endpoint when running

---

*This project is part of my Web3 development portfolio, showcasing backend development skills for blockchain applications.*
