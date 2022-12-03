package redis

import (
	"log"
	"time"
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
	id uint64

	// All storage in this db.
	storage map[string]Item

	// TTL of each keys
	ttl map[string]time.Time

	// TODO: Some statistics about the database that might be useful
	// when we have eviction policies and stuff like that.

	redis *Redis
}

// The item interface. An item is the value of a key.
type Item interface {
	// The pointer to the value.
	Value() interface{}

	// The id of the type of the Item.
	// This need to be constant for the type because it is
	// used when de-/serializing item from/to disk.
	Type() uint64
	TypeFancy() string
}

// NewRedisDb creates a new db.
func NewRedisDb(id uint64, r *Redis) *RedisDb {
	return &RedisDb{
		id:      id,
		redis:   r,
		storage: make(map[string]Item, 0),
		ttl:     make(map[string]time.Time, 0),
	}
}

// RedisDbs gets all redis databases.
func (r *Redis) RedisDbs() map[uint64]*RedisDb {
	return r.redisDbs
}

// Redis gets the redis instance.
func (db *RedisDb) Redis() *Redis {
	return db.redis
}

// Id gets the db id.
func (db *RedisDb) Id() uint64 {
	return db.id
}

// Sets a key with an item which can have an expiration time.
func (db *RedisDb) Set(key string, i Item, ttl time.Time) Item {
	// Empty item is considered a delete operation because
	// operations on non-existent key is equivalent to zeroth of that
	// object type.
	// TODO: Should this be behavior of set or the specific commands?
	if i.Type() == ValueTypeString {
		// Except or strings?
	} else if i.Type() == ValueTypeList {
		list := i.(*List)

		if list.Len() == 0 {
			db.Delete(key)
			return nil
		}
	} else if i.Type() == ValueTypeSet {
		set := i.(*Set)

		if set.Len() == 0 {
			db.Delete(key)
			return nil
		}
	} else if i.Type() == ValueTypeZSet {
		str := i.(*ZSet)

		if str.Len() == 0 {
			db.Delete(key)
			return nil
		}
	}

	old, exists := db.storage[key]

	// Insert new value to a key will overwrite everything about it
	db.storage[key] = i
	db.ttl[key] = ttl

	if exists {
		return old
	} else {
		return nil
	}
}

// Get returns the item by the key or nil if key does not exists.
// TODO: Should this returns the exists bool?
func (db *RedisDb) Get(key string) Item {
	return db.storage[key]
}

// GetExpiry returns the item by the key or nil if key does not exists.
func (db *RedisDb) GetExpiry(key string) (time.Time, bool) {
	v, e := db.ttl[key]
	return v, e
}

// SetExpiry sets the expiry of a key
func (db *RedisDb) SetExpiry(key string, ttl time.Time) (time.Time, bool) {
	if time.Time.IsZero(ttl) {
		delete(db.ttl, key)
		return time.Time{}, false
	} else {
		old, exists := db.ttl[key]
		db.ttl[key] = ttl
		return old, exists
	}
}

// Deletes a key, returns number of deleted keys.
func (db *RedisDb) Delete(keys ...string) int {
	var c int
	for _, k := range keys {
		_, itemExists := db.storage[k]
		_, ttlExists := db.ttl[k]
		delete(db.storage, k)
		delete(db.ttl, k)

		if itemExists != ttlExists {
			log.Printf("Invariant failure: %t != %t when checking if key exists in storage and expiringKeys", itemExists, ttlExists)
		}

		if itemExists && ttlExists {
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
// TODO: Should this return the exists bool or its enough to return nil?
func (db *RedisDb) GetOrExpire(key string, deleteIfExpired bool) (Item, time.Time) {
	value, exists := db.storage[key]
	if !exists {
		return nil, time.Time{}
	}
	if db.Expired(key) {
		if deleteIfExpired {
			db.Delete(key)
		}
		return nil, time.Time{}
	}
	return value, db.ttl[key]
}

// IsEmpty checks if db is empty.
func (db *RedisDb) IsEmpty() bool {
	return len(db.storage) == 0
}

// HasExpiringKeys checks if db has any expiring keys.
func (db *RedisDb) HasExpiringKeys() bool {
	return len(db.ttl) != 0
}

// Exists return whether or not a key exists.
// Internally, it has the side effect of evicting keys that
// expires.
func (db *RedisDb) Exists(key string) bool {
	maybeItem, _ := db.GetOrExpire(key, true)
	return maybeItem != nil
}

// Check if key has an expiry set.
func (db *RedisDb) Expires(key string) bool {
	_, ok := db.ttl[key]
	return ok
}

// Expired only check if a key can and is expired.
func (db *RedisDb) Expired(key string) bool {
	ttl, exists := db.Expiry(key)
	// Since we always write ttl in Set, we need to
	// check if its zero.
	if !exists || time.Time.IsZero(ttl) {
		return false
	}
	return db.Expires(key) && time.Now().After(ttl)
}

// Expiry gets the expiry of the key has one.
func (db *RedisDb) Expiry(key string) (time.Time, bool) {
	val, ok := db.ttl[key]
	return val, ok
}

// DeleteExpiredKeys will delete all the keys that have expired TTL.
func (db *RedisDb) DeleteExpiredKeys() int {
	count := 0
	for k := range db.ttl {
		count += db.DeleteExpired(k)
	}
	return count
}

func (db *RedisDb) Clear() {
	for k := range db.storage {
		delete(db.storage, k)
		delete(db.ttl, k)
	}
}

// Number of keys in the storage
func (db *RedisDb) Len() int {
	return len(db.storage)
}
