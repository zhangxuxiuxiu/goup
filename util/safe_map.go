package util

import "sync"

type SafeMap[K comparable, V any] interface {
	Load(k K) (V, bool)
	Store(k K, v V)
	Delete(k K)
}

type syncMap[K comparable, V any] struct {
	data sync.Map
}

func NewSyncMap[K comparable, V any]() SafeMap[K, V] {
	return &syncMap[K, V]{}
}

func (sm *syncMap[K, V]) Load(k K) (V, bool) {
	v, b := sm.data.Load(k)
	return v.(V), b
}

func (sm *syncMap[K, V]) Store(k K, v V) {
	sm.data.Store(k, v)
}

func (sm *syncMap[K, V]) Delete(k K) {
	sm.data.Delete(k)
}

type mutexMap[K comparable, V any] struct {
	mu   sync.Mutex
	data map[K]V
}

func NewMutexMap[K comparable, V any]() SafeMap[K, V] {
	return &mutexMap[K, V]{data: make(map[K]V)}
}

func (sm *mutexMap[K, V]) Load(k K) (V, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	v, b := sm.data[k]
	return v, b
}

func (sm *mutexMap[K, V]) Store(k K, v V) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[k] = v
}

func (sm *mutexMap[K, V]) Delete(k K) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, k)
}

type bucketSyncMap[K comparable, V any] struct {
	buckets []SafeMap[K, V]
	hash    func(K) int
}

func NewBucketSyncMap[K comparable, V any](bucketCnt int, h func(K) int, gen func() SafeMap[K, V]) SafeMap[K, V] {
	bsm := bucketSyncMap[K, V]{buckets: make([]SafeMap[K, V], bucketCnt), hash: h}
	for idx := range bsm.buckets {
		bsm.buckets[idx] = gen()
	}
	return &bsm
}

func (bsm *bucketSyncMap[K, V]) bucketOf(k K) SafeMap[K, V] {
	return bsm.buckets[bsm.hash(k)%len(bsm.buckets)]
}

func (bsm *bucketSyncMap[K, V]) Load(k K) (V, bool) {
	return bsm.bucketOf(k).Load(k)
}

func (bsm *bucketSyncMap[K, V]) Store(k K, v V) {
	bsm.bucketOf(k).Store(k, v)
}

func (bsm *bucketSyncMap[K, V]) Delete(k K) {
	bsm.bucketOf(k).Delete(k)
}
