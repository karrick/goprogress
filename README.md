# goprogress

command line progress bar for Go programs using ANSI escape sequences

# Use

## Percentage

```Go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/karrick/goprogress"
)

const waitKey = false

func main() {
	message := os.Args[len(os.Args)-1]
	cols := 80

	p, err := goprogress.NewPercentage(cols)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}

	for i := 0; i <= 101; i++ {
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
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/karrick/goprogress"
)

const waitKey = false

func main() {
	message := os.Args[len(os.Args)-1]
	cols := 80

	s, err := goprogress.NewSpinner(cols)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}

	for i := 0; i <= 101; i++ {
		s.Update(fmt.Sprintf("%s: %d", message, i))
		s.WriteTo(os.Stdout)
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

# TODO

1. Make compatible with runes that consume more than a single byte.
2. Make compatible with runes that consume more than a single display column.
