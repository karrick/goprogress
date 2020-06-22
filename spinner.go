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

// NewSpinner returns a progress bar with specified width, to be used when the
// percentage complete is not known for every update.
//
//     func main() {
//         cols := flag.Int("columns", 80, "number of columns to use")
//         flag.Parse()
//
//         s, err := goprogress.NewSpinner(*cols)
//         if err != nil {
//             fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
//             os.Exit(1)
//         }
//
//         message := flag.Arg(flag.NArg() - 1)
//
//         for i := 0; i <= 42; i++ {
//             s.Update(fmt.Sprintf("%s: %d", message, i))
//             s.WriteTo(os.Stdout)
//             time.Sleep(10 * time.Millisecond)
//         }
//         s.Update(fmt.Sprintf("%s: complete", message))
//         s.WriteTo(os.Stdout)
//         fmt.Println() // newline after spinner
//     }
func NewSpinner(width int) (*Spinner, error) {
	if width < 1 {
		return nil, fmt.Errorf("cannot create width less than 1: %d", width)
	}
	return &Spinner{width: width}, nil
}

// Update will update the Spinner progress bar with the provided message and
// update the spinner character.
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

// WriteTo will send the sequence of ANSI characters required to redraw the
// Spinner progress bar to the specified io.Writer.
func (p *Spinner) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(p.formatted)
	return int64(n), err
}
