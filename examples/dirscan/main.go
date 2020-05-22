package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/karrick/godirwalk"
	"github.com/karrick/goprogress"
)

const waitKey = false

func main() {
	cols := flag.Int("columns", 80, "number of columns to use")
	flag.Parse()

	dirname := "."
	if flag.NArg() > 0 {
		dirname = flag.Arg(0)
	}

	//
	// Count the number of file system entries in a hierarchy.
	//
	var totalDirents int

	s, err := goprogress.NewSpinner(*cols)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}

	options := &godirwalk.Options{
		Callback: func(_ string, _ *godirwalk.Dirent) error {
			totalDirents++
			s.Update(fmt.Sprintf("Counting: found %d entries so far...", totalDirents))
			_, err := s.WriteTo(os.Stderr)
			return err
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			totalDirents++
			s.Update(fmt.Sprintf("Counting: %s: %s", osPathname, err))
			s.WriteTo(os.Stderr)
			return godirwalk.SkipNode
		},
	}

	if err = godirwalk.Walk(dirname, options); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
	s.Update(fmt.Sprintf("Counting: found %d entries", totalDirents))
	s.WriteTo(os.Stderr)
	fmt.Println() // newline after spinner progress bar

	//
	// Now present a progress bar with the percentage complete.
	//
	var complete int

	p, err := goprogress.NewPercentage(*cols)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}

	options.Callback = func(osPathname string, _ *godirwalk.Dirent) error {
		complete += 100 // increment by 100 so do not need to multiple progress by 100
		p.Update(fmt.Sprintf("Doing work: %s", osPathname), complete/totalDirents)
		_, err := p.WriteTo(os.Stderr)
		return err
	}
	options.ErrorCallback = func(osPathname string, err error) godirwalk.ErrorAction {
		complete += 100 // increment by 100 so do not need to multiple progress by 100
		p.Update(fmt.Sprintf("Doing work: %s: %s", osPathname, err), complete/totalDirents)
		p.WriteTo(os.Stderr)
		return godirwalk.SkipNode
	}

	if err = godirwalk.Walk(dirname, options); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
	p.Update("Doing work: complete", 100)
	p.WriteTo(os.Stderr)
	fmt.Println() // newline after percentage progress bar
}
