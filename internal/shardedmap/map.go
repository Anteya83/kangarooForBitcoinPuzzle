package shardedmap

import (
	"math/big"
	"sync"
)

type Shard struct {
	mu sync.RWMutex
	m  map[string]*big.Int
}

type ShardedMap struct {
	shards []*Shard
}

func NewShardedMap(numShards int) *ShardedMap {
	shards := make([]*Shard, numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = &Shard{m: make(map[string]*big.Int)}
	}
	return &ShardedMap{shards: shards}
}

func (sm *ShardedMap) getShard(key string) *Shard {
	h := uint64(0)
	for _, c := range key {
		h = h*31 + uint64(c)
	}
	return sm.shards[h%uint64(len(sm.shards))]
}

func (sm *ShardedMap) Set(key string, value *big.Int) {
	shard := sm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.m[key] = new(big.Int).Set(value)
}

func (sm *ShardedMap) Get(key string) (*big.Int, bool) {
	shard := sm.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	val, ok := shard.m[key]
	if !ok {
		return nil, false
	}
	return new(big.Int).Set(val), true
}

func (sm *ShardedMap) Len() int {
	total := 0
	for _, shard := range sm.shards {
		shard.mu.RLock()
		total += len(shard.m)
		shard.mu.RUnlock()
	}
	return total
}

func (sm *ShardedMap) Delete(key string) {
	shard := sm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	delete(shard.m, key)
}

func (sm *ShardedMap) ForEach(fn func(key string, val *big.Int)) {
	for _, shard := range sm.shards {
		shard.mu.RLock()
		for k, v := range shard.m {
			fn(k, v)
		}
		shard.mu.RUnlock()
	}
}
