package sharded

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkShardedMap(b *testing.B) {
	for _, shards := range []int{1, 4, 16, 64, 256} {
		for _, concurrency := range []int{1, 16, 256} {
			runTest(b, shards, concurrency)
		}
	}
}

func runTest(b *testing.B, shards, concurrency int) {
	shardedMap := NewShardedMap(shards)
	b.SetParallelism(concurrency)
	b.Run(fmt.Sprintf(`sharded-%d-%d`, shards, concurrency), func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			o := rand.Intn(1e6)
			for i := 0; pb.Next(); i++ {
				shardedMap.Set(fmt.Sprintf("key%d", i+o), i)
			}
		})
	})
}
