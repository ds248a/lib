# lib

## Bytes Buffer

Pool represents byte buffer pool.
Distinct pools may be used for distinct types of byte buffers.
Properly determined byte buffer types with their own pools may help reducing memory waste.

ByteBuffer may be used with functions appending data to the given []byte slice.

```go
package main

import (
	"fmt"
	"github.com/ds247a/pool"
)

func ExampleByteBuffer() {
	bb := pool.Get()

	bb.WriteString("first line\n")
	bb.Write([]byte("second line\n"))
	bb.B = append(bb.B, "third line\n"...)

	fmt.Printf("bytebuffer contents=%q", bb.B)

	// It is safe to release byte buffer now, since it is no longer used.
	pool.Put(bb)
}
```