package goprogress

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/karrick/goutfs"
)

const debugP = 0
const prefix = "\033[G\033[7m"
const suffix = "\033[27m"
const lprefix = 7
const lsuffix = 5

type Percentage struct {
	formatted []byte // formatted is the formatted and printable bytes.
	width     int    // width is the number of columns the progress bar should consume.
}

// NewPercentage returns a progress bar with specified width that will include
// an indication of the percentage complete every time it is updated.
//
//     func main() {
//         cols := flag.Int("columns", 80, "number of columns to use")
//         flag.Parse()
//
//         p, err := goprogress.NewPercentage(*cols)
//         if err != nil {
//             fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
//             os.Exit(1)
//         }
//
//         message := flag.Arg(flag.NArg() - 1)
//
//         for i := 0; i <= 100; i++ {
//             p.Update(message, i)
//             p.WriteTo(os.Stdout)
//             time.Sleep(10 * time.Millisecond)
//         }
//         fmt.Println() // newline after spinner
//     }
func NewPercentage(width int) (*Percentage, error) {
	if width < 4 {
		return nil, fmt.Errorf("cannot create width less than 4: %d", width)
	}
	return &Percentage{width: width}, nil
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

// Update will update the Percentage progress bar with the provided message and
// update the percentage complete.
func (p *Percentage) Update(message string, percentage int) {
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

	ms := goutfs.NewString(message)
	lms := ms.Len()

	// fmt.Fprintf(os.Stderr, "\np.width: %d; percentage: %d; lpercent: %d; message columns: %d; space columns: %d\n", p.width, percentage, lpercent, messageColumns, messageColumns-lms)

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
	if required := lprefix + len(ms.Bytes()) + lsuffix + spaceColumns + lpercent; cap(p.formatted) < required {
		p.formatted = make([]byte, required) // grow
	} else if cap(p.formatted) > required {
		p.formatted = p.formatted[:required] // trim
	}
	if debugP > 0 {
		memfill(p.formatted, '?', cap(p.formatted))
		fmt.Fprintf(os.Stderr, "%s\n", strings.Repeat("-", p.width))
	}

	// Start with escape sequences to return to first column and reverse video,
	// and track the number of bytes copied for later on appending to the
	// formatted byte slice.
	idx := copy(p.formatted, prefix)

	mi := lms - reverseColumns
	if mi > 0 {
		idx += copy(p.formatted[idx:], ms.Slice(0, reverseColumns))
		idx += copy(p.formatted[idx:], suffix)
		idx += copy(p.formatted[idx:], ms.Slice(reverseColumns, -1))
		idx += memfill(p.formatted[idx:], ' ', spaceColumns)
		p.appendPercentage(percentage, 0, idx+lpercent-1)
	} else if reverseColumns == lms {
		idx += copy(p.formatted[idx:], ms.Bytes())
		idx += copy(p.formatted[idx:], suffix)
		idx += memfill(p.formatted[idx:], ' ', spaceColumns)
		p.appendPercentage(percentage, 0, idx+lpercent-1)
	} else if reverseColumns < messageColumns {
		idx += copy(p.formatted[idx:], ms.Bytes())
		mi = -mi
		idx += memfill(p.formatted[idx:], ' ', mi)
		idx += copy(p.formatted[idx:], suffix)
		idx += memfill(p.formatted[idx:], ' ', spaceColumns-mi)
		p.appendPercentage(percentage, 0, idx+lpercent-1)
	} else if reverseColumns == messageColumns {
		idx += copy(p.formatted[idx:], ms.Bytes())
		idx += memfill(p.formatted[idx:], ' ', spaceColumns)
		idx += copy(p.formatted[idx:], suffix)
		p.appendPercentage(percentage, 0, idx+lpercent-1)
	} else if reverseColumns < (messageColumns + lpercent) {
		idx += copy(p.formatted[idx:], ms.Bytes())
		idx += memfill(p.formatted[idx:], ' ', spaceColumns)
		idx += copy(p.formatted[idx:], suffix)
		p.appendPercentage(percentage, lpercent+messageColumns-reverseColumns, idx+lpercent-1)
	} else {
		idx += copy(p.formatted[idx:], ms.Bytes())
		idx += memfill(p.formatted[idx:], ' ', spaceColumns)
		p.appendPercentage(percentage, 0, idx+lpercent-1)
		idx += lpercent
		idx += copy(p.formatted[idx:], suffix)
	}
}

// WriteTo will send the sequence of ANSI characters required to redraw the
// Percentage progress bar to the specified io.Writer.
func (p *Percentage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(p.formatted)
	return int64(n), err
}
