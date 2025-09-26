package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client interface {
	Connect(host string, port int, password string) error
	Close() error
	HealthCheck(ctx context.Context) error

	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	GetClient() *redis.Client
	GetExpiration(key string) (time.Duration, error)
	SetJSON(key string, value interface{}, expiration time.Duration) error
	GetJSON(key string, v interface{}) error
	GetInt(key string) (int, error)
	GetInt64(key string) (int64, error)
	GetUint64(key string) (uint64, error)
	GetFloat32(key string) (float32, error)
	GetFloat64(key string) (float64, error)
	GetBool(key string) (bool, error)
	GetBytes(key string) ([]byte, error)
	GetTime(key string) (time.Time, error)

	Del(key string) error
	Expire(key string, expiration time.Duration) error
	Exists(key string) (bool, error)
}

type client struct {
	rdb *redis.Client
	ctx context.Context
}

func New() Client {
	return &client{
		ctx: context.Background(),
	}
}

func (c *client) WithContext(ctx context.Context) Client {
	return &client{
		rdb: c.rdb,
		ctx: ctx,
	}
}

func (c *client) Connect(host string, port int, password string) error {
	c.rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprint(host, ":", port),
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	fmt.Println("Pinging redis server...")
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return err
	}

	fmt.Println("Redis server is up and running...")
	return nil
}

func (c *client) Close() error {
	return c.rdb.Close()
}

func (c *client) HealthCheck(ctx context.Context) error {
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return err
	}

	return nil
}

func (c *client) Set(key string, value interface{}, expiration time.Duration) error {
	return c.rdb.Set(c.ctx, key, value, expiration).Err()
}

func (c *client) Get(key string) (string, error) {
	val, err := c.rdb.Get(c.ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (c *client) GetClient() *redis.Client {
	return c.rdb
}

func (c *client) GetExpiration(key string) (time.Duration, error) {
	duration, err := c.rdb.TTL(c.ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return duration, nil
}

func (c *client) SetJSON(key string, value interface{}, expiration time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.rdb.Set(c.ctx, key, v, expiration).Err()
}

func (c *client) GetJSON(key string, v interface{}) error {
	val, err := c.rdb.Get(c.ctx, key).Bytes()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(val, v); err != nil {
		return err
	}

	return nil
}

func (c *client) GetInt(key string) (int, error) {
	val, err := c.rdb.Get(c.ctx, key).Int()
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (c *client) GetInt64(key string) (int64, error) {
	val, err := c.rdb.Get(c.ctx, key).Int64()
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (c *client) GetBool(key string) (bool, error) {
	val, err := c.rdb.Get(c.ctx, key).Bool()
	if err != nil {
		return false, err
	}

	return val, nil
}

func (c *client) GetFloat32(key string) (float32, error) {
	val, err := c.rdb.Get(c.ctx, key).Float32()
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (c *client) GetBytes(key string) ([]byte, error) {
	val, err := c.rdb.Get(c.ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (c *client) GetTime(key string) (time.Time, error) {
	val, err := c.rdb.Get(c.ctx, key).Time()
	if err != nil {
		return time.Time{}, err
	}

	return val, nil
}

func (c *client) GetUint64(key string) (uint64, error) {
	val, err := c.rdb.Get(c.ctx, key).Uint64()
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (c *client) GetFloat64(key string) (float64, error) {
	val, err := c.rdb.Get(c.ctx, key).Float64()
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (c *client) Del(key string) error {
	return c.rdb.Del(c.ctx, key).Err()
}

func (c *client) Expire(key string, expiration time.Duration) error {
	return c.rdb.Expire(c.ctx, key, expiration).Err()
}

func (c *client) Exists(key string) (bool, error) {
	v, err := c.rdb.Exists(c.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return v == 1, nil
}
