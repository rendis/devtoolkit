package csv

import "strings"

type Row interface {
	Value(fieldName string) (string, bool)
	Fields() []*RowField
	Values() []string
	AsMap() map[string]string
	LineNumber() int
	ToObject(obj any) error
}

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
