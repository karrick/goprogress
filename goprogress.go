package goprogress

import (
	"fmt"
	"io"
	"os"
)

type Progress struct {
	Formatted            []byte
	size, width, spinner int
}

const cr = "\033[G"
const prefix = "\033[G\033[7m"
const suffix = "\033[27m"
const spinner = "-\\|/"

const lcr = 3
const lprefix = 7
const lsuffix = 5
const lspinner = 4

func NewProgress(width int) (*Progress, error) {
	if width < 1 {
		return nil, fmt.Errorf("cannot create width less than 1: %d", width)
	}
	size := width + lprefix + lsuffix + 1

	p := &Progress{
		Formatted: make([]byte, size),
		size:      size,
		width:     width,
	}
	return p, nil
}

type flusher interface {
	Flush() error
}

func (p *Progress) Display(w io.Writer) error {
	if _, err := w.Write(p.Formatted); err != nil {
		return err
	}
	if f, ok := w.(flusher); ok {
		return f.Flush()
	}
	return nil
}

func (p *Progress) appendPercentage(percentage, foo, start int) {
	p.Formatted[start] = '%'
	foo--
	pc := percentage

loop:
	start--
	if foo == 0 {
		copy(p.Formatted[start-lsuffix+1:], suffix)
		start -= lsuffix
	}
	foo--

	p.Formatted[start] = byte(pc%10) + '0'

	if pc /= 10; pc > 0 {
		goto loop
	}
}

const debug = 1

func (p *Progress) UpdatePercentage(message string, percentage int) {
	p.Formatted = p.Formatted[:p.size]
	if debug > 0 {
		memset(p.Formatted, '?', p.size-1)
	}

	copy(p.Formatted, prefix)

	lpc := 1
	mcap := p.width - 1
	pc := percentage

loop:
	lpc++
	mcap--
	if pc /= 10; pc > 0 {
		goto loop
	}

	var rc int
	if percentage < 100 {
		rc = p.width * percentage / 100
	} else {
		rc = p.width
	}

	lmessage := len(message)

	lblanks := mcap - lmessage
	if lblanks < 0 {
		fmt.Fprintf(os.Stderr, "TODO: trim message too long\n")
		os.Exit(1)
	}

	pstart := p.size - 2

	mi := lmessage - rc
	if mi > 0 {
		copy(p.Formatted[lprefix:], message[:rc])
		copy(p.Formatted[lprefix+rc:], suffix)
		copy(p.Formatted[lprefix+rc+lsuffix:], message[rc:])
		memset(p.Formatted[lprefix+lmessage+lsuffix:], ' ', lblanks)
		p.appendPercentage(percentage, -1, pstart)
	} else if rc == lmessage {
		copy(p.Formatted[lprefix:], message)
		copy(p.Formatted[lprefix+lmessage:], suffix)
		memset(p.Formatted[lprefix+lmessage+lsuffix:], ' ', lblanks)
		p.appendPercentage(percentage, -1, pstart)
	} else if rc < mcap {
		copy(p.Formatted[lprefix:], message)
		mi = -mi
		memset(p.Formatted[lprefix+lmessage:], ' ', mi)
		copy(p.Formatted[lprefix+lmessage+mi:], suffix)
		memset(p.Formatted[lprefix+lmessage+mi+lsuffix:], ' ', lblanks-mi)
		p.appendPercentage(percentage, -1, pstart)
	} else if rc == mcap {
		copy(p.Formatted[lprefix:], message)
		memset(p.Formatted[lprefix+lmessage:], ' ', lblanks)
		copy(p.Formatted[lprefix+lmessage+lblanks:], suffix)
		p.appendPercentage(percentage, -1, pstart)
	} else if rc < (mcap + lpc) {
		foo := lpc - rc + mcap
		copy(p.Formatted[lprefix:], message)
		memset(p.Formatted[lprefix+lmessage:], ' ', lblanks)
		copy(p.Formatted[lprefix+lmessage+lblanks:], suffix)
		p.appendPercentage(percentage, foo, pstart)
	} else {
		copy(p.Formatted[lprefix:], message)
		memset(p.Formatted[lprefix+lmessage:], ' ', lblanks)
		p.appendPercentage(percentage, -1, p.size-lsuffix-2)
		copy(p.Formatted[lprefix+lmessage+lblanks+lpc:], suffix)
	}
}

func (p *Progress) UpdateSpinner(message string) {
	c := spinner[p.spinner]
	if p.spinner++; p.spinner == lspinner {
		p.spinner = 0
	}

	copy(p.Formatted, cr)

	lmessage := len(message)
	copy(p.Formatted[lcr:], message)

	mcap := p.width - 1

	lblanks := mcap - lmessage
	if lblanks < 0 {
		fmt.Fprintf(os.Stderr, "TODO: trim message too long\n")
		os.Exit(1)
	}

	memset(p.Formatted[lcr+lmessage:], ' ', lblanks)
	p.Formatted[lcr+lmessage+lblanks] = c
	p.Formatted = p.Formatted[:lcr+lmessage+lblanks+1]
}

func memset(buf []byte, b byte, l int) {
	for i := 0; i < l; i++ {
		buf[i] = b
	}
}
