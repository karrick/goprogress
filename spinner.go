package goprogress

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"

	"github.com/karrick/goutfs"
)

const debugS = 0
const cr = "\033[G"
const spinner = "-\\|/"
const lcr = 3
const lspinner = 4

type Spinner struct {
	formatted []byte // formatted is the formatted and printable bytes.
	spinner   int32  // spinner is an index into the spinner string.
	width     int32  // width is the number of columns the progress bar should consume.
}

// NewSpinner returns a progress bar with specified width, to be used when the
// percentage complete is not known for every update.
//
//	func main() {
//	    cols := flag.Int("columns", 80, "number of columns to use")
//	    flag.Parse()
//
//	    s, err := goprogress.NewSpinner(*cols)
//	    if err != nil {
//	        fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
//	        os.Exit(1)
//	    }
//
//	    message := flag.Arg(flag.NArg() - 1)
//
//	    for i := 0; i <= 42; i++ {
//	        s.Update(fmt.Sprintf("%s: %d", message, i))
//	        s.WriteTo(os.Stdout)
//	        time.Sleep(10 * time.Millisecond)
//	    }
//	    s.Update(fmt.Sprintf("%s: complete", message))
//	    s.WriteTo(os.Stdout)
//	    fmt.Println() // newline after spinner
//	}
func NewSpinner(width int) (*Spinner, error) {
	if width < 1 {
		return nil, fmt.Errorf("cannot create width less than 1: %d", width)
	}
	return &Spinner{width: int32(width)}, nil
}

// Update will update the Spinner progress bar with the provided message and
// update the spinner character.
func (p *Spinner) Update(message string) {
	// Determine number of columns dedicated for the message and for empty
	// spaces before the spinner.
	width := int(atomic.LoadInt32(&p.width))

	messageColumns := width - 1 // subtract one for width of spinner character

	ms := goutfs.NewString(message)
	lms := ms.Len()

	spaceColumns := messageColumns - lms
	if spaceColumns < 0 {
		// Truncate message string to the number of characters allotted for it.
		ms.Trunc(messageColumns)
		lms = ms.Len()
		spaceColumns = 0
	}

	// After potentially resizing message string and calculating number of
	// columns allotted for spaces, calculate the size of the formatted byte
	// slice.
	messageBytes := ms.Bytes()
	required := lcr + len(messageBytes) + spaceColumns + 1

	if required > len(p.formatted) {
		p.formatted = make([]byte, required) // grow
	} else if required < len(p.formatted) {
		p.formatted = p.formatted[:required] // trim
	}

	if debugS > 0 {
		memfill(p.formatted, '?', len(p.formatted))
		fmt.Fprintf(os.Stderr, "\n%s\n", strings.Repeat("-", width))
	}

	// Start with escape sequence to return to first column.
	idx := copy(p.formatted, cr)

	// Append message.
	idx += copy(p.formatted[idx:], messageBytes)

	// Append any spaces which may be required.
	idx += memfill(p.formatted[idx:], ' ', spaceColumns)

	// Then finally the spinner.
	p.formatted[idx] = spinner[p.spinner]
	if p.spinner++; p.spinner == lspinner {
		p.spinner = 0
	}
}

// Width updates the number of columns allotted for the spinner to use.
func (p *Spinner) Width(width int) {
	atomic.StoreInt32(&p.width, int32(width))
}

// WriteTo will send the sequence of ANSI characters required to redraw the
// Spinner progress bar to the specified io.Writer.
func (p *Spinner) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(p.formatted)
	return int64(n), err
}
