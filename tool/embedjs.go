// +build embedjs
package tool

import (
	_ "embed"
	"os"
)

//go:embed fotoDen.min.js
var fotoDenJS []byte

func writefotoDenJS(j string) error {
	f, err := os.Create(j)
	if checkError(err) {
		return err
	}

	_, err = f.Write(fotoDenJS)
	if checkError(err) {
		return err
	}

	return nil
}

func init() {
	isEmbed = true
}
