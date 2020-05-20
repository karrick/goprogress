package goprogress

import (
	"fmt"
	"io"
)

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
	return &Spinner{
		formatted: make([]byte, lcr+width),
		width:     width,
	}, nil
}

func (p *Spinner) Update(message string) {
	if debug > 0 {
		memfill(p.formatted, '?', cap(p.formatted))
	}

	b := spinner[p.spinner]
	if p.spinner++; p.spinner == lspinner {
		p.spinner = 0
	}

	// Start with escape sequence to return to first column.
	idx := copy(p.formatted, cr)

	// Determine number of columns dedicated for the message and for empty
	// spaces before the spinner.
	messageColumns := p.width - 1

	var spaceColumns int
	if sc := messageColumns - len(message); sc > 0 {
		spaceColumns = sc
	} else {
		message = message[:messageColumns] // message width exceeds number of columns allotted for it.
	}

	idx += copy(p.formatted[idx:], message)
	idx += memfill(p.formatted[idx:], ' ', spaceColumns)
	p.formatted[idx] = b
}

func (p *Spinner) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(p.formatted)
	return int64(n), err
}
