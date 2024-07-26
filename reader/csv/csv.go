package csv

import (
	"encoding/csv"
	"github.com/jszwec/csvutil"
	"io"
	"os"
	"strings"
)

type ReaderSeparator rune

const (
	CommaSeparator     ReaderSeparator = ','
	SemicolonSeparator ReaderSeparator = ';'
	TabSeparator       ReaderSeparator = '\t'
	PipeSeparator      ReaderSeparator = '|'
)

func NewCSVReaderFromPath(path string, optFns ...func(*ReaderOptions)) (Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return NewCSVReader(file, optFns...)
}

func NewCSVReader(r io.Reader, optFns ...func(*ReaderOptions)) (Reader, error) {
	opt := &ReaderOptions{
		HasNoHeader: false,
		Separator:   CommaSeparator,
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
