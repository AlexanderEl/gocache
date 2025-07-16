package gocache

import (
	"bytes"
	"testing"
	"time"
)

var ttl time.Duration = time.Millisecond * 200

func awaitTtlExpiration() {
	time.Sleep(ttl + time.Millisecond*10)
}

func TestGet(t *testing.T) {
	cache := New()

	// Missing key
	_, err := cache.Get("missing")
	if err == nil {
		t.Errorf("Expecting an error when fetching by non-existing key but no error raised")
	}

	// Existing key
	key, value := "testKey", "testValue"
	cache.Set(key, value, ttl)

	byteValue, err := cache.Get(key)
	if err != nil {
		t.Errorf("Failure to retrieve value from cache with valid key")
	}

	if !bytes.Equal(byteValue, []byte(value)) {
		t.Errorf("Byte slice comparison failed")
	}

	if string(byteValue) != value {
		t.Errorf("Retrieved string value doesn't match input string value")
	}

	awaitTtlExpiration()

	// TTL expired, key should be missing
	_, err = cache.Get(key)
	if err == nil {
		t.Errorf("Element should be removed after TTL expires, but element still found.")
	}
}

func TestSet(t *testing.T) {
	cache := New()

	key, value := "testKey", "testValue"
	cache.Set(key, value, ttl)

	// Existing element
	if !cache.Has(key) {
		t.Errorf("Failed to find freshly set element")
	}

	awaitTtlExpiration()

	// Expired element
	if cache.Has(key) {
		t.Errorf("Element with expired TTL is still found")
	}

	// Override values
	cache.Set(key, value, time.Second*10) // Long TTL
	cache.Set(key, value, ttl)            // Short TTL

	awaitTtlExpiration() // Await short TTL

	if cache.Has(key) {
		t.Errorf("Short TTL expired but long did not, element is still found")
	}
}

func TestHas(t *testing.T) {
	cache := New()

	// Non-existing element
	if cache.Has("missing") {
		t.Errorf("Found a non-existing element")
	}

	key, value := "testKey", "testValue"
	cache.Set(key, value, ttl)

	// Existing element
	if !cache.Has(key) {
		t.Errorf("Failed to find freshly set element")
	}

	awaitTtlExpiration()

	// Expired element
	if cache.Has(key) {
		t.Errorf("Element with expired TTL is still found")
	}
}

func TestDelete(t *testing.T) {
	cache := New()

	// Delete a non-existing element
	nonExistingKey := "missing"
	cache.Delete(nonExistingKey)
	_, err := cache.Get(nonExistingKey)
	if err == nil {
		t.Errorf("Retrieved a result for deleted non-existing element")
	}

	key, value := "testKey", "testValue"
	cache.Set(key, value, 0)

	// Validate it exists
	if !cache.Has(key) {
		t.Errorf("Failure to find newly inserted element")
	}

	cache.Delete(key)

	// Validate it is gone
	if cache.Has(key) {
		t.Errorf("Found a deleted element")
	}
}

func TestE2EIntegration(t *testing.T) {
	cache := New()

	// Populate elements, with and without TTL
	cache.Set("1", "one", 0)
	cache.Set("2", "two", 0)
	cache.Set("3", "three", ttl)
	cache.Set("4", "four", ttl)
	cache.Set("5", "five", ttl)

	// Validate existance of values
	if !(cache.Has("1") || cache.Has("3")) {
		t.Errorf("Failed to find set elements")
	}

	// Validate existance and value
	byteValue, err := cache.Get("1")
	if err != nil {
		t.Errorf("Failed to retrieve existing key")
	}
	if string(byteValue) != "one" {
		t.Errorf("Invalid value retrieved")
	}

	// Allow all TTL elements to expire
	awaitTtlExpiration()

	// Validate expiration
	if cache.Has("3") || cache.Has("4") || cache.Has("5") {
		t.Errorf("Cache still has expired elements")
	}

	// Validate non-expiring elements
	byteValue, err = cache.Get("2")
	if err != nil {
		t.Errorf("Element without expiration is not found after sleep")
	}
	if string(byteValue) != "two" {
		t.Errorf("Invalid value found after sleep")
	}

	// Validate deletion
	cache.Delete("1")
	if cache.Has("1") {
		t.Errorf("Deleted element still found")
	}

	cache.Delete("2")
	byteValue, err = cache.Get("2")
	if err == nil {
		t.Errorf("Deleted element was found")
	}
	if byteValue != nil {
		t.Errorf("Deleted element returned value")
	}
}
