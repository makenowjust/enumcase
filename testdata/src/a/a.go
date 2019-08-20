package a

import (
	"a/b"
)

func main() {
	b1 := b.B1

	switch b1 { // want "missing case\\(s\\) to a/b.B value\\(s\\): B2"
	case b.B1:
	case b.B3:
	}
}
