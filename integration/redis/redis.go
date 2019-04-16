package redis

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

var pool redis.Pool

func init() {
	pool = redis.Pool{
		MaxIdle:   50,
		MaxActive: 500, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "cityos.utdxew.ng.0001.euw1.cache.amazonaws.com:6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// IncrementByFloat will increment the key by the value
func IncrementByFloat(key, value string) {
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("INCRBYFLOAT", key, value); err != nil {
		log.Fatal(err)
	}
}

// IncrementBy will increment the key by the value
func IncrementBy(key, value string) {
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("INCRBY", key, value); err != nil {
		log.Fatal(err)
	}
}

// IncrementByHashFloat will increment the key by the value for the hash
func IncrementByHashFloat(hash, key, value string) {
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("HINCRBYFLOAT", hash, key, value); err != nil {
		log.Fatal(err)
	}
}

// IncrementByHash will increment the key by the value for the hash
func IncrementByHash(hash, key, value string) {
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("HINCRBY", hash, key, value); err != nil {
		log.Fatal(err)
	}
}

// GetInt will return the int value
func GetInt(key string) (int, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("GET", key))
}

// GetIntHash will return the int value from the hash
func GetIntHash(hash, key string) (int, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("HGET", hash, key))
}

// GetFloat will return the float value
func GetFloat(key string) (float64, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Float64(conn.Do("GET", key))
}

// GetFloatHash will return the float value from the hash
func GetFloatHash(hash, key string) (float64, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Float64(conn.Do("HGET", hash, key))
}

// GetString will return the string value
func GetString(key string) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.String(conn.Do("GET", key))
}

// GetStringList will return the list of strings
func GetStringList(key string) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("GET", key))
}

// Keys will return the list of Keys based on teh provided filter
func Keys(key string) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("KEYS", key))
}

// Set will save the value
func Set(key, value string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	return err
}

// SetHash will save the value for the hash
func SetHash(hash, key, value string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", hash, key, value)
	return err
}

// Del will delete the provided key
func Del(key string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

// FlushDB will delete all keys in teh currently loaded DB
func FlushDB() error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("FLUSHDB")
	return err
}

// Expire will Expire the provided key
func Expire(key, seconds string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", key, seconds)
	return err
}
