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
