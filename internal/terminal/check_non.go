// +build js nacl plan9

package terminal

import (
	"io"
)

func Check(w io.Writer) bool {
	return false
}
