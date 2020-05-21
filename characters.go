package goprogress

import "golang.org/x/text/unicode/norm"

type String struct {
	sequence []byte // all bytes from the string
	offsets  []int  // offsets[n] is the offset for the start of character N in sequence.
}

func NewString(s string) *String {
	sequence := make([]byte, 0, len(s)) // pre-allocate byte slice at least as long as s
	var offsets []int
	var offset int

	var ni norm.Iter
	ni.InitString(norm.NFKD, s)

	for !ni.Done() {
		// Initial testing revealed that norm.Iter reuses the same byte slice
		// each time, so will need to copy the bytes from the returned slice
		// each time through the loop into our own sequence.
		b := ni.Next()
		sequence = append(sequence, b...)
		offsets = append(offsets, offset)
		offset += len(b)
	}
	return &String{sequence: sequence, offsets: offsets}
}

// Bytes returns the entire slice of bytes that encode all characters.
func (s *String) Bytes() []byte {
	return s.sequence
}

// Char returns the slice of bytes that encode the ith character.
func (s *String) Char(i int) []byte {
	if i < len(s.offsets)-1 {
		return s.sequence[s.offsets[i]:s.offsets[i+1]]
	}
	return s.sequence[s.offsets[i]:]
}

// Len returns the number of display-able characters in String.
func (s *String) Len() int {
	return len(s.offsets)
}

// Slice returns the slice of bytes that encode the ith thru jth-1
// characters. As a special case, when j is -1, this returns from the ith
// character to the end of the string.
func (s *String) Slice(i, j int) []byte {
	if i >= len(s.offsets) {
		return nil
	}
	istart := s.offsets[i]

	if j == -1 || j >= len(s.offsets) {
		return s.sequence[istart:]
	}

	return s.sequence[istart:s.offsets[j]]
}

// Trunc truncates String to max of i characters. No operation is performed when
// i is greater than the number of characters in String.
func (s *String) Trunc(i int) {
	if i >= len(s.offsets) {
		return
	}
	s.sequence = s.sequence[:s.offsets[i]]
	s.offsets = s.offsets[:i]
}
