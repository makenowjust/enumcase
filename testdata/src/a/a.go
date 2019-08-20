package a

import (
	"github.com/MakeNowJust/enumcase/testdata/src/a/b"
)

func main() {
	b1 := b.B1

	switch b1 {
	case b.B1:
	case b.B3:
	}
}
