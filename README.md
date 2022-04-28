# lib

### Strconv
String conversions
```go
func B2S(b []byte) string
func S2B(s string) (b []byte)
```

Benchmark
```go
Benchmark_B2S-4   	1000000000	         0.5577 ns/op	       0 B/op	       0 allocs/op
Benchmark_S2B-4   	1000000000	         0.5477 ns/op	       0 B/op	       0 allocs/o
```


### Random
Generating a random set of characters of the latin alphabet
```go
func Rand(dst []byte) []byte
```

Benchmark
```go
Benchmark_Rand-4   	 10000000	      790 ns/op	       0 B/op	       0 allocs/op
```


### BPool
Represents byte buffer pool.
Distinct pools may be used for distinct types of byte buffers.
Properly determined byte buffer types with their own pools may help reducing memory waste.

```go
package main

import (
	"fmt"
	"github.com/ds248a/lib/bpool"
)

func main() {
	b := bpool.Get()

	b.WriteString("line 1\n")
	b.Write([]byte("line 2\n"))
	b.B = append(b.B, "line 3\n"...)
	
	fmt.Printf("b.B=%q", b.B)
	bpool.Put(b)
}
```

Benchmark
```go
BenchmarkBBpoolBuf-4   	11564948	  99.83 ns/op	   0 B/op	   0 allocs/op
```
