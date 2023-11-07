package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/karrick/goprogress"
	"github.com/karrick/gows"
	"golang.org/x/sys/unix"
)

func main() {
	flag.Parse()

	message := flag.Arg(flag.NArg() - 1)

	width, _, err := gows.GetWinSize()
	if err != nil {
		bail(1, err)
	}

	spinner, err := goprogress.NewSpinner(width)
	if err != nil {
		bail(1, err)
	}

	// Spin off a goroutine to receive SIGWINCH signal, and when one arrives,
	// fetch the terminal width and update the spinner width.
	var wg sync.WaitGroup
	done := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()

		signals := make(chan os.Signal, 2)
		signal.Notify(signals, unix.SIGWINCH)

		for {
			select {
			case <-done:
				return
			case <-signals:
				width, _, err := gows.GetWinSize()
				if err != nil {
					bail(1, err)
				}
				spinner.Width(width)
			}
		}
	}()

	var i int

	for {
		spinner.Update(fmt.Sprintf("%s: %c", message, byte(i+'a')))
		spinner.WriteTo(os.Stdout)
		if true {
			time.Sleep(100 * time.Millisecond)
		} else {
			var r rune
			fmt.Scanf("%c", &r)
		}
		i = (i + 1) % 26
	}

	close(done)
	wg.Wait()

	spinner.Update(fmt.Sprintf("%s: complete", message))
	spinner.WriteTo(os.Stdout)
	fmt.Println() // newline after spinner
}

func bail(status int, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
	os.Exit(status)
}
