package csv

import (
	"encoding/csv"
	"github.com/jszwec/csvutil"
	"io"
	"os"
	"strings"
)

// ReaderSeparator defines the type for the separator used in the CSV file.
type ReaderSeparator rune

const (
	// CommaSeparator is used to separate fields with a comma.
	CommaSeparator ReaderSeparator = ','

	// SemicolonSeparator is used to separate fields with a semicolon.
	SemicolonSeparator ReaderSeparator = ';'

	// TabSeparator is used to separate fields with a tab.
	TabSeparator ReaderSeparator = '\t'

	// PipeSeparator is used to separate fields with a pipe.
	PipeSeparator ReaderSeparator = '|'
)

// NewCSVReaderFromPath creates a new CSV Reader from a file path with optional ReaderOptions.
func NewCSVReaderFromPath(path string, optFns ...func(*ReaderOptions)) (Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return NewCSVReader(file, optFns...)
}

// NewCSVReader creates a new CSV Reader from an io.Reader with optional ReaderOptions.
func NewCSVReader(r io.Reader, optFns ...func(*ReaderOptions)) (Reader, error) {
	opt := &ReaderOptions{
		NoHeader:  false,
		Separator: CommaSeparator,
	}

	for _, o := range optFns {
		o(opt)
	}

	localReader := &csvReader{}
	reader := csv.NewReader(r)
	reader.Comma = rune(opt.Separator)

	if err := localReader.loadRows(reader, opt); err != nil {
		return nil, err
	}

	return localReader, nil
}

// ToReaderSeparator converts a string to a ReaderSeparator.
func ToReaderSeparator(separator string) (ReaderSeparator, bool) {
	separator = strings.TrimSpace(separator)
	switch separator {
	case ",":
		return CommaSeparator, true
	case ";":
		return SemicolonSeparator, true
	case "\t":
		return TabSeparator, true
	case "|":
		return PipeSeparator, true
	default:
		return 0, false
	}
}

func decodeObject(csvStr string, obj any) error {
	reader := csv.NewReader(strings.NewReader(csvStr))
	dec, err := csvutil.NewDecoder(reader)
	if err != nil {
		return err
	}

	return dec.Decode(obj)
}
