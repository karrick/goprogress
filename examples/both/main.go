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
	cols := 80

	p, err := goprogress.NewPercentage(cols)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
	s, err := goprogress.NewSpinner(cols)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}

	if debug > 0 {
		fmt.Println(strings.Repeat("-", cols))
	}

	const waitKey = false

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

	for i := 0; i <= 100; i++ {
		s.Update(fmt.Sprintf("Doing some other stuff: %d", i))
		s.WriteTo(os.Stdout)
		if waitKey {
			var r rune
			fmt.Scanf("%c", &r)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
	fmt.Println()

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
