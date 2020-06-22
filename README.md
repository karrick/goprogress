# goprogress

Progress bar for command-line Go programs using ANSI escape
sequences. Respects UTF-8 character boundries for characters that are
encoded with multiple unicode code points.

# Use

Documentation is available via
[![GoDoc](https://godoc.org/github.com/karrick/goprogress?status.svg)](https://godoc.org/github.com/karrick/goprogress)
and
[https://pkg.go.dev/github.com/karrick/goprogress?tab=doc](https://pkg.go.dev/github.com/karrick/goprogress?tab=doc).

## Percentage

```Go
package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/karrick/goprogress"
)

func main() {
    cols := flag.Int("columns", 80, "number of columns to use")
    flag.Parse()

    p, err := goprogress.NewPercentage(*cols)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
        os.Exit(1)
    }

    message := flag.Arg(flag.NArg() - 1)

    for i := 0; i <= 100; i++ {
        p.Update(message, i)
        p.WriteTo(os.Stdout)
        time.Sleep(10 * time.Millisecond)
    }
    fmt.Println() // newline after spinner
}
```

## Spinner

```Go
package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/karrick/goprogress"
)

func main() {
    cols := flag.Int("columns", 80, "number of columns to use")
    flag.Parse()

    s, err := goprogress.NewSpinner(*cols)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
        os.Exit(1)
    }

    message := flag.Arg(flag.NArg() - 1)

    for i := 0; i <= 42; i++ {
        s.Update(fmt.Sprintf("%s: %d", message, i))
        s.WriteTo(os.Stdout)
        time.Sleep(10 * time.Millisecond)
    }
    s.Update(fmt.Sprintf("%s: complete", message))
    s.WriteTo(os.Stdout)
    fmt.Println() // newline after spinner
}
```
