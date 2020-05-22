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
	arrow := flag.Bool("arrows", false, "use arrows to control")
	flag.Parse()

	p, err := goprogress.NewPercentage(*cols)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}

	message := flag.Arg(flag.NArg() - 1)

	var i int

	for {
		p.Update(message, i)
		p.WriteTo(os.Stdout)
		if *arrow {
			b := make([]byte, 1)
			if _, err = os.Stdin.Read(b); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
				os.Exit(1)
			}
			switch b[0] {
			case 'h':
				i--
			case 'l':
				i++
			case 'q':
				goto end
			default:
				fmt.Fprintf(os.Stderr, "b: %q\n", b)
			}
		} else {
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
end:
	fmt.Println() // newline after spinner
}
