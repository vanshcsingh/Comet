package cache

import (
	"testing"
)

func setupBatcher() {
	
}

func setupCache() {

}

// TestFetch tests the Fetch method of our local cache
func TestFetch(t *testing.T) {
	setupBatcher()
	setupCache()
}

func TestRequest(t *testing.T) {
	setupBatcher()
	cache := setupCache()
}
