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
