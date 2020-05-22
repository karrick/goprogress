package goprogress

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/karrick/goutfs"
)

const debugS = 0
const cr = "\033[G"
const spinner = "-\\|/"
const lcr = 3
const lspinner = 4

type Spinner struct {
	formatted []byte // formatted is the formatted and printable bytes.
	spinner   int    // spinner is an index into the spinner string.
	width     int    // width is the number of columns the progress bar should consume.
}

func NewSpinner(width int) (*Spinner, error) {
	if width < 1 {
		return nil, fmt.Errorf("cannot create width less than 1: %d", width)
	}
	return &Spinner{width: width}, nil
}

func (p *Spinner) Update(message string) {
	// Determine number of columns dedicated for the message and for empty
	// spaces before the spinner.
	messageColumns := p.width - 1
	var spaceColumns int

	ms := goutfs.NewString(message)
	lms := ms.Len()

	if sc := messageColumns - lms; sc >= 0 {
		spaceColumns = sc
	} else {
		// Truncate message string to the number of characters allotted for it.
		ms.Trunc(messageColumns)
		lms = ms.Len()
	}

	// After potentially resizing message string and calculating number of
	// columns allotted for spaces, calculate the size of the formatted byte
	// slice.
	if required := lcr + len(ms.Bytes()) + spaceColumns + 1; cap(p.formatted) < required {
		p.formatted = make([]byte, required) // grow
	} else if cap(p.formatted) > required {
		p.formatted = p.formatted[:required] // trim
	}
	if debugS > 0 {
		memfill(p.formatted, '?', cap(p.formatted))
		fmt.Fprintf(os.Stderr, "\n%s\n", strings.Repeat("-", p.width))
	}

	// Start with escape sequence to return to first column.
	idx := copy(p.formatted, cr)

	// Append message
	idx += copy(p.formatted[idx:], ms.Bytes())

	// Append spaces
	idx += memfill(p.formatted[idx:], ' ', spaceColumns)

	// Then finally the spinner
	p.formatted[idx] = spinner[p.spinner]
	if p.spinner++; p.spinner == lspinner {
		p.spinner = 0
	}
}

func (p *Spinner) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(p.formatted)
	return int64(n), err
}
