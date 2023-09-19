package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}
type RedisParameter struct {
	Network          string
	Address          string
	Password         string
	DB               int
	DialTimeout      time.Duration
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	PoolSize         int
	SubscribeChannel string
	PublishChannel   string
}

func NewRedis(param *RedisParameter) (*Redis, error) {

	client := redis.NewClient(&redis.Options{
		Network:      param.Network,
		Addr:         param.Address,
		Password:     param.Password,
		DB:           param.DB,
		DialTimeout:  param.DialTimeout,
		ReadTimeout:  param.ReadTimeout,
		WriteTimeout: param.WriteTimeout,
		PoolSize:     param.PoolSize,
	})

	_, err := client.Ping(context.TODO()).Result()

	return &Redis{client: client}, err

}

func (rds *Redis) Set(key string, value interface{}) error {
	err := rds.client.Set(context.TODO(), key, value, 0).Err()

	if err != nil {
		return err
	}

	return nil
}

func (rds *Redis) GetSting(key string, defaultValue string) (string, error) {
	value, err := rds.client.Get(context.TODO(), key).Result()

	if err != nil {
		return defaultValue, err
	}

	return value, nil
}

func (rds *Redis) GetInt(key string, defaultValue int) (int, error) {
	value, err := rds.client.Get(context.TODO(), key).Result()

	if err != nil {
		return defaultValue, err
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue, err
	}

	return intValue, nil
}

func (rds *Redis) GetInt64(key string, defaultValue int64) (int64, error) {
	value, err := rds.client.Get(context.TODO(), key).Result()

	if err != nil {
		return defaultValue, err
	}

	int64Value, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue, err
	}

	return int64Value, nil
}

func (rds *Redis) GetFloat64(key string, defaultValue float64) (float64, error) {
	value, err := rds.client.Get(context.TODO(), key).Result()

	if err != nil {
		return defaultValue, err
	}

	float64Value, err := strconv.ParseFloat(value, 10)
	if err != nil {
		return defaultValue, err
	}

	return float64Value, nil
}

func (rds *Redis) HGetSting(key, field string, defaultValue string) (string, error) {
	value, err := rds.client.HGet(context.TODO(), key, field).Result()

	if err != nil {
		return defaultValue, err
	}

	return value, nil
}

func (rds *Redis) HGetInt(key, field string, defaultValue int) (int, error) {
	value, err := rds.client.HGet(context.TODO(), key, field).Result()

	if err != nil {
		return defaultValue, err
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue, err
	}

	return intValue, nil
}

func (rds *Redis) HGetInt64(key, field string, defaultValue int64) (int64, error) {
	value, err := rds.client.HGet(context.TODO(), key, field).Result()

	if err != nil {
		return defaultValue, err
	}

	int64Value, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue, err
	}

	return int64Value, nil
}

func (rds *Redis) HGetFloat64(key, field string, defaultValue float64) (float64, error) {
	value, err := rds.client.HGet(context.TODO(), key, field).Result()

	if err != nil {
		return defaultValue, err
	}

	float64Value, err := strconv.ParseFloat(value, 10)
	if err != nil {
		return defaultValue, err
	}

	return float64Value, nil
}

func (rds *Redis) HMSet(key string, data map[string]interface{}) {
	rds.client.HMSet(context.TODO(), key, data)
}

func (rds *Redis) HMGet(key string, fields []string) (map[string]interface{}, error) {

	m := make(map[string]interface{})

	for _, f := range fields {
		if r, err := rds.client.HMGet(context.TODO(), key, f).Result(); err == nil {
			m[f] = r[0]
		}
	}

	return m, nil
}

func (rds *Redis) HMGetAll(key string) (map[string]string, error) {
	m := make(map[string]string)
	m, err := rds.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		return m, err
	}

	return m, nil
}

func (rds *Redis) IncrBy(key string) (int64, error) {

	val, err := rds.client.Incr(context.TODO(), key).Result()
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (rds *Redis) DecrBy(key string) (int64, error) {

	val, err := rds.client.Decr(context.TODO(), key).Result()
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (rds *Redis) HIncrBy(key, field string, incr int64) (int64, error) {

	val, err := rds.client.HIncrBy(context.TODO(), key, field, incr).Result()
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (rds *Redis) HMGetByFields(key string, fields ...string) (map[string]interface{}, error) {

	if result, err := rds.client.HMGet(context.TODO(), key, fields...).Result(); err != nil {
		return nil, err

	} else {
		m := make(map[string]interface{})

		for i, r := range result {
			m[fields[i]] = r
		}

		return m, nil
	}
}

func (rds *Redis) Exist(key string) int64 {
	result, _ := rds.client.Exists(context.TODO(), key).Result()

	return result
}

func (rds *Redis) HExistAndGetString(key, fields string) (string, bool, error) {
	isExist, _ := rds.client.HExists(context.TODO(), key, fields).Result()
	if isExist {
		result, err := rds.client.HGet(context.TODO(), key, fields).Result()
		if err != nil {
			return "", false, err
		} else {
			return result, true, nil
		}
	} else {
		return "", false, nil
	}
}

func (rds *Redis) LPush(key string, value ...interface{}) error {
	err := rds.client.LPush(context.TODO(), key, value...).Err()

	if err != nil {
		return err
	}

	return nil
}

func (rds *Redis) LPop(key, defaultValue string) (string, error) {
	result, err := rds.client.LPop(context.TODO(), key).Result()

	if err != nil {
		return defaultValue, err
	}

	return result, nil
}

func (rds *Redis) RPush(key string, value ...interface{}) error {
	err := rds.client.RPush(context.TODO(), key, value...).Err()

	if err != nil {
		return err
	}

	return nil
}

func (rds *Redis) RPop(key, defaultValue string) (string, error) {
	result, err := rds.client.RPop(context.TODO(), key).Result()

	if err != nil {
		return defaultValue, err
	}

	return result, nil
}

func (rds *Redis) LLen(key string) (int64, error) {
	counts, err := rds.client.LLen(context.TODO(), key).Result()
	if err != nil {
		return 0, err
	}

	return counts, nil
}

func (rds *Redis) LRange(key string, start, stop int64, defaultValue []string) ([]string, error) {
	listLength, err := rds.client.LLen(context.TODO(), key).Result()
	if err != nil {
		return []string{}, errors.New("is not list")
	}

	if listLength == 0 {
		return defaultValue, errors.New("key is not exist")
	}

	if start > (listLength - 1) {
		return defaultValue, errors.New("index out of range")
	}

	if start >= 0 && stop >= 0 && start > stop {
		return defaultValue, errors.New("illegal index")
	} else if start < 0 && stop >= 0 {
		return defaultValue, errors.New("illegal index")
	} else if start < 0 && stop < 0 && start > stop {
		return defaultValue, errors.New("illegal index")
	}

	total, err := rds.client.LRange(context.TODO(), key, start, stop).Result()
	if err != nil {
		return defaultValue, err
	}

	return total, nil
}

func (rds *Redis) Expire(key string, expire time.Duration) bool {
	result, _ := rds.client.Expire(context.TODO(), key, expire).Result()

	return result
}

func (rds *Redis) TTL(key string) (time.Duration, error) {
	expire, err := rds.client.TTL(context.TODO(), key).Result()
	if err != nil {
		return 0, err
	}

	return expire, nil
}

func (rds *Redis) Scan(cursor uint64, match string, count int64) ([]string, uint64, error) {
	keys, newCursor, err := rds.client.Scan(context.TODO(), cursor, match, count).Result()

	if err != nil {
		return []string{}, 0, err
	}

	return keys, newCursor, nil
}

func (rds *Redis) Delete(key ...string) (int64, error) {
	numberOfKeyRemove, err := rds.client.Del(context.TODO(), key...).Result()

	if err != nil {
		return 0, err
	}

	return numberOfKeyRemove, nil
}

func (rds *Redis) Publish(channel string, message interface{}) error {
	_, err := rds.client.Publish(context.TODO(), channel, message).Result()

	if err != nil {
		return err
	}

	return nil
}

func (rds *Redis) SubscribeChanel(channel string) <-chan *redis.Message {
	subscriber := rds.client.Subscribe(context.TODO(), channel)
	ch := subscriber.Channel()

	return ch
}

func (rds *Redis) GetClient() *redis.Client {
	return rds.client
}

func (rds *Redis) Close() error {
	err := rds.client.Close()

	if err != nil {
		return err
	}

	return nil
}
