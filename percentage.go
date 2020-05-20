package goprogress

import (
	"fmt"
	"io"
)

const prefix = "\033[G\033[7m"
const suffix = "\033[27m"
const lprefix = 7
const lsuffix = 5

type Percentage struct {
	formatted []byte // formatted is the formatted and printable bytes.
	width     int    // width is the number of columns the progress bar should consume.
}

func NewPercentage(width int) (*Percentage, error) {
	if width < 1 {
		return nil, fmt.Errorf("cannot create width less than 1: %d", width)
	}
	return &Percentage{
		formatted: make([]byte, lprefix+width+lsuffix),
		width:     width,
	}, nil
}

func (p *Percentage) appendPercentage(percentage, foo, start int) {
	p.formatted[start] = '%'
	foo--

loop:
	start--
	if foo == 0 {
		copy(p.formatted[start-lsuffix+1:], suffix)
		start -= lsuffix
	}
	foo--

	p.formatted[start] = byte(percentage%10) + '0'

	if percentage /= 10; percentage > 0 {
		goto loop
	}
}

func (p *Percentage) Update(message string, percentage int) {
	if debug > 0 {
		memfill(p.formatted, '?', cap(p.formatted))
	}

	// Start with escape sequences to return to first column and reverse video,
	// and track the number of bytes copied for later on appending to the
	// formatted byte slice.
	idx := copy(p.formatted, prefix)

	// Determine number of columns that should be displayed in reverse video.
	reverseColumns := p.width * percentage / 100
	if reverseColumns > p.width {
		reverseColumns = p.width // handle when given percentage greater than 100
	}

	// Determine number of columns the percentage indication consumes.
	lpercent := 1
	percent := percentage // mutate a copy of the percent
loop:
	lpercent++
	if percent /= 10; percent > 0 {
		goto loop
	}

	// Determine number of columns dedicated for the message and for empty
	// spaces before the percentage indication.
	messageColumns := p.width - lpercent
	var spaceColumns int

	if sc := messageColumns - len(message); sc > 0 {
		spaceColumns = sc
	} else {
		message = message[:messageColumns] // message width exceeds number of columns allotted for it.
	}

	// Determine at which index of the formatted string the percent sign will appear.
	pstart := cap(p.formatted) - 1

	mi := len(message) - reverseColumns
	if mi > 0 {
		idx += copy(p.formatted[idx:], message[:reverseColumns])
		idx += copy(p.formatted[idx:], suffix)
		idx += copy(p.formatted[idx:], message[reverseColumns:])
		memfill(p.formatted[idx:], ' ', spaceColumns)
		p.appendPercentage(percentage, 0, pstart)
	} else if reverseColumns == len(message) {
		idx += copy(p.formatted[idx:], message)
		idx += copy(p.formatted[idx:], suffix)
		memfill(p.formatted[idx:], ' ', spaceColumns)
		p.appendPercentage(percentage, 0, pstart)
	} else if reverseColumns < messageColumns {
		idx += copy(p.formatted[idx:], message)
		mi = -mi
		idx += memfill(p.formatted[idx:], ' ', mi)
		idx += copy(p.formatted[idx:], suffix)
		memfill(p.formatted[idx:], ' ', spaceColumns-mi)
		p.appendPercentage(percentage, 0, pstart)
	} else if reverseColumns == messageColumns {
		idx += copy(p.formatted[idx:], message)
		idx += memfill(p.formatted[idx:], ' ', spaceColumns)
		copy(p.formatted[idx:], suffix)
		p.appendPercentage(percentage, 0, pstart)
	} else if reverseColumns < (messageColumns + lpercent) {
		idx += copy(p.formatted[idx:], message)
		idx += memfill(p.formatted[idx:], ' ', spaceColumns)
		copy(p.formatted[idx:], suffix)
		p.appendPercentage(percentage, lpercent+messageColumns-reverseColumns, pstart)
	} else {
		idx += copy(p.formatted[idx:], message)
		idx += memfill(p.formatted[idx:], ' ', spaceColumns)
		p.appendPercentage(percentage, 0, pstart-lsuffix)
		idx += lpercent
		copy(p.formatted[idx:], suffix)
	}
}

func (p *Percentage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(p.formatted)
	return int64(n), err
}
