package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/karrick/goprogress"
)

const debug = 1

func main() {
	message := os.Args[len(os.Args)-1]
	_ = message

	cols := 100

	p, err := goprogress.NewProgress(cols)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
	_ = p

	if debug > 0 {
		fmt.Println(strings.Repeat("-", cols))
	}

	const waitKey = false

	for i := 0; i <= 100; i++ {
		p.UpdatePercentage(message, i)
		p.Display(os.Stdout)
		if waitKey {
			var r rune
			fmt.Scanf("%c", &r)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
	fmt.Println()

	for i := 0; i <= 100; i++ {
		p.UpdateSpinner(fmt.Sprintf("Doing some other stuff: %d", i))
		p.Display(os.Stdout)
		if waitKey {
			var r rune
			fmt.Scanf("%c", &r)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
	fmt.Println()

	for i := 0; i <= 100; i++ {
		p.UpdatePercentage(message, i)
		p.Display(os.Stdout)
		if waitKey {
			var r rune
			fmt.Scanf("%c", &r)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
	fmt.Println()
}
