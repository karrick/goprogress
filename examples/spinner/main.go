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

	for i := 0; i < 26; i++ {
		s.Update(fmt.Sprintf("%s: %c", message, byte(i+'a')))
		s.WriteTo(os.Stdout)
		if true {
			time.Sleep(100 * time.Millisecond)
		} else {
			var r rune
			fmt.Scanf("%c", &r)
		}
	}
	s.Update(fmt.Sprintf("%s: complete", message))
	s.WriteTo(os.Stdout)
	fmt.Println() // newline after spinner
}
