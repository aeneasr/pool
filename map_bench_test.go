package pool

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func sliceGenerator(size uint) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	result := make([]int, size)

	for i := uint(0); i < size; i++ {
		result[i] = int(r.Int63())
	}

	return result
}

var bgctx = context.Background()

func BenchmarkMap10000x4(b *testing.B) {
	var err error
	in := sliceGenerator(10000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err = Map(bgctx, in, func(ctx context.Context, t int, i int) (int, error) {
			return t * 2, nil
		}, WithWorkers(1))
		if err != nil {
			b.Error(err)
		}
	}
}
