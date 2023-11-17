package main

import (
	"fmt"
	"github.com/matinkhosravani/cloud-native-go-book-notes/vertical-sharding/sharded"
	"sync"
	"time"
)

type unshardedMap struct {
	sync.RWMutex
	m map[string]interface{}
}

func newUnshardedMap() *unshardedMap {
	return &unshardedMap{m: make(map[string]interface{})}
}

func main() {
	unshardedMap := newUnshardedMap()
	t := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			unshardedMap.Lock()
			unshardedMap.m[fmt.Sprintf("key%d", i)] = i
			unshardedMap.Unlock()
		}(i)
	}

	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func(i int) {
			wg.Done()
			unshardedMap.RLock()
			_ = unshardedMap.m[fmt.Sprintf("key%d", i)]
			unshardedMap.RUnlock()
		}(i)
	}
	wg.Wait()

	fmt.Println(time.Since(t))
	t = time.Now()
	shardedMap := sharded.NewShardedMap(10)
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func(i int) {
			wg.Done()
			shardedMap.Set(fmt.Sprintf("key%d", i), i)
		}(i)
	}

	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = shardedMap.Get(fmt.Sprintf("key%d", i))
		}(i)
	}
	wg.Wait()
	fmt.Println(time.Since(t))
}
