package service

import (
	"math"
	"sync"
	"sync/atomic"
)

type ConcurrentMap struct {
	shard      []*MapShard
	count      int32
	shardCount uint32
}

type MapShard struct {
	m     map[string]interface{}
	mutex sync.RWMutex
}

//该参数转成二进制，每个位都赋为1
func computeCapacity(param int) int {
	if param <= 16 {
		return 16
	}
	n := param - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return math.MaxInt32
	}
	return n + 1
}

// MakeConcurrentMap 返回一个分段锁的实例
func MakeConcurrentMap(shardCount int) *ConcurrentMap {
	shardCount = computeCapacity(shardCount)
	shard := make([]*MapShard, shardCount)
	for idx := range shard {
		shard[idx] = &MapShard{
			m:     make(map[string]interface{}),
			mutex: sync.RWMutex{},
		}
	}
	return &ConcurrentMap{
		shard:      shard,
		count:      0,
		shardCount: uint32(shardCount),
	}
}

const prime32 = uint32(16777619)

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
func (dict *ConcurrentMap) getShardMap(key string) *MapShard {
	hashCode := fnv32(key)
	idx := (dict.shardCount - 1) & hashCode
	return dict.shard[idx]
}
func (dict *ConcurrentMap) Get(key string) (val interface{}, exists bool) {
	shard := dict.getShardMap(key)
	shard.mutex.RLock()
	defer shard.mutex.RUnlock()
	val, exists = shard.m[key]
	return
}
func (dict *ConcurrentMap) Len() int {
	return int(atomic.LoadInt32(&dict.count))
}

// Set 插入
func (dict *ConcurrentMap) Set(key string, val interface{}) int {
	shard := dict.getShardMap(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	if _, ok := shard.m[key]; ok {
		shard.m[key] = val
		return 0
	}
	shard.m[key] = val
	atomic.AddInt32(&dict.count, 1)
	return 1
}

// Remove 删除一个key
func (dict *ConcurrentMap) Remove(key string) int {
	shard := dict.getShardMap(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	if _, ok := shard.m[key]; ok {
		delete(shard.m, key)
		atomic.AddInt32(&dict.count, -1)
		return 1
	}
	return 0
}


