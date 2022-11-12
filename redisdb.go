package redis

import (
	"time"
)

const (
	keysMapSize           = 32
	redisDbMapSizeDefault = 3
)

const (
	ValueTypeList = iota
	ValueTypeString
	ValueTypeSet
	ValueTypeZSet
)

const (
	ValueTypeFancyList   = "list"
	ValueTypeFancyString = "string"
	ValueTypeFancySet    = "set"
	ValueTypeFancyZSet   = "zset"
)

// Key-value pair.
// Will be used when serializing/deserializing redis objects.
type Kvp struct {
	Key   string      `json:"key"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// A redis database.
// There can be more than one in a redis instance.
type RedisDb struct {
	// Database id
	id DatabaseId

	// All storage in this db.
	storage KeyValue

	// Keys with expire timestamp.
	expiringKeys ExpiringKeys

	// TODO long long avg_ttl;          /* Average TTL, just for stats */

	redis *Redis
}

// Database id
type DatabaseId uint64

// Key-Item map
type KeyValue map[string]Item

// Keys with expire timestamp.
type ExpiringKeys map[string]time.Time

// The item interface. An item is the value of a key.
type Item interface {
	// The pointer to the value.
	Value() interface{}

	// The id of the type of the Item.
	// This need to be constant for the type because it is
	// used when de-/serializing item from/to disk.
	Type() uint64
	TypeFancy() string

	// OnDelete is triggered before the key of the item is deleted.
	// db is the affected database.
	OnDelete(key string, db RedisDb)
}

// NewRedisDb creates a new db.
func NewRedisDb(id DatabaseId, r *Redis) *RedisDb {
	return &RedisDb{
		id:           id,
		redis:        r,
		storage:      make(KeyValue, keysMapSize),
		expiringKeys: make(ExpiringKeys, keysMapSize),
	}
}

// RedisDbs gets all redis databases.
func (r *Redis) RedisDbs() map[DatabaseId]*RedisDb {
	return r.redisDbs
}

// Redis gets the redis instance.
func (db *RedisDb) Redis() *Redis {
	return db.redis
}

// Id gets the db id.
func (db *RedisDb) Id() DatabaseId {
	return db.id
}

// Sets a key with an item which can have an expiration time.
func (db *RedisDb) Set(key string, i Item, expiry time.Time) Item {
	old := db.storage[key]
	db.storage[key] = i
	if !time.Time.IsZero(expiry) {
		db.expiringKeys[key] = expiry
	}
	return old
}

// Returns the item by the key or nil if key does not exists.
func (db *RedisDb) Get(key string) Item {
	return db.storage[key]
}

// Deletes a key, returns number of deleted keys.
func (db *RedisDb) Delete(keys ...string) int {
	do := func(k string) bool {
		// value.OnDelete(k, *db)
		delete(db.storage, k)
		delete(db.expiringKeys, k)
		return true
	}

	var c int
	for _, k := range keys {
		if do(k) {
			c++
		}
	}

	return c
}

func (db *RedisDb) DeleteExpired(keys ...string) int {
	var c int
	for _, k := range keys {
		if db.Expired(k) && db.Delete(k) > 0 {
			c++
		}
	}
	return c
}

// GetOrExpire gets the item or nil if expired or not exists. If 'deleteIfExpired' is true the key will be deleted.
func (db *RedisDb) GetOrExpire(key string, deleteIfExpired bool) (Item, time.Time) {
	value, exists := db.storage[key]
	if !exists {
		return nil, time.Time{}
	}
	if db.expired(key) {
		if deleteIfExpired {
			db.Delete(key)
		}
		return nil, time.Time{}
	}
	return value, db.expiringKeys[key]
}

// IsEmpty checks if db is empty.
func (db *RedisDb) IsEmpty() bool {
	return len(db.storage) == 0
}

// HasExpiringKeys checks if db has any expiring keys.
func (db *RedisDb) HasExpiringKeys() bool {
	return len(db.expiringKeys) != 0
}

// Check if key exists.
func (db *RedisDb) Exists(key *string) bool {
	return db.exists(key)
}

func (db *RedisDb) exists(key *string) bool {
	_, ok := db.storage[*key]
	return ok
}

// Check if key has an expiry set.
func (db *RedisDb) Expires(key string) bool {
	return db.expires(key)
}

func (db *RedisDb) expires(key string) bool {
	_, ok := db.expiringKeys[key]
	return ok
}

// Expired only check if a key can and is expired.
func (db *RedisDb) Expired(key string) bool {
	return db.expired(key)
}

func (db *RedisDb) expired(key string) bool {
	return db.expires(key) && TimeExpired(db.expiry(key))
}

// Expiry gets the expiry of the key has one.
func (db *RedisDb) Expiry(key string) time.Time {
	return db.expiry(key)
}

func (db *RedisDb) expiry(key string) time.Time {
	return db.expiringKeys[key]
}

// Keys gets all keys in this db.
func (db *RedisDb) Keys() KeyValue {
	return db.storage
}

// ExpiringKeys gets keys with an expiry set and their timeout.
func (db *RedisDb) ExpiringKeys() ExpiringKeys {
	return db.expiringKeys
}

// TimeExpired check if a timestamp is older than now.
func TimeExpired(expireAt time.Time) bool {
	return time.Now().After(expireAt)
}

func (db *RedisDb) SyncFlushAll() {
	for k, i := range db.storage {
		i.OnDelete(k, *db)
		delete(db.storage, k)
		delete(db.expiringKeys, k)
	}
}
