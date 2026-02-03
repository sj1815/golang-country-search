package cache

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewInMemoryCache tests the creation of a new InMemoryCache instance.
func TestNewInMemoryCache(t *testing.T) {
	c := NewInMemoryCache()

	assert.NotNil(t, c)
	assert.NotNil(t, c.store)
}

// TestCache_SetAndGet tests setting and getting values in the cache.
func TestCache_SetAndGet(t *testing.T) {
	c := NewInMemoryCache()

	// Test Set and Get
	c.Set("key1", "value1")
	value, found := c.Get("key1")

	assert.True(t, found)
	assert.Equal(t, "value1", value)
}

// TestCache_GetNotFound tests getting a value that does not exist in the cache.
func TestCache_GetNotFound(t *testing.T) {
	c := NewInMemoryCache()

	value, found := c.Get("nonexistent")

	assert.False(t, found)
	assert.Nil(t, value)
}

// TestCache_OverwriteValue tests overwriting an existing value in the cache.
func TestCache_OverwriteValue(t *testing.T) {
	c := NewInMemoryCache()

	c.Set("key1", "value1")
	c.Set("key1", "value2")

	value, found := c.Get("key1")

	assert.True(t, found)
	assert.Equal(t, "value2", value)
}

// TestCache_DifferentTypes tests setting and getting values of different types in the cache.
func TestCache_DifferentTypes(t *testing.T) {
	c := NewInMemoryCache()

	// Test with different types
	c.Set("string", "hello")
	c.Set("int", 42)
	c.Set("struct", struct{ Name string }{"test"})

	strVal, _ := c.Get("string")
	intVal, _ := c.Get("int")
	structVal, _ := c.Get("struct")

	assert.Equal(t, "hello", strVal)
	assert.Equal(t, 42, intVal)
	assert.Equal(t, struct{ Name string }{"test"}, structVal)
}

// TestCache_ConcurrentAccess tests concurrent access to the cache.
func TestCache_ConcurrentAccess(t *testing.T) {
	c := NewInMemoryCache()
	var wg sync.WaitGroup

	// Launch multiple goroutines to read and write concurrently
	for i := 0; i < 100; i++ {
		wg.Add(2)

		// Writer goroutine
		go func(val int) {
			defer wg.Done()
			c.Set("key", val)
		}(i)

		// Reader goroutine
		go func() {
			defer wg.Done()
			c.Get("key")
		}()
	}

	wg.Wait()
}

// TestCache_ConcurrentDifferentKeys tests concurrent access to different keys in the cache.
func TestCache_ConcurrentDifferentKeys(t *testing.T) {
	c := NewInMemoryCache()
	var wg sync.WaitGroup

	// Launch multiple goroutines to read and write different keys concurrently
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			key := string(rune('a' + val%26))
			c.Set(key, val)
			c.Get(key)
		}(i)
	}

	wg.Wait()
}
