package config

import (
	"fmt"
	"io"
	"io/fs"
	"testing"
)

var errOpen = fmt.Errorf("Open is broken!")
var errMkdirAll = fmt.Errorf("MkdirAll is broken!")
var errReadAll = fmt.Errorf("ReadAll is broken!")

type fakeFS struct {
	openError     bool
	mkdirAllError bool
	readAllError  bool
	contents      []byte
}

type fakeFile struct {
	file
}

func (f fakeFile) Close() error {
	return nil
}

func (f fakeFS) Open(path string) (file, error) {
	if f.openError {
		return nil, errOpen
	}
	return fakeFile{}, nil
}
func (f fakeFS) MkdirAll(path string, perm fs.FileMode) error {
	if f.mkdirAllError {
		return errMkdirAll
	}
	return nil
}

func (f fakeFS) ReadAll(r io.Reader) ([]byte, error) {
	if f.readAllError {
		return nil, errReadAll
	}
	return f.contents, nil
}

func TestParse(t *testing.T) {

	var valid_toml = `location="/tmp"

[[podcasts]]
name="Philosophize This!"
url="https://philosophizethis.libsyn.com/rss"

[[podcasts]]
name="Lex Fridman Podcast"
url="https://lexfridman.com/feed/podcast"
`

	var tests = []struct {
		file         string
		openError    bool
		readAllError bool
		want         error
	}{
		{valid_toml, true, false, errOpen},
		{valid_toml, false, true, errReadAll},
		{valid_toml, true, true, errOpen},
	}

	for _, test := range tests {
		// Create mock OS for testing
		fsym = fakeFS{
			test.openError,
			false,
			test.readAllError,
			[]byte(test.file),
		}
		if _, got := Parse("fake path"); got != test.want {
			t.Errorf("Parse() = %v, wanted: %v", got, test.want)
		}
	}
}

func TestSetup(t *testing.T) {
	fsym = fakeFS{
		false,
		false,
		false,
		[]byte(``),
	}

	config, err := Parse("fake path")
	if err != nil {
		t.Errorf("Error constructing TestSetup initialization")
	}

	var tests = []struct {
		mkdirError bool
		want       error
	}{
		{false, nil},
		{true, errMkdirAll},
	}

	for _, test := range tests {
		// Create mock OS for testing
		fsym = fakeFS{
			false,
			test.mkdirError,
			false,
			[]byte(nil),
		}
		if got := config.Setup(); got != test.want {
			t.Errorf("Parse() = %v, wanted: %v", got, test.want)
		}
	}
}
