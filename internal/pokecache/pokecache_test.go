package pokecache

import (
	"testing"
	"time"
)

// Test the basic functionality of the cache including:
// Add, Get, and automatic reaping of old entries.
func TestPokecache(t *testing.T) {
	const interval = time.Second
	cases := []struct { // key and value pairs to test
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("example data"),
		},
		{
			key: "https://example.com/2",
			val: []byte("example data 2"),
		},
	}

	for i, c := range cases {
		// Make a new cache instance for each test case
		cache := NewCache(interval)

		// Add and retrieve the key-value pair
		cache.Add(c.key, c.val)
		val, ok := cache.Get(c.key)
		if !ok {
			t.Errorf("Test case %d: Expected to find key", i)
		}
		// Check correctness of the value
		if string(val) != string(c.val) {
			t.Errorf("Test case %d: Expected to find value", i)
		}

		// Wait for the interval to pass.
		// The reapLoop should have removed the entry afterwards.
		time.Sleep(interval + time.Millisecond*100)
		val, ok = cache.Get(c.key)
		if ok {
			t.Errorf("Test case %d: Expected key to be reaped", i)
		}
	}
}
