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
