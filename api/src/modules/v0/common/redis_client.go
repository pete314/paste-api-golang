//Author: Peter Nagy <https://peternagy.ie>
//Since: 06, 2017
//Description: wrapper for redis

package common

import (
	"gopkg.in/redis.v5"
	"strconv"
	"time"
)

//RedisClient - the data belonging to this object
type RedisClient struct {
	client *redis.Client
}

var (
	instance *RedisClient
)

//GetRedisClient - singleton for this
func GetRedisClient() *RedisClient {
	once.Do(func() {
		instance = &RedisClient{client: newRedisClient()}
	})
	return instance
}

//NewRedisClient - Get new redis client
func NewRedisClient(config *Config) *RedisClient {
	return &RedisClient{client: newRedisClient()}
}

//Get new initialized redis client
func newRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     RuntimeConfig.RedisAddress + ":" + strconv.Itoa(RuntimeConfig.RedisPort),
		Password: RuntimeConfig.RedisPassword,
		DB:       RuntimeConfig.RedisDb,
	})
}

//Info - Get cluster info
func (rc RedisClient) Info(param ...string) (string, error) {
	return rc.client.Info(param[0]).Result()
}

//SetKey - Set key value pair
func (rc RedisClient) SetKey(key string, value interface{}, ttl int) bool {
	result, err := rc.client.SetNX(key, value, time.Duration(ttl)*time.Second).Result()
	CheckError("Redis client: can not create hash", err, false)

	return result
}

//GetKey - get the key value
func (rc RedisClient) GetKey(key string) []byte {
	result, _ := rc.client.Get(key).Bytes()

	return result
}

//DeleteKey - delete key value
func (rc RedisClient) DeleteKey(key string) int64 {
	result, err := rc.client.Del(key).Result()
	CheckError("Redis client: can not delete hash", err, false)

	return result
}

//GetKeyNames - get the key names
func (rc RedisClient) GetKeyNames(key string) []string {
	result, err := rc.client.Keys(key).Result()
	CheckError("Redis client: can not get keys", err, false)

	return result
}

//KeyExists - check if key exists
func (rc RedisClient) KeyExists(key string) (bool, error) {
	return rc.client.Exists(key).Result()
}

//GetListLength - Get the list length for key
func (rc RedisClient) GetListLength(key string) int64 {
	result, err := rc.client.LLen(key).Result()
	CheckError("Redis client: can not get key length", err, false)

	return result
}

//GetListRange - Get list item data in range
func (rc RedisClient) GetListRange(key string, from, until int64) []string {
	result, err := rc.client.LRange(key, from, until).Result()
	CheckError("Redis client: can not get list items", err, false)

	return result
}

//PushToList - push to list
func (rc RedisClient) PushToList(key, data string) bool {
	result, err := rc.client.LPush(key, data).Result()
	CheckError("redis_client::PushToList can not push items to list", err, false)

	return result > 0
}
