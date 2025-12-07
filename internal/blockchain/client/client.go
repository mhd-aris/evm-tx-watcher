package client

import (
	"context"
	"evm-tx-watcher/internal/config"
	"evm-tx-watcher/internal/domain"
	"evm-tx-watcher/internal/util"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

// Client represents a blockchain client for a specific network
type Client struct {
	NetworkConfig config.NetworkConfig
	ethClient     *ethclient.Client
	logger        *util.Logger
}

// ERC20TransferEvent represents the Transfer event signature
var ERC20TransferEventSignature = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

// TransactionDetails contains both transaction and its token transfers
type TransactionDetails struct {
	Transaction    *domain.Transaction
	TokenTransfers []domain.TokenTransfer
}

// New creates a new blockchain client
func New(networkConfig config.NetworkConfig, logger *util.Logger) (*Client, error) {
	ethClient, err := ethclient.Dial(networkConfig.RPC)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", networkConfig.RPC, err)
	}

	client := &Client{
		NetworkConfig: networkConfig,
		ethClient:     ethClient,
		logger:        logger,
	}

	// Test connection and verify chain ID
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chainID, err := client.ethClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID for %s: %w", networkConfig.Name, err)
	}

	if chainID.Int64() != networkConfig.ChainID {
		return nil, fmt.Errorf("chain ID mismatch for %s: expected %d, got %d", 
			networkConfig.Name, networkConfig.ChainID, chainID.Int64())
	}

	logger.Infof("Connected to %s (Chain ID: %d)", networkConfig.Name, chainID.Int64())
	return client, nil
}

// Close closes the client connection
func (c *Client) Close() {
	c.ethClient.Close()
}

// GetLatestBlockNumber returns the latest block number
func (c *Client) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	return c.ethClient.BlockNumber(ctx)
}

// SubscribeNewHeads subscribes to new block headers
func (c *Client) SubscribeNewHeads(ctx context.Context, ch chan<- *types.Header) (chan error, error) {
	sub, err := c.ethClient.SubscribeNewHead(ctx, ch)
	if err != nil {
		return nil, err
	}

	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		select {
		case err := <-sub.Err():
			errCh <- err
		case <-ctx.Done():
			return
		}
	}()

	return errCh, nil
}

// GetBlockWithTransactions retrieves a block with all its transactions and receipts
func (c *Client) GetBlockWithTransactions(ctx context.Context, blockNumber *big.Int) (*types.Block, []*TransactionDetails, error) {
	block, err := c.ethClient.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get block %d: %w", blockNumber.Int64(), err)
	}

	var transactionDetails []*TransactionDetails

	for _, tx := range block.Transactions() {
		// Get transaction receipt for logs and status
		receipt, err := c.ethClient.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			c.logger.Warnf("[%s] Failed to get receipt for tx %s: %v", c.NetworkConfig.Name, tx.Hash().Hex(), err)
			continue
		}

		// Convert to domain transaction
		domainTx, err := c.convertToDomainTransaction(block, tx, receipt)
		if err != nil {
			c.logger.Warnf("[%s] Failed to convert transaction %s: %v", c.NetworkConfig.Name, tx.Hash().Hex(), err)
			continue
		}

		// Generate UUID for the transaction
		domainTx.ID = uuid.New()

		// Extract token transfers from logs
		tokenTransfers := c.extractTokenTransfers(domainTx.ID, receipt.Logs)

		transactionDetails = append(transactionDetails, &TransactionDetails{
			Transaction:    domainTx,
			TokenTransfers: tokenTransfers,
		})
	}

	return block, transactionDetails, nil
}

// convertToDomainTransaction converts go-ethereum types to domain transaction
func (c *Client) convertToDomainTransaction(block *types.Block, tx *types.Transaction, receipt *types.Receipt) (*domain.Transaction, error) {
	// Get the sender address using the correct signer for this chain
	signer := types.LatestSignerForChainID(big.NewInt(c.NetworkConfig.ChainID))
	from, err := types.Sender(signer, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender: %w", err)
	}

	domainTx := &domain.Transaction{
		Hash:             tx.Hash().Hex(),
		BlockNumber:      block.Number().Int64(),
		BlockHash:        block.Hash().Hex(),
		TransactionIndex: int(receipt.TransactionIndex),
		ChainID:          c.NetworkConfig.ChainID,
		FromAddress:      strings.ToLower(from.Hex()),
		Value:            tx.Value(),
		GasUsed:          new(int64),
		GasPrice:         tx.GasPrice(),
		TxType:           int(tx.Type()),
		Status:           int(receipt.Status),
		BlockTimestamp:   time.Unix(int64(block.Time()), 0),
		CreatedAt:        time.Now(),
	}

	*domainTx.GasUsed = int64(receipt.GasUsed)

	if tx.To() != nil {
		toAddr := strings.ToLower(tx.To().Hex())
		domainTx.ToAddress = &toAddr
	}

	return domainTx, nil
}

// extractTokenTransfers extracts ERC-20 transfer events from transaction logs
func (c *Client) extractTokenTransfers(transactionID uuid.UUID, logs []*types.Log) []domain.TokenTransfer {
	var transfers []domain.TokenTransfer

	for _, log := range logs {
		// Check if this is an ERC-20 Transfer event
		if len(log.Topics) != 3 || log.Topics[0] != ERC20TransferEventSignature {
			continue
		}

		// Extract from and to addresses from topics
		from := common.BytesToAddress(log.Topics[1].Bytes())
		to := common.BytesToAddress(log.Topics[2].Bytes())

		// Extract value from data (should be 32 bytes)
		if len(log.Data) != 32 {
			continue
		}
		value := new(big.Int).SetBytes(log.Data)

		transfer := domain.TokenTransfer{
			ID:            uuid.New(),
			TransactionID: transactionID,
			LogIndex:      int(log.Index),
			TokenAddress:  strings.ToLower(log.Address.Hex()),
			FromAddress:   strings.ToLower(from.Hex()),
			ToAddress:     strings.ToLower(to.Hex()),
			Value:         value,
			CreatedAt:     time.Now(),
		}

		transfers = append(transfers, transfer)
	}

	return transfers
}
