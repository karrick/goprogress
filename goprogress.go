package goprogress

const debug = 0

// memfill byte behaves exactly like memset, but returns the number of bytes.
func memfill(buf []byte, b byte, l int) int {
	for i := 0; i < l; i++ {
		buf[i] = b
	}
	return l
}
