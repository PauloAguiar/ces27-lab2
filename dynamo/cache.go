package dynamo

import (
    "log"
    "sync"
)

///////////////////////////////
// THIS FILE IS MISSING CODE //
///////////////////////////////

// Cache is the struct that handle all the data storage for the dynamo server.
type Cache struct {
    data map[string]string
    sync.Mutex
}

// Create a new cache object and return a pointer to it.
func NewCache() *Cache {
    var s Cache

    s.data = make(map[string]string)

    return &s
}

// Get the value of a key from the storage. This will handle concurrent get
// requests by locking the structure.
func (cache *Cache) Get(key string) (value string, timestamp int64) {
    cache.Lock()
    value = cache.data[key]
    timestamp = 0
    cache.Unlock()

    log.Printf("[CACHE] Getting Key '%v' with Value '%v' @ timestamp '%v'\n", key, value, timestamp)
    return
}

// Put a value to a key in the storage. This will handle concurrent put
// requests by locking the structure.
func (cache *Cache) Put(key string, value string, timestamp int64) {
    log.Printf("[CACHE] Putting Key '%v' with Value '%v' @ timestamp '%v'\n", key, value, timestamp)

    cache.Lock()
    cache.data[key] = value
    cache.Unlock()

    return
}

// Retrieve all information from the server. This shouldn't be used in any way
// except for testing purposes.
func (cache *Cache) getAll() (data map[string]string, timestamps map[string]int64) {
    data = cache.data
    timestamps = make(map[string]int64)
    return data, timestamps
}
