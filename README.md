# goprogress

command line progress bar for Go programs using ANSI escape sequences

# Use

## Percentage

```Go
package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/karrick/goprogress"
)

func main() {
    cols := flag.Int("columns", 80, "number of columns to use")
    flag.Parse()

    p, err := goprogress.NewPercentage(*cols)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
        os.Exit(1)
    }

    message := flag.Arg(flag.NArg() - 1)

    for i := 0; i <= 100; i++ {
        p.Update(message, i)
        p.WriteTo(os.Stdout)
        time.Sleep(10 * time.Millisecond)
    }
    fmt.Println() // newline after spinner
}
```

## Spinner

```Go
package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/karrick/goprogress"
)

func main() {
    cols := flag.Int("columns", 80, "number of columns to use")
    flag.Parse()

    s, err := goprogress.NewSpinner(*cols)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
        os.Exit(1)
    }

    message := flag.Arg(flag.NArg() - 1)

    for i := 0; i <= 42; i++ {
        s.Update(fmt.Sprintf("%s: %d", message, i))
        s.WriteTo(os.Stdout)
        time.Sleep(10 * time.Millisecond)
    }
    s.Update(fmt.Sprintf("%s: complete", message))
    s.WriteTo(os.Stdout)
    fmt.Println() // newline after spinner
}
```

# TODO

Make fully UTF-8 compatible.

## Definitions

For the purposes of this library, I have attempted to adopt the
universal and Go specific terminology for characters, code points,
runes, and bytes. There is a chance that I misread a resource and have
an error in my terminology, but a best effort has been attempted.

### Character

Each character occupies a single column in the output, and roughly
corresponds to what a human sees when they look at the printed text. A
human might see the latin letter e with an accent grave over it, for
example.

Characters are stored and transmitted using some encoding. In unicode
those encodings are called code points. Because of how combining
characters work in unicode, some characters could have multiple code
point representations. For instance, the lower case letter e with an
accent grave could be encoded as a single unicode code point, or
alternatively encoded by two code points: the first one being the
lower case latin letter e, the second as what is known as a combining
code point, in this case the combining code point for accent
grave. Both of these representations result in the same character
being displayed, but have two byte encodings. There are libraries to
normalize these encodings to one of various canonical
standards. However, I am not certain character normalization needs to
be addressed in this library.

### Code Point, a.k.a. Go rune

A code point is called a rune in Go parlance. A Go rune is stored as
an int32 value. Remember a rune is not necessarily a single
character. Some characters have multiple unicode encodings, each of
which could be single or multiple code points.

Another point--no pun intended--is there are look alike characters in
unicode. Not just different code point sequences that represent the
same character, but two different characters that happen to look
alike. For instance, the latin capitol K looks identical to the
unicode code point for the Kelvin symbol. This library need not worry
itself with look alike characters. In order to function correctly,
this library merely needs to know at one byte offset a particular
character ends and the next character begins.

### Strings

Go has no restrictions on the sequence of bytes stored in a
string. The only restriction Go puts on the bytes in a string are that
Go source code is defined as UTF-8, which means most string literal
values are valid UTF-8 encodings. This is not always the case,
however, as Go allows byte level escapes to be included in string
literals, which may or may not represent valid UTF-8 encoded data.

Iterating over a UTF-8 string will result in some runes that require
multiple bytes, and other runes that require a single byte.

### Starting vs Non-Staring (Combining) Rune

Unicode defines many code points that are called starting code
points. They may be displayed independently of any other code point,
and the may be modified indefinitely by appending non-starting code
points. These non-starting code points are more frequently called
combining code points in literature.

## Plan

I need to figure out how to iterate over a string on a character by
character basis. Go provides the ability to iterate over a string rune
by rune, but some runes are combining code points to be applied to a
previous starting rune. The first step will be to iterate over a
sequence of code points (Go runes) in a string, and checking each one
to determine whether it is a starting or a combining code point. This
will allow the library to count the display width of a string.

This logic will need to be used to iterate through a string on a
character by character basis, so the library does not append the
reverse video reset byte sequence before any combining code points for
a character.

Furthermore, when truncating long strings to fit in a narrow progress
bar, it is important to truncate based on character width rather than
truncating on bytes or even code points.

One precaution to keep in mind is the unicode/utf-8 function
DecodeRune will return an error when the rune is not the shortest
possible UTF-8 encoding for the value. If this function is needed,
then using a code point normalization library, such as one referenced
below, may be required when accepting a string from a client, to
normalize the UTF-8 string first, then allow normal enumeration over
the string's code points.

My first attempt will be in repeatedly using the norm.FirstBoundary
function to find the next character in a string.

## References

https://blog.golang.org/strings
https://blog.golang.org/normalization
https://pkg.go.dev/golang.org/x/text/transform?tab=doc
