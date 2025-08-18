package client

import (
	"context"
	"evm-tx-watcher/internal/util"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	Network string
	URL     string
	Eth     *ethclient.Client
	rpc     *rpc.Client
}

// RawTransaction represents a raw transaction from RPC
type RawTransaction struct {
	Hash     string `json:"hash"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Type     string `json:"type"`
	Nonce    string `json:"nonce"`
}

// RawBlock represents a raw block from RPC with transactions
type RawBlock struct {
	Number           string           `json:"number"`
	Hash             string           `json:"hash"`
	Timestamp        string           `json:"timestamp"`
	Transactions     []RawTransaction `json:"transactions"`
	TransactionCount int              `json:"-"` // calculated field
}

func New(network, url string, logger *util.Logger) (*Client, error) {
	c, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rpc dial %s: %w", url, err)
	}

	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("raw rpc dial %s: %w", url, err)
	}

	return &Client{Network: network, URL: url, Eth: c, rpc: rpcClient}, nil
}

// GetRPC returns the raw RPC client for manual calls
func (c *Client) GetRPC() *rpc.Client {
	return c.rpc
}

// GetTransactionByHash gets a transaction by its hash
func (c *Client) GetTransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	return c.Eth.TransactionByHash(ctx, hash)
}

// GetBlockTransactionCount gets the number of transactions in a block
func (c *Client) GetBlockTransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	return c.Eth.TransactionCount(ctx, blockHash)
}

// GetTransactionInBlock gets a transaction by block hash and index
func (c *Client) GetTransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	return c.Eth.TransactionInBlock(ctx, blockHash, index)
}

// GetBlockWithTransactionsRaw gets block with transaction details using raw RPC calls
func (c *Client) GetBlockWithTransactionsRaw(ctx context.Context, blockNumber *big.Int) (*RawBlock, error) {
	var result RawBlock

	blockNumberHex := fmt.Sprintf("0x%x", blockNumber)
	err := c.rpc.CallContext(ctx, &result, "eth_getBlockByNumber", blockNumberHex, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get block with transactions: %w", err)
	}

	result.TransactionCount = len(result.Transactions)
	return &result, nil
}
