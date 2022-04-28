package rand

import (
	"sync"
	"testing"

	"github.com/ds248a/lib/strings"
)

func testRand(t *testing.T) {
	t.Helper()
	values := make([]string, 0)

	for i := 0; i < 10000; i++ {
		n := 32
		dst := make([]byte, n)

		Rand(dst)

		if strings.Include(values, string(dst)) {
			t.Error("Rand() returns the same value")
			return
		}

		values = append(values, string(dst))

		for i := range dst {
			if string(rune(i)) == "" {
				t.Errorf("RandBytes() invalid char '%v'", dst[i])
			}
		}

		if len(dst) != n {
			t.Errorf("RandBytes() length '%d', want '%d'", len(dst), n)
		}
	}
}

func Test_Rand(t *testing.T) {
	testRand(t)
}

func Test_RandConcurrent(t *testing.T) {
	n := 32
	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testRand(t)
		}()
	}

	wg.Wait()
}

func Benchmark_Rand(b *testing.B) {
	n := 32
	dst := make([]byte, n)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Rand(dst)
		}
	})
}
