// +build !appengine,!js,!windows,!nacl,!plan9

package terminal

import (
	"io"
	"os"
)

func Check(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return isTerminal(int(v.Fd()))
	default:
		return false
	}
}
