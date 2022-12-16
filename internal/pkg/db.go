package pkg

import (
	"sync"
	"time"

	"github.com/hbina/radish/internal/types"
)

// Key-value pair.
// Will be used when serializing/deserializing redis objects.
type Kvp struct {
	Key  string `json:"key"`
	Type string `json:"type"`
	Data []byte `json:"value"`
}

// A redis database.
// There can be more than one in a redis instance.
type Db struct {
	mu           *sync.RWMutex // Lock to the database
	id           uint64
	storage      map[string]types.Item
	storageTtl   map[string]time.Time
	redis        *Redis
	bcmd         map[*Client]BlockedCommand
	bcmdTtl      chan *Client
	modifiedKeys chan string
}

// NewRedisDb creates a new db.
func NewRedisDb(id uint64, r *Redis) *Db {
	db := &Db{
		mu:         new(sync.RWMutex),
		id:         id,
		storage:    make(map[string]types.Item, 0),
		storageTtl: make(map[string]time.Time, 0),
		redis:      r,
	}
	return db
}

// RedisDbs gets all redis databases.
func (r *Redis) RedisDbs() map[uint64]*Db {
	return r.dbs
}

// Redis gets the redis instance.
func (db *Db) Redis() *Redis {
	return db.redis
}

// Id gets the db id.
func (db *Db) Id() uint64 {
	return db.id
}

// Sets a key with an item which can have an expiration time.
func (db *Db) Set(key string, i types.Item, ttl time.Time) types.Item {
	// Empty item is considered a delete operation because
	// operations on non-existent key is equivalent to zeroth of that
	// object type.
	// TODO: Should this be behavior of set or the specific commands?
	if i.Type() == types.ValueTypeString {
		// Except or strings?
	} else if i.Type() == types.ValueTypeList {
		list := i.(*types.List)

		if list.Len() == 0 {
			db.Delete(key)
			return nil
		}
	} else if i.Type() == types.ValueTypeSet {
		set := i.(*types.Set)

		if set.Len() == 0 {
			db.Delete(key)
			return nil
		}
	} else if i.Type() == types.ValueTypeZSet {
		str := i.(*types.ZSet)

		if str.Len() == 0 {
			db.Delete(key)
			return nil
		}
	}

	old, exists := db.storage[key]

	// Insert new value to a key will overwrite everything about it
	db.storage[key] = i
	db.storageTtl[key] = ttl

	if exists {
		return old
	} else {
		return nil
	}
}

// GetExpiry returns the item by the key or nil if key does not exists.
func (db *Db) GetExpiry(key string) (time.Time, bool) {
	v, e := db.storageTtl[key]
	return v, e
}

// SetExpiry sets the expiry of a key
func (db *Db) SetExpiry(key string, ttl time.Time) (time.Time, bool) {
	old, exists := db.storageTtl[key]
	db.storageTtl[key] = ttl
	return old, exists
}

// Deletes a key, returns number of deleted keys.
func (db *Db) Delete(keys ...string) int {
	var c int
	for _, k := range keys {
		_, itemExists := db.storage[k]
		_, ttlExists := db.storageTtl[k]
		delete(db.storage, k)
		delete(db.storageTtl, k)

		if itemExists && ttlExists {
			c++
		}
	}

	return c
}

func (db *Db) DeleteExpired(keys ...string) int {
	var c int
	for _, k := range keys {
		if db.Expired(k) && db.Delete(k) > 0 {
			c++
		}
	}
	return c
}

// Get gets the item or nil if expired or not exists. If 'deleteIfExpired' is true the key will be deleted.
// TODO: Should this return the exists bool or its enough to return nil?
func (db *Db) Get(key string) (types.Item, time.Time) {
	value, exists := db.storage[key]
	if !exists {
		return nil, time.Time{}
	}
	if db.Expired(key) {
		db.Delete(key)
		return nil, time.Time{}
	}
	return value, db.storageTtl[key]
}

// IsEmpty checks if db is empty.
func (db *Db) IsEmpty() bool {
	return len(db.storage) == 0
}

// HasExpiringKeys checks if db has any expiring keys.
func (db *Db) HasExpiringKeys() bool {
	return len(db.storageTtl) != 0
}

// Exists return whether or not a key exists.
// Internally, it has the side effect of evicting keys that
// expires.
func (db *Db) Exists(key string) bool {
	maybeItem, _ := db.Get(key)
	return maybeItem != nil
}

// Expired only check if a key can and is expired.
func (db *Db) Expired(key string) bool {
	ttl, exists := db.Expiry(key)
	// Since we always write ttl in Set, we need to
	// check if its zero.
	if !exists || time.Time.IsZero(ttl) {
		return false
	}
	return time.Now().After(ttl)
}

// Expiry gets the expiry of the key has one.
func (db *Db) Expiry(key string) (time.Time, bool) {
	val, ok := db.storageTtl[key]
	return val, ok
}

// DeleteExpiredKeys will delete all the keys that have expired TTL.
func (db *Db) DeleteExpiredKeys() int {
	count := 0
	for k := range db.storageTtl {
		count += db.DeleteExpired(k)
	}
	return count
}

func (db *Db) Clear() {
	for k := range db.storage {
		delete(db.storage, k)
		delete(db.storageTtl, k)
	}
}

// Number of keys in the storage
func (db *Db) Len() int {
	return len(db.storage)
}

func (db *Db) Lock() {
	db.mu.Lock()
}

func (db *Db) UnLock() {
	db.mu.Unlock()
}

func (db *Db) StartKeyExpiryJob(tick time.Duration) {
	ticker := time.NewTicker(tick)
	for range ticker.C {
		db.Lock()
		db.DeleteExpiredKeys()
		db.UnLock()
	}
}

func (db *Db) StartBcmdTimeoutJob(tick time.Duration) {
	for c := range db.bcmdTtl {
		db.Lock()
		delete(db.bcmd, c)
		db.UnLock()
	}
}
