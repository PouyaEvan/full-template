package integration

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRedisIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// Start Redis Container
	redisContainer, err := redis.Run(ctx, "redis:7-alpine")
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	// Clean up the container
	defer func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Get the connection string
	endpoint, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})
	defer rdb.Close()

	// Test Set/Get
	key := "test_key"
	value := "hello_world"

	err = rdb.Set(ctx, key, value, 10*time.Second).Err()
	assert.NoError(t, err)

	val, err := rdb.Get(ctx, key).Result()
	assert.NoError(t, err)
	assert.Equal(t, value, val)
}
