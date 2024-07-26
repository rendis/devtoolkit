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

type RowIterator func(yield func(Row) bool)

type Reader interface {
	SetHeader(header []string)
	Iterator() RowIterator
	GetHeaders() []string
	TotalRows() int
	GroupByColumnIndex(columnIndex int) map[string][]Row
	GroupByColumnIndexes(columnIndexes ...int) map[string][]Row
	GroupByColumnName(columnName string) map[string][]Row
	GroupByColumnNames(columnNames ...string) map[string][]Row
	GetRow(index int) (Row, bool)
	RowToObjet(index int, obj any) (bool, error)
	GetNextIndex(currentIndex int, cycle bool) int
	ToObjects(objs []any) error
}

type Row interface {
	Value(fieldName string) (string, bool)
	Fields() []*RowField
	Values() []string
	AsMap() map[string]string
	LineNumber() int
	ToObject(obj any) error
}

type ReaderOptions struct {
	HasNoHeader bool
	Separator   ReaderSeparator
}

type RowField struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

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

type csvReader struct {
	headers        []string
	headerPosition map[string]int
	records        [][]string
}

func (c *csvReader) SetHeader(header []string) {
	c.headerPosition = make(map[string]int)
	c.headers = header
	for i, v := range header {
		c.headerPosition[v] = i
	}
}

func (c *csvReader) Iterator() RowIterator {
	return func(yield func(Row) bool) {
		for i, record := range c.records {
			r := &row{
				row:            record,
				headers:        c.headers,
				headerPosition: c.headerPosition,
				lineNumber:     i + 1,
			}

			if !yield(r) {
				return
			}
		}
	}
}

func (c *csvReader) GetHeaders() []string {
	headers := make([]string, len(c.headerPosition))
	for k, v := range c.headerPosition {
		headers[v] = k
	}
	return headers
}

func (c *csvReader) TotalRows() int {
	return len(c.records)
}

func (c *csvReader) GroupByColumnIndex(columnIndex int) map[string][]Row {
	if len(c.records) == 0 || columnIndex < 0 || columnIndex >= len(c.records[0]) {
		return nil
	}

	grouped := make(map[string][]Row)
	for i, record := range c.records {
		value := record[columnIndex]
		if _, ok := grouped[value]; !ok {
			grouped[value] = make([]Row, 0)
		}
		r := &row{
			row:            record,
			headers:        c.headers,
			headerPosition: c.headerPosition,
			lineNumber:     i + 1,
		}
		grouped[value] = append(grouped[value], r)
	}
	return grouped
}

func (c *csvReader) GroupByColumnIndexes(columnIndexes ...int) map[string][]Row {
	if len(columnIndexes) == 0 || len(c.records) == 0 {
		return nil
	}

	grouped := make(map[string][]Row)
	var recordLength = len(c.records[0])

	var groupKeyBuilder = func(record []string, columnIndexes []int) string {
		var groupValues []string
		for _, columnIndex := range columnIndexes {
			if recordLength > columnIndex {
				value := record[columnIndex]
				groupValues = append(groupValues, value)
			}
		}
		return strings.Join(groupValues, ":")
	}

	for i, record := range c.records {
		// build group key
		groupKey := groupKeyBuilder(record, columnIndexes)

		// add to group
		if _, ok := grouped[groupKey]; !ok {
			grouped[groupKey] = make([]Row, 0)
		}
		r := &row{
			row:            record,
			headers:        c.headers,
			headerPosition: c.headerPosition,
			lineNumber:     i + 1,
		}
		grouped[groupKey] = append(grouped[groupKey], r)
	}
	return grouped
}

func (c *csvReader) GroupByColumnName(columnName string) map[string][]Row {
	if i, ok := c.headerPosition[columnName]; ok {
		return c.GroupByColumnIndex(i)
	}
	return nil
}

func (c *csvReader) GroupByColumnNames(columnNames ...string) map[string][]Row {
	var columnIndexes []int
	for _, columnName := range columnNames {
		if i, ok := c.headerPosition[columnName]; ok {
			columnIndexes = append(columnIndexes, i)
		}
	}
	return c.GroupByColumnIndexes(columnIndexes...)
}

func (c *csvReader) GetRow(index int) (Row, bool) {
	if index < 0 || index >= len(c.records) {
		return nil, false
	}

	return &row{
		row:            c.records[index],
		headers:        c.headers,
		headerPosition: c.headerPosition,
		lineNumber:     index + 1,
	}, true
}

func (c *csvReader) RowToObjet(index int, obj any) (bool, error) {
	r, ok := c.GetRow(index)
	if !ok {
		return false, nil
	}
	return true, r.ToObject(obj)
}

func (c *csvReader) GetNextIndex(currentIndex int, cycle bool) int {
	if currentIndex+1 >= len(c.records) {
		if cycle {
			return 0
		}
		return -1
	}
	return currentIndex + 1
}

func (c *csvReader) ToObjects(objs []any) error {
	var csvStr = ""
	if len(c.headers) > 0 {
		csvStr = strings.Join(c.headers, ",") + "\n"
	}

	for _, record := range c.records {
		csvStr += strings.Join(record, ",") + "\n"
	}

	return decodeObject(csvStr, objs)
}

func (c *csvReader) loadRows(reader *csv.Reader, opts *ReaderOptions) error {
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return nil
	}

	if !opts.HasNoHeader {
		c.SetHeader(records[0])
		records = records[1:]
	}

	c.records = records
	return nil
}

type row struct {
	row            []string
	headers        []string
	headerPosition map[string]int
	lineNumber     int
}

func (r *row) Fields() []*RowField {
	fields := make([]*RowField, len(r.row))
	for i, v := range r.headerPosition {
		fields[v] = &RowField{
			Name:  i,
			Value: r.row[v],
		}
	}
	return fields

}

func (r *row) Value(field string) (string, bool) {
	if i, ok := r.headerPosition[field]; ok {
		return r.row[i], true
	}
	return "", false
}

func (r *row) Values() []string {
	return r.row
}

func (r *row) AsMap() map[string]string {
	m := make(map[string]string)
	for i, v := range r.headerPosition {
		m[i] = r.row[v]
	}
	return m
}

func (r *row) LineNumber() int {
	return r.lineNumber
}

func (r *row) ToObject(obj any) error {
	var csvStr = ""
	if len(r.headers) > 0 {
		csvStr = strings.Join(r.headers, ",") + "\n"
	}
	csvStr += strings.Join(r.row, ",")

	return decodeObject(csvStr, obj)
}

func decodeObject(csvStr string, obj any) error {
	reader := csv.NewReader(strings.NewReader(csvStr))
	dec, err := csvutil.NewDecoder(reader)
	if err != nil {
		return err
	}

	return dec.Decode(obj)
}
