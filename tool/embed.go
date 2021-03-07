// +build embed
package tool

import (
	_ "embed"
	"bytes"
	"os"
)

//go:embed build/fotoDen.min.js
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

//go:embed build/default_theme.zip
var defaultThemeZip []byte
var defaultThemeZipLen int64

func defaultThemeZipReader() *bytes.Reader {
	return bytes.NewReader(defaultThemeZip)
}

func init() {
	isEmbed = true
	defaultThemeZipLen = int64(len(defaultThemeZip))
}
