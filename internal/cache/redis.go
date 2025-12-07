package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"evm-tx-watcher/internal/config"
	"evm-tx-watcher/internal/domain"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps redis client with application-specific methods
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg *config.RedisConfig) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{client: rdb}, nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Cache keys
const (
	WatchedAddressesKey = "watched_addresses"
	WebhookQueueKey     = "webhook_queue"
	ProcessedBlockKey   = "processed_block:%s:%d" // network:block_number
)

// CacheWatchedAddresses caches the list of watched addresses
func (r *RedisClient) CacheWatchedAddresses(ctx context.Context, addresses []domain.WatchedAddress) error {
	data, err := json.Marshal(addresses)
	if err != nil {
		return fmt.Errorf("failed to marshal watched addresses: %w", err)
	}

	return r.client.Set(ctx, WatchedAddressesKey, data, 5*time.Minute).Err()
}

// GetWatchedAddresses retrieves cached watched addresses
func (r *RedisClient) GetWatchedAddresses(ctx context.Context) ([]domain.WatchedAddress, error) {
	data, err := r.client.Get(ctx, WatchedAddressesKey).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get watched addresses from cache: %w", err)
	}

	var addresses []domain.WatchedAddress
	if err := json.Unmarshal([]byte(data), &addresses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal watched addresses: %w", err)
	}

	return addresses, nil
}

// InvalidateWatchedAddresses removes watched addresses from cache
func (r *RedisClient) InvalidateWatchedAddresses(ctx context.Context) error {
	return r.client.Del(ctx, WatchedAddressesKey).Err()
}

// QueueWebhookDelivery adds a webhook delivery to the queue
func (r *RedisClient) QueueWebhookDelivery(ctx context.Context, delivery *domain.WebhookDelivery) error {
	data, err := json.Marshal(delivery)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook delivery: %w", err)
	}

	return r.client.LPush(ctx, WebhookQueueKey, data).Err()
}

// DequeueWebhookDelivery gets a webhook delivery from the queue (blocking)
func (r *RedisClient) DequeueWebhookDelivery(ctx context.Context, timeout time.Duration) (*domain.WebhookDelivery, error) {
	result, err := r.client.BRPop(ctx, timeout, WebhookQueueKey).Result()
	if err == redis.Nil {
		return nil, nil // Timeout
	}
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue webhook delivery: %w", err)
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("invalid queue result format")
	}

	var delivery domain.WebhookDelivery
	if err := json.Unmarshal([]byte(result[1]), &delivery); err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhook delivery: %w", err)
	}

	return &delivery, nil
}

// SetProcessedBlock marks a block as processed for a network
func (r *RedisClient) SetProcessedBlock(ctx context.Context, network string, blockNumber int64) error {
	key := fmt.Sprintf(ProcessedBlockKey, network, blockNumber)
	return r.client.Set(ctx, key, "1", 24*time.Hour).Err()
}

// IsBlockProcessed checks if a block has been processed for a network
func (r *RedisClient) IsBlockProcessed(ctx context.Context, network string, blockNumber int64) (bool, error) {
	key := fmt.Sprintf(ProcessedBlockKey, network, blockNumber)
	_, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check processed block: %w", err)
	}
	return true, nil
}

// GetQueueLength returns the length of the webhook queue
func (r *RedisClient) GetQueueLength(ctx context.Context) (int64, error) {
	return r.client.LLen(ctx, WebhookQueueKey).Result()
}

