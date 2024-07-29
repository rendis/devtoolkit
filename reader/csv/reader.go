package csv

import (
	"encoding/csv"
	"strings"
)

// RowIterator defines a function type for iterating over rows.
type RowIterator func(yield func(Row) bool)

// Reader defines the interface for reading CSV files and provides various methods to work with the data.
type Reader interface {
	// SetHeader sets the header of the CSV file.
	SetHeader(header []string)

	// Iterator returns a RowIterator for iterating over rows.
	Iterator() RowIterator

	// GetHeaders returns the headers of the CSV file.
	GetHeaders() []string

	// TotalRows returns the total number of rows in the CSV file.
	TotalRows() int

	// GroupByColumnIndex groups rows by the value at the specified column index.
	GroupByColumnIndex(columnIndex int) map[string][]Row

	// GroupByColumnIndexes groups rows by the values at the specified column indexes.
	GroupByColumnIndexes(columnIndexes ...int) map[string][]Row

	// GroupByColumnName groups rows by the value of the specified column name.
	GroupByColumnName(columnName string) map[string][]Row

	// GroupByColumnNames groups rows by the values of the specified column names.
	GroupByColumnNames(columnNames ...string) map[string][]Row

	// GetRow returns the row at the specified index.
	GetRow(index int) (Row, bool)

	// RowToObjet converts the row at the specified index to the specified object.
	RowToObjet(index int, obj any) (bool, error)

	// GetNextIndex returns the next index based on the current index and cycle option.
	GetNextIndex(currentIndex int, cycle bool) int

	// ToObjects converts all rows to the specified slice of objects.
	ToObjects(objs []any) error
}

// ReaderOptions holds options for configuring the CSV Reader.
type ReaderOptions struct {
	NoHeader  bool
	Separator ReaderSeparator
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

	if !opts.NoHeader {
		c.SetHeader(records[0])
		records = records[1:]
	}

	c.records = records
	return nil
}
