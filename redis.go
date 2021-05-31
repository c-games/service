package service

import (
	"errors"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

type Redis struct {
	client *redis.Client
}
type RedisParameter struct {
	Network      string
	Address      string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
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

	_, err := client.Ping().Result()

	return &Redis{client: client}, err

}

func (rds *Redis) Set(key string, value interface{}) error {
	err := rds.client.Set(key, value, 0).Err()

	if err != nil {
		return err
	}

	return nil
}

func (rds *Redis) GetSting(key string, defaultValue string) (string, error) {
	value, err := rds.client.Get(key).Result()

	if err != nil {
		return defaultValue, err
	}

	return value, nil
}

func (rds *Redis) GetInt(key string, defaultValue int) (int, error) {
	value, err := rds.client.Get(key).Result()

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
	value, err := rds.client.Get(key).Result()

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
	value, err := rds.client.Get(key).Result()

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
	value, err := rds.client.HGet(key, field).Result()

	if err != nil {
		return defaultValue, err
	}

	return value, nil
}

func (rds *Redis) HGetInt(key, field string, defaultValue int) (int, error) {
	value, err := rds.client.HGet(key, field).Result()

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
	value, err := rds.client.HGet(key, field).Result()

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
	value, err := rds.client.HGet(key, field).Result()

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
	rds.client.HMSet(key, data)
}

func (rds *Redis) HMGet(key string, fields []string) (map[string]interface{}, error) {

	m := make(map[string]interface{})

	for _, f := range fields {
		if r, err := rds.client.HMGet(key, f).Result(); err == nil {
			m[f] = r[0]
		}
	}

	return m, nil
}

func (rds *Redis) IncrBy(key string) (int64, error) {

	val, err := rds.client.Incr(key).Result()
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (rds *Redis) HIncrBy(key, field string, incr int64) (int64, error) {

	val, err := rds.client.HIncrBy(key, field, incr).Result()
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (rds *Redis) HMGetByFields(key string, fields ...string) (map[string]interface{}, error) {

	if result, err := rds.client.HMGet(key, fields...).Result(); err != nil {
		return nil, err

	} else {
		m := make(map[string]interface{})

		for i, r := range result {
			m[fields[i]] = r
		}

		return m, nil
	}
}

func (rds *Redis) Exist (key string) int64 {
	result, _ := rds.client.Exists(key).Result()

	return result
}

func (rds *Redis) HExistAndGetString (key, fields string) (string, bool, error) {
	isExist, _ := rds.client.HExists(key, fields).Result()
	if isExist {
		result, err := rds.client.HGet(key, fields).Result()
		if err != nil {
			return "", false, err
		} else {
			return  result, true, nil
		}
	} else {
		return "", false, nil
	}
}

func (rds *Redis) LPush(key string, value ...interface{}) error {
	err := rds.client.LPush(key, value...).Err()

	if err != nil {
		return err
	}

	return nil
}

func (rds *Redis) LRange(key string, start, stop int64, defaultValue []string) ([]string, error) {
	listLength, err := rds.client.LLen(key).Result()
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

	total, err := rds.client.LRange(key, start, stop).Result()
	if err != nil {
		return defaultValue, err
	}

	return total, nil
}

func (rds *Redis) Expire (key string, expire time.Duration) bool {
	result, _ := rds.client.Expire(key, expire).Result()

	return result
}

func (rds *Redis) Scan (cursor uint64, match string, count int64) ([]string, uint64, error) {
	keys, newCursor, err := rds.client.Scan(cursor, match, count).Result()

	if err != nil {
		return []string{}, 0, err
	}

	return keys, newCursor, nil
}

func (rds *Redis) Delete (key ...string) (int64, error) {
	numberOfKeyRemove, err := rds.client.Del(key...).Result()

	if err != nil {
		return 0, err
	}

	return numberOfKeyRemove, nil
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
