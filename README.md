# lib

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
