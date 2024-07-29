# CSV Reader Library

This library is built on top of Go's `encoding/csv` package, providing additional functionality to read, process, and construct CSV files in a structured and easy-to-use manner. 
It enhances the basic CSV handling capabilities with features such as header management, row iteration, grouping, and object conversion, making it easier to work with CSV data in various applications. 

Additionally, it leverages the [github.com/jszwec/csvutil](https://github.com/jszwec/csvutil) library for deserializing CSV rows into objects, simplifying the conversion process.


## Features

- Read CSV files with different separators.
- Manage headers of CSV files.
- Iterate over the rows of the CSV file.
- Group rows by columns.
- Convert rows to specific objects.
- Handle CSV files with or without headers.
- Construct CSV files with specified data and options.

## Construction Methods

### `NewCSVReaderFromPath`

Creates a new CSV reader from a file path with optional `ReaderOptions`.

```go
func NewCSVReaderFromPath(path string, optFns ...func (*ReaderOptions)) (Reader, error)
```


- `path string`: The path to the CSV file.
- `optFns ...func(*ReaderOptions)`: Optional functions to set `ReaderOptions`.

### `NewCSVReader`

Creates a new CSV reader from an `io.Reader` with optional `ReaderOptions`.

```go
func NewCSVReader(r io.Reader, optFns ...func(*ReaderOptions)) (Reader, error)
```

- `r io.Reader`: The reader to read CSV data from.
- `optFns ...func(*ReaderOptions)`: Optional functions to set `ReaderOptions`.


### `ToReaderSeparator`

Converts a string to a `ReaderSeparator`.

```go
func ToReaderSeparator(separator string) (ReaderSeparator, bool)
```

- `separator string`: The string representation of the separator.
- Returns `ReaderSeparator`: The corresponding `ReaderSeparator` constant.
- Returns `bool`: Indicates if the conversion was successful.

Additionally, you can create a `ReaderSeparator` by casting a `rune` type.

#### Example

```go
separator := ReaderSeparator(',')
```

## Available Methods

### `Reader`

#### Methods

- `SetHeader(header []string)`: Sets the header of the CSV file.
- `Iterator() RowIterator`: Returns an iterator for iterating over rows.
- `GetHeaders() []string`: Returns the headers of the CSV file.
- `TotalRows() int`: Returns the total number of rows in the CSV file.
- `GroupByColumnIndex(columnIndex int) map[string][]Row`: Groups rows by the value at the specified column index.
- `GroupByColumnIndexes(columnIndexes ...int) map[string][]Row`: Groups rows by the values at the specified column
  indexes.
- `GroupByColumnName(columnName string) map[string][]Row`: Groups rows by the value of the specified column name.
- `GroupByColumnNames(columnNames ...string) map[string][]Row`: Groups rows by the values of the specified column names.
- `GetRow(index int) (Row, bool)`: Returns the row at the specified index.
- `RowToObjet(index int, obj any) (bool, error)`: Converts the row at the specified index to the specified object.
- `GetNextIndex(currentIndex int, cycle bool) int`: Returns the next index based on the current index and cycle option.
- `ToObjects(objs []any) error`: Converts all rows to the specified slice of objects.

### `Row`

#### Methods

- `Value(columnName string) (string, bool)`: Returns the value of the specified column name.
- `Fields() []*RowField`: Returns the fields of the row.
- `Values() []string`: Returns the values of the row.
- `AsMap() map[string]string`: Returns the row as a map with column names as keys.
- `LineNumber() int`: Returns the line number of the row in the CSV file.
- `ToObject(obj any) error`: Converts the row to the specified object.

### `RowField`

#### Fields

- `Name string`: The name of the field.
- `Value string`: The value of the field.

### `ReaderOptions`

#### Fields

- `NoHeader bool`: Indicates if the CSV file has no header.
- `Separator ReaderSeparator`: The separator used in the CSV file.

### `ReaderSeparator`
#### Constants

- `CommaSeparator`: Separator for comma (`,`).
- `SemicolonSeparator`: Separator for semicolon (`;`).
- `TabSeparator`: Separator for tab (`\t`).
- `PipeSeparator`: Separator for pipe (`|`).


## Example Usage

Here is an example demonstrating how to use the CSV Reader Library to read a CSV file and convert its rows into objects.

```go
type ExampleStruct struct {
	Value1 string `csv:"value1"`
	Value2 string `csv:"value2"`
}

func main() {
	path := "./example/test.csv"

	// Create a new CSV reader from a file path with optional ReaderOptions
	reader, err := csvreader.NewCSVReaderFromPath(path)
	if err != nil {
		log.Fatalf("Error creating CSV reader: %v", err)
	}

	// Iterate over the rows and convert each row to an ExampleStruct object
	for item := range reader.Iterator() {
		var example ExampleStruct
		if err := item.ToObject(&example); err != nil {
			log.Fatalf("Error converting row to object: %v", err)
		}
		fmt.Printf("Value1: %s, Value2: %s\n", example.Value1, example.Value2)
	}
}
```