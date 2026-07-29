package shardedmap

import (
	"math/big"
	"sync"
	"testing"
)

func TestSetGet(t *testing.T) {
	sm := NewShardedMap(4)
	key := "test_key"
	val := big.NewInt(12345)

	sm.Set(key, val)
	got, ok := sm.Get(key)
	if !ok {
		t.Fatal("key not found")
	}
	if got.Cmp(val) != 0 {
		t.Errorf("got %v, want %v", got, val)
	}
}

func TestLen(t *testing.T) {
	sm := NewShardedMap(4)
	sm.Set("a", big.NewInt(1))
	sm.Set("b", big.NewInt(2))
	sm.Set("c", big.NewInt(3))

	if sm.Len() != 3 {
		t.Errorf("Len = %d, want 3", sm.Len())
	}
}

func TestDelete(t *testing.T) {
	sm := NewShardedMap(4)
	sm.Set("x", big.NewInt(10))
	sm.Delete("x")
	_, ok := sm.Get("x")
	if ok {
		t.Error("key should be deleted")
	}
}

// / if not panic test passed
func TestConcurrency(t *testing.T) {
	sm := NewShardedMap(8)
	const numGoroutines = 100
	const numOps = 1000
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := "key"
				val := big.NewInt(int64(id*1000 + j))
				sm.Set(key, val)
				sm.Get(key)
				sm.Len()
			}
		}(i)
	}
	wg.Wait()

}
