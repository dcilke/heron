# heron

A JSON stream parser.

Heron provides a convenience wrapper when dealing with JSON streams. Given a stream, it will read from the stream and emit any valid JSON objects or arrays.

# Examples

```go
package main

import (
  "fmt"
  "os"

  "github.com/dcilke/heron"
)

func main() {
  h := heron.New(
    heron.WithJSON(func(a any) {
      fmt.Printf("JSON: %v\n", a)
    }),
  )

  h.Process(os.Stdin)
}
```
