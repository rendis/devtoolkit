package csv

import "strings"

// Row defines the interface for a row in the CSV file.
type Row interface {
	// Value returns the value of the specified column name.
	Value(columnName string) (string, bool)

	// Fields returns the fields of the row.
	Fields() []*RowField

	// Values returns the values of the row.
	Values() []string

	// AsMap returns the row as a map with column names as keys.
	AsMap() map[string]string

	// LineNumber returns the line number of the row in the CSV file.
	LineNumber() int

	// ToObject converts the row to the specified object.
	ToObject(obj any) error
}

// RowField represents a field in a row with a name and value.
type RowField struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
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

func (r *row) Value(columnName string) (string, bool) {
	if i, ok := r.headerPosition[columnName]; ok {
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
