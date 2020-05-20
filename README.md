# goprogress

command line progress bar for Go programs using ANSI escape sequences

# Use

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

const waitKey = false

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
		if waitKey {
			var r rune
			fmt.Scanf("%c", &r)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
	fmt.Println()
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

# TODO

1. Make compatible with runes that consume more than a single byte,
   with single runes that consume more than one display column, and
   with multiple runes that consume only a single display column.
   https://blog.golang.org/normalization
