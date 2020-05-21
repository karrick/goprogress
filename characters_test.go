package goprogress

import "testing"

// Characters returns a slice of characters, each character being a slice of
// bytes of the encoded character.
func (s *String) characters() [][]byte {
	// fmt.Fprintf(os.Stderr, "sequence: %v; offsets: %v\n", s.sequence, s.offsets)
	l := len(s.offsets)

	var characters [][]byte
	for i := 0; i < l; i++ {
		// fmt.Fprintf(os.Stderr, "i: %d: %v\n", i, s.Char(i))
		characters = append(characters, s.Char(i))
	}
	return characters
}

func ensureByteSlicesMatch(tb testing.TB, got, want []byte) {
	tb.Helper()

	la, lb := len(got), len(want)

	max := la
	if max < lb {
		max = lb
	}

	for i := 0; i < max; i++ {
		if i < la && i < lb {
			if g, w := got[i], want[i]; g != w {
				tb.Errorf("%d: GOT: %q; WANT: %q", i, got, want)
			}
		} else if i < la {
			tb.Errorf("%d: GOT: extra byte: %q", i, got[i])
		} else /* i < lb */ {
			tb.Errorf("%d: WANT extra byte: %q", i, want[i])
		}
	}
}

func ensureSlicesOfByteSlicesMatch(tb testing.TB, got, want [][]byte) {
	tb.Helper()

	la, lb := len(got), len(want)

	max := la
	if max < lb {
		max = lb
	}

	for i := 0; i < max; i++ {
		if i < la && i < lb {
			ensureByteSlicesMatch(tb, got[i], want[i])
		} else if i < la {
			tb.Errorf("%d: GOT: extra slice: %v", i, got[i])
		} else /* i < lb */ {
			tb.Errorf("%d: WANT: extra slice: %v", i, want[i])
		}
	}
	if tb.Failed() {
		tb.Logf("GOT: %v; WANT: %v", got, want)
	}
}

func TestString(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		got, want := NewString("").characters(), [][]byte(nil)
		ensureSlicesOfByteSlicesMatch(t, got, want)
	})

	t.Run("a", func(t *testing.T) {
		got, want := NewString("a").characters(), [][]byte{[]byte{'a'}}
		ensureSlicesOfByteSlicesMatch(t, got, want)
	})

	t.Run("cafés", func(t *testing.T) {
		got := NewString("cafés").characters()
		want := [][]byte{[]byte{'c'}, []byte{'a'}, []byte{'f'}, []byte{101, 204, 129}, []byte{'s'}}
		ensureSlicesOfByteSlicesMatch(t, got, want)
	})

	t.Run("slice", func(t *testing.T) {
		t.Run("i too large", func(t *testing.T) {
			s := NewString("cafés")

			if got, want := string(s.Slice(6, 13)), ""; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			if got, want := string(s.Slice(6, -1)), ""; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("j too large", func(t *testing.T) {
			s := NewString("cafés")
			if got, want := string(s.Slice(0, 13)), "cafés"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("j is -1", func(t *testing.T) {
			s := NewString("cafés")

			if got, want := string(s.Slice(0, -1)), "cafés"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			if got, want := string(s.Slice(1, -1)), "afés"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			if got, want := string(s.Slice(2, -1)), "fés"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			if got, want := string(s.Slice(3, -1)), "és"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			if got, want := string(s.Slice(4, -1)), "s"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			if got, want := string(s.Slice(5, -1)), ""; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("i and j within range", func(t *testing.T) {
			s := NewString("cafés")

			if got, want := string(s.Slice(0, 5)), "cafés"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			if got, want := string(s.Slice(0, 4)), "café"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			if got, want := string(s.Slice(0, 3)), "caf"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			if got, want := string(s.Slice(0, 2)), "ca"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			if got, want := string(s.Slice(0, 1)), "c"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			if got, want := string(s.Slice(0, 0)), ""; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
	})

	t.Run("trunc", func(t *testing.T) {
		t.Run("index zero", func(t *testing.T) {
			s := NewString("cafés")
			s.Trunc(0)
			if got, want := s.Len(), 0; got != want {
				t.Fatalf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := string(s.sequence), ""; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("index one", func(t *testing.T) {
			s := NewString("cafés")
			s.Trunc(1)
			if got, want := s.Len(), 1; got != want {
				t.Fatalf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := string(s.sequence), "c"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("index two", func(t *testing.T) {
			s := NewString("cafés")
			s.Trunc(2)
			if got, want := s.Len(), 2; got != want {
				t.Fatalf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := string(s.sequence), "ca"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("index before multi-code-point", func(t *testing.T) {
			s := NewString("cafés")
			s.Trunc(3)
			if got, want := s.Len(), 3; got != want {
				t.Fatalf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := string(s.sequence), "caf"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("index after multi-code-point", func(t *testing.T) {
			s := NewString("cafés")
			s.Trunc(4)
			if got, want := s.Len(), 4; got != want {
				t.Fatalf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := string(s.sequence), string([]byte{99, 97, 102, 101, 204, 129}); got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("index equals length", func(t *testing.T) {
			s := NewString("cafés")
			s.Trunc(5)
			if got, want := s.Len(), 5; got != want {
				t.Fatalf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := string(s.sequence), "cafés"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
		t.Run("index greater than length", func(t *testing.T) {
			s := NewString("cafés")
			s.Trunc(6)
			if got, want := s.Len(), 5; got != want {
				t.Fatalf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := string(s.sequence), "cafés"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
	})
}
