# Devtoolkit

[![Go Report Card](https://goreportcard.com/badge/github.com/rendis/devtoolkit)](https://goreportcard.com/report/github.com/rendis/devtoolkit)
[![Go Reference](https://pkg.go.dev/badge/github.com/rendis/devtoolkit.svg)](https://pkg.go.dev/github.com/rendis/devtoolkit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/rendis/devtoolkit.svg?style=flat-square)](https://github.com/rendis/devtoolkit/releases)

Devtoolkit is a powerful and ever-expanding toolkit designed to streamline daily tasks in Golang software development. 
Within this library, you'll find an array of features, such as tools for working with yml or json prop files, slices, handling generic objects, managing concurrency, and more. 
As Devtoolkit continues to evolve, it will encompass even more functionalities to cater to a variety of programming needs.

## Table of Contents
- [Devtoolkit](#devtoolkit)
    * [Installation](#installation)
    * [Usage](#usage)
        + [Concurrent solutions](#concurrent-solutions)
            - [Running concurrent functions](#running-concurrent-functions)
            - [Running concurrent workers](#running-concurrent-workers)
        + [Load properties from a file (JSON/YAML) with environment variable injections and validations](#load-properties-from-a-file-with-environment-variable-injections-and-validations)
        + [Resilience](#resilience)
            - [RetryOperation](#retryoperation)
        + [Design Patterns](#design-patterns)
            - [Process Chain](#process-chain)
        + [Data structures](#data-structures)
            - [Pair](#pair)
            - [Triple](#triple)
        + [Readers](#readers)
            - [CSV Reader](#csv-reader)
        + [Generators](#generators)
            - [struct-guard](#struct-guard)
        + [Working with Generic Objects](#working-with-generic-objects)
            - [ToPtr](#toptr)
            - [IsZero](#iszero)
            - [StructToMap](#structtomap)
            - [MapToStruct](#maptostruct)
            - [CastToPointer](#casttopointer)
            - [IfThenElse](#ifthenelse)
            - [IfThenElseFn](#ifthenelsefn)
            - [DefaultIfNil](#defaultifnil)
            - [ZeroValue](#zerovalue)
            - [ToInt](#toint)
            - [ToFloat64](#tofloat64)
            - [StrToStruct](#strtostruct)
        + [Working with Slices](#working-with-slices)
            - [Contains](#contains)
            - [ContainsWithPredicate](#containswithpredicate)
            - [IndexOf](#indexof)
            - [IndexOfWithPredicate](#indexofwithpredicate)
            - [LastIndexOf](#lastindexof)
            - [LastIndexOfWithPredicate](#lastindexofwithpredicate)
            - [Remove](#remove)
            - [RemoveWithPredicate](#removewithpredicate)
            - [RemoveAll](#removeall)
            - [RemoveAllWithPredicate](#removeallwithpredicate)
            - [RemoveAt](#removeat)
            - [RemoveRange](#removerange)
            - [RemoveIf](#removeif)
            - [Filter](#filter)
            - [FilterNot](#filternot)
            - [Map](#map)
            - [RemoveDuplicates](#removeduplicates)
            - [Reverse](#reverse)
            - [Difference](#difference)
            - [Intersection](#intersection)
            - [Union](#union)
            - [GetMapKeys](#getmapkeys)
            - [GetMapValues](#getmapvalues)
    * [Contributions](#contributions)
    * [License](#license)

## Installation

```
go get github.com/rendis/devtoolkit
```

## Usage

### Concurrent solutions

#### ConcurrentExec

`ConcurrentExec` is a utility for executing a series of functions concurrently.

```go
var fns []devtoolkit.ConcurrentFn

fns = append(fns, func(ctx context.Context) (any, error) {
   // Implement function logic
   return "Result1", nil
})

fns = append(fns, func(ctx context.Context) (any, error) {
   // Implement function logic
   return "Result2", nil
})

ctx := context.Background()
ce, err := devtoolkit.NewConcurrentExec().ExecuteFns(ctx, fns...)

if err != nil {
   fmt.Println(err)
}

// errors is a slice of errors returned by the functions, where each index corresponds to the function at the same index in the fns slice
errors := ce.Errors()
for _, err := range errors {
   if err != nil {
      fmt.Println(err)
   }
}

// results is a slice of results returned by the functions, where each index corresponds to the function at the same index in the fns slice
// Note: results are of type any, so you'll need to cast them to the appropriate type
results := ce.Results()
for _, res := range results {
   fmt.Println(res)
}
```

Note: This example does not include error handling, be sure to do so in your implementations.

#### ConcurrentWorkers

`ConcurrentWorkers` is a utility for executing a series of functions concurrently using a pool of workers.

```go
var maxWorkers = 5
var cw = NewConcurrentWorkers(maxWorkers)

for i := 0; i < 10; i++ {
   cw.Execute(func() {
      // do something cool
   })
}

// Stop the workers with an error
cw.Stop(fmt.Errorf("Something went wrong"))

// Stop the workers without an error
cw.Stop(nil)

// Wait for all workers to finish
cw.Wait()
```

#### ConcurrentManager

`ConcurrentManager` is a struct that dynamically manages a pool of workers within set limits. 
It adjusts worker count based on load, offering functions to allocate, release, and wait for workers, 
optimizing concurrent task handling.

```go
// NewConcurrentManager creates a new instance of ConcurrentManager with specified parameters.
// It ensures that the provided parameters are within acceptable ranges and initializes the manager.
func NewConcurrentManager(minWorkers, maxWorkers int, workerIncreaseRate float64, timeIncreasePeriod time.Duration) *ConcurrentManager

// Allocate requests a new worker to be allocated.
// It blocks if the maximum number of workers has been reached, until a worker is released.
func (c *ConcurrentManager) Allocate()

// Release frees up a worker, making it available for future tasks.
// It only releases a worker if the release condition is met, ensuring resources are managed efficiently.
func (c *ConcurrentManager) Release()

// Wait blocks until all workers have finished their tasks.
// It ensures that all resources are properly cleaned up before shutting down or reinitializing the manager.
func (c *ConcurrentManager) Wait()
```


#### AtomicNumber

`AtomicNumber` is a utility for managing an number atomically.

```go
// initialize an atomic number with default value 0
var atomic AtomicNumber[int]

atomic.Get() // get the current value
atomic.Set(5) // set the value to 5
atomic.Increment() // increment the value by 1
atomic.IncrementBy(10) // add 10 to the value
atomic.IncrementAndGet() // increment the value by 1 and return the new value
atomic.IncrementByAndGet(10) // add 10 to the value and return the new value
var incremented = atomic.IncrementIf(func(v int) {return v > 0}) // increment the value by 1 if the condition is met
var incrementedBy = atomic.IncrementByIf(10, func(v int) {return v > 0}) // add 10 to the value if the condition is met

// same functions are available for decrementing the value
atomic.Decrement() // decrement the value by 1
...
```



---

### Load properties from a file with environment variable injections and validations

Utility functions for loading configuration properties from JSON or YAML files.
This functionality supports the injection of environment variables directly into the configuration properties.

`LoadPropFile` supports field validation using struct tags provided by the [go-playground/validator](https://github.com/go-playground/validator/v10) library.

You can register your own custom validators using the `RegisterCustomValidator` function.

`devtoolkit` provides the following built-in validators:
- `trimmed-non-empty` - checks whether a string is not empty after trimming whitespace
 
```yaml
dbConfig:
  host: ${DB_HOST}
  port: ${DB_PORT}
  username: ${DB_USERNAME}
  password: ${DB_PASSWORD}
  description: "YAML config file"
```

```json
{
  "dbConfig": {
    "host": "${DB_HOST}",
    "port": "${DB_PORT}",
    "username": "${DB_USERNAME}",
    "password": "${DB_PASSWORD}",
    "description": "JSON config file"
  }
}
```

```go
type Config struct {
    DBConfig `json:"dbConfig" yaml:"dbConfig" validate:"required"`
}

type DBConfig struct {
    Host string `json:"host" yaml:"host" validate:"required"`
    Port int `json:"port" yaml:"port" validate:"required,min=1,max=65535"`
    Username string `json:"username" yaml:"username" validate:"required,trimmed-non-empty"`
    Password string `json:"password" yaml:"password" validate:"required,trimmed-non-empty"`
    Description string `json:"description" yaml:"description" validate:"required"`
}

func (p *Config) SetDefaults() {
    p.DBConfig.Host = "localhost"
    p.DBConfig.Port = 3306
}

func trimmedNonEmpty(fl validator.FieldLevel) bool {
    s := fl.Field().String()
    trimmed := strings.TrimSpace(s)
    return len(trimmed) > 0
}

cfg := &Config{}
props := []devtoolkit.ToolKitProp{cfg}
devtoolkit.RegisterCustomValidator("trimmed-non-empty", trimmedNonEmpty)
err := devtoolkit.LoadPropFile("config.json", props)

```

---

### Resilience

Utility functions and patterns to ensure resilient operation execution.

#### RetryOperation

The `RetryOperation` function retries an operation for a specified number of times with optional exponential backoff. It's useful when operations have a tendency to fail temporarily and may succeed on a subsequent attempt.

```go
type ResilienceOptions struct {
    MaxRetries       int                     // indicates the maximum number of retries. Default is 3.
    WaitTime         time.Duration           // indicates the wait time between retries. Default is 100ms.
    Backoff          bool                    // indicates whether to use exponential backoff. Default is false.
    RawError         bool                    // indicates whether to return the raw error or wrap it in a new error. Default is false.
    IsIgnorableErrorHandler func(error) bool // indicates whether to ignore the error or not. Default is nil.
    ReturnIgnorable  bool                    // indicates whether to return the ignorable error or not. Default is false.
}


func NewResilience(options *ResilienceOptions) (Resilience, error)
```

Example:
```go
operation := func() error {
	return networkCall() // Some hypothetical network operation
}
options := &devtoolkit.ResilienceOptions{
	MaxRetries: 3, 
	WaitTime:   2 * time.Second,
	Backoff:    true,
}

resilienceHandler, err := devtoolkit.NewResilience(options)

if err != nil {
	panic(err)
}

err = resilienceHandler.RetryOperation(operation) // wrapped error returned

if err != nil {
	fmt.Println("Operation failed.", err)
} else {
	fmt.Println("Operation succeeded!")
}
```

With the `RetryOperation` function, users can easily add resiliency to their operations and ensure that temporary failures don't lead to complete system failures.

---

### Design Patterns

#### Process Chain
`ProcessChain` is an implementation that enables the orderly execution of operations on data within a Go application.
Leveraging the "Chain of Responsibility" and "Command" design patterns, it allows for the addition of operations (links)
and an optional save step to ensure data integrity after each operation. Ideal for scenarios requiring a series of
actions on data with the flexibility to add, remove, or modify steps as needed.

```go
type ProcessChain[T any] interface {
    AddLink(string, LinkFn[T]) error
    SetSaveStep(SaveStep[T])
    GetChain() []string
    Execute(context.Context, T) ([]string, error)
    ExecuteWithIgnorableLinks(context.Context, T, []string) ([]string, error)
}

func NewProcessChain[T any](opts *ProcessChainOptions) ProcessChain[T]
```

Example:
```go
type Data struct {
    // Your data fields here
}

func step1(d Data) error {
    // Define a process operation
}

func step2(d Data) error {
    // Define a process operation
}

func saveData(d Data) error {
    // Define a save operation
}

chain := NewProcessChain[Data]()
chain.AddLink("step1", step1)
chain.AddLink("step2", step2)
chain.SetSaveStep(saveData)
err := chain.Execute(Data{})
```
---

### Working with Generic Objects

#### ToPtr

The `ToPtr` function takes a value of any type and returns a pointer to it.

```go
val := 5
ptr := devtoolkit.ToPtr(val)
fmt.Println(*ptr) // Returns 5
```

#### IsZero

The `IsZero` function checks whether a value is the zero value of its type.

```go
fmt.Println(devtoolkit.IsZero(0)) // Returns true
fmt.Println(devtoolkit.IsZero(1)) // Returns false
fmt.Println(devtoolkit.IsZero("")) // Returns true
```

#### StructToMap

The `StructToMap` function converts a struct to a `map[string]any`.

```go
type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

p := Person{
    Name:  "John",
    Age:   30,
    Email: "john@example.com",
}

personMapData, err := devtoolkit.StructToMap(p)
```


#### MapToStruct

The `MapToStruct` function converts a `map[string]any` to a pointer to a struct.

```go
// personMapData is a map[string]any containing the data of a Person struct, see StructToMap example
ptrToNewPerson, err := devtoolkit.MapToStruct[Person](personMapData)
```

#### CastToPointer
`CastToPointer` casts a value to a pointer of the same type.

```go
func CastToPointer[T any](v any) (*T, bool)
```

Rules:
- v must be a pointer.
- if v not a pointer, returns false.
- if v is nil, returns false.
- if v is a pointer but not of the given type, returns false.
- if v is a pointer of the given type, returns true.


#### IfThenElse
`IfThenElse` returns the first value if the condition is true, otherwise it returns the second value.

```go
func IfThenElse[T any](condition bool, first, second T) T
```

#### IfThenElseFn
`IfThenElseFn` returns the first value if the condition is true, otherwise it returns the result of the second value.

```go
func IfThenElseFn[T any](condition bool, first T, second func() T) T
```

#### DefaultIfNil
`DefaultIfNil` returns the first value if it is not nil, otherwise it returns the second value.

```go
func DefaultIfNil[T any](first *T, second T) T
```

#### ZeroValue
`ZeroValue` returns the zero value of a type.

```go
func ZeroValue[T any]() T
```

#### ToInt
`ToInt` converts a value to an int.

```go
func ToInt(value any) (int, bool)
```


#### ToFloat64
`ToFloat64` converts a value to a float64.

```go
func ToFloat64(value any) (float64, bool)
```


#### StrToStruct
`StrToStruct` converts a string to a struct.

```go
func StrToStruct[T any](s string) (*T, error)
```

---

### Data structures

#### Pair

The `Pair` type represents a pair of values.

```go
type Pair[F any, S any] struct {
    First  F
    Second S
}
```

#### Triple

The `Triple` type represents a triple of values.

```go
type Triple[F any, S any, T any] struct {
    First  F
    Second S
    Third  T
}
```

---

### Readers

#### CSV Reader

The CSV reader provides a simple and efficient way to read CSV files in Go.

More details can be found in the [CSV Reader documentation](reader/csv/README.md).

---

### Generators

#### struct-guard

The struct-guard tool generates wrapper structs in Go for tracking changes to the fields of the original struct.

More details can be found in the [struct-generator documentation](generator/struct-guard/README.md).

---

### Working with Slices

Common utility functions for working with slices.

#### Contains

The `Contains` function checks whether a slice contains an item. The item must be comparable.

```go
func Contains[T comparable](slice []T, item T) bool
```

Example:
```go
slice := []int{1, 2, 3}
item := 2
fmt.Println(devtoolkit.Contains(slice, item)) // Returns true
```

#### ContainsWithPredicate

The `ContainsWithPredicate` function checks whether a slice contains an item using a predicate to compare items.

```go
func ContainsWithPredicate[T any](slice []T, item T, predicate func(T, T) bool) bool
```

Example:
```go
slice := []int{1, 2, 3}
item := 2
predicate := func(a, b int) bool { return a == b }
fmt.Println(devtoolkit.ContainsWithPredicate(slice, item, predicate)) // Returns true
```

#### IndexOf

`IndexOf` returns the index of the first instance of an item in a slice, or -1 if the item is not present in the slice.

```go
func IndexOf[T comparable](slice []T, item T) int
```

Example:

```go
index := IndexOf([]int{1, 2, 3, 2, 1}, 2)
fmt.Println(index) // Output: 1
```

#### IndexOfWithPredicate

`IndexOfWithPredicate` returns the index of the first instance of an item in a slice, or -1 if the item is not present in the slice. It uses a predicate function to compare items.

```go
func IndexOfWithPredicate[T any](slice []T, item T, predicate func(T, T) bool) int
```

Example:

```go
index := IndexOfWithPredicate([]string{"apple", "banana", "cherry"}, "APPLE", strings.EqualFold)
fmt.Println(index) // Output: 0
```

#### LastIndexOf

`LastIndexOf` returns the index of the last instance of an item in a slice, or -1 if the item is not present in the slice.

```go
func LastIndexOf[T comparable](slice []T, item T) int
```

Example:

```go
index := LastIndexOf([]int{1, 2, 3, 2, 1}, 2)
fmt.Println(index) // Output: 3
```

#### LastIndexOfWithPredicate

`LastIndexOfWithPredicate` returns the index of the last instance of an item in a slice, or -1 if the item is not present in the slice. It uses a predicate function to compare items.

```go
func LastIndexOfWithPredicate[T any](slice []T, item T, predicate func(T, T) bool) int
```

Example:

```go
index := LastIndexOfWithPredicate([]string{"apple", "banana", "cherry", "apple"}, "APPLE", strings.EqualFold)
fmt.Println(index) // Output: 3
```

#### Remove

`Remove` removes the first instance of an item from a slice, if present. It returns true if the item was removed, false otherwise.

```go
func Remove[T comparable](slice []T, item T) bool
```

Example:

```go
removed := Remove([]int{1, 2, 3, 2, 1}, 2)
fmt.Println(removed) // Output: true
```

#### RemoveWithPredicate

`RemoveWithPredicate` removes the first instance of an item from a slice, if present. It uses a predicate function to compare items. It returns true if the item was removed, false otherwise.

```go
func RemoveWithPredicate[T any](slice []T, item T, predicate func(T, T) bool) bool
```

Example:

```go
removed := RemoveWithPredicate([]string{"apple", "banana", "cherry"}, "APPLE", strings.EqualFold)
fmt.Println(removed) // Output: true
```

#### RemoveAll

`RemoveAll` removes all instances of an item from a slice, if present. It returns true if the item was removed, false otherwise.

```go
func RemoveAll[T comparable](slice []T, item T) bool
```

Example:

```go
removed := RemoveAll([]int{1, 2, 3, 2, 1}, 2)
fmt.Println(removed) // Output: true
```

#### RemoveAllWithPredicate

`RemoveAllWithPredicate` removes all instances of an item from a slice, if present. It uses a predicate function to compare items. It returns true if the item was removed, false otherwise.

```go
func RemoveAllWithPredicate[T any](slice []T, item T, predicate func(T, T) bool) bool
```

Example:

```go
removed := RemoveAllWithPredicate([]string{"apple", "banana", "cherry", "apple"}, "APPLE", strings.EqualFold)
fmt.Println(removed) // Output: true
```

#### RemoveAt

`RemoveAt` removes the item at a given index from a slice. It returns true if the item was removed, false otherwise.

```go
func RemoveAt[T any](slice []T, index int) bool
```

Example:

```go
removed := RemoveAt([]int{1, 2, 3}, 1)
fmt.Println(removed) // Output: true
```

#### RemoveRange

`RemoveRange` removes the items in a given range from a slice. It returns true if items were removed, false otherwise.

```go
func RemoveRange[T any](slice []T, start, end int) bool
```

Example:

```go
removed := RemoveRange([]int{1, 2, 3, 4, 5}, 1, 3)
fmt.Println(removed) // Output: true
```

#### RemoveIf

`RemoveIf` removes all items from a slice for which a predicate function returns true. It returns true if any items were removed, false otherwise.

```go
func RemoveIf[T any](slice []T, predicate func(T) bool) bool
```

Example:

```go
removed := RemoveIf([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 })
fmt.Println(removed) // Output: true
```

#### Filter

`Filter` returns a new slice containing all items from the original slice for which a predicate function returns true.

```go
func Filter[T any](slice []T, predicate func(T) bool) []T
```

Example:

```go
filtered := Filter([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 })
fmt.Println(filtered) // Output: [2 4]
```

#### FilterNot

`FilterNot` returns a new slice containing all items from the original slice for which a predicate function returns false.

```go
func FilterNot[T any](slice []T, predicate func(T) bool) []T
```

Example:

```go
filtered := FilterNot([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 })
fmt.Println(filtered) // Output: [1 3 5]
```


#### Map

`Map` applies a transformation function to all items in a slice and returns a new slice containing the results.

```go
func Map[T, R any](slice []T, mapper func(T) R) []R 
```

Example:

```go
mapped := Map([]int{1, 2, 3}, func(n int) int { return n * 2 })
fmt.Println(mapped) // Output: [2 4 6]
```

#### RemoveDuplicates

`RemoveDuplicates` removes all duplicate items from a slice. It returns true if any items were removed, false otherwise.

```go
func RemoveDuplicates[T comparable](slice []T) bool
```

Example:

```go
removed := RemoveDuplicates([]int{1, 2, 3, 2, 1})
fmt.Println(removed) // Output: true
```

#### Reverse

`Reverse` reverses the order of items in a slice.

```go
func Reverse[T any](slice []T)
```

Example:

```go
data := []int{1, 2, 3}
Reverse(data)
fmt.Println(data) // Output: [3 2 1]
```

#### Difference

`Difference` returns a new slice containing all items from the original slice that are not present in the other slice.

```go
func Difference[T comparable](slice, other []T) []T
```

Example:

```go
diff := Difference([]int{1, 2, 3, 4, 5}, []int{3, 4, 5, 6, 7})
fmt.Println(diff) // Output: [1 2]
```

#### Intersection

`Intersection` returns a new slice containing all items from the original slice that are also present in the other slice.

```go
func Intersection[T comparable](slice, other []T) []T
```

Example:

```go
inter := Intersection([]int{1, 2, 3, 4, 5}, []int{3, 4, 5, 6, 7})
fmt.Println(inter) // Output: [3 4 5]
```

#### Union

`Union` returns a new slice containing all unique items from both the original slice and the other slice.

```go
func Union[T comparable](slice, other []T) []T
```

Example:

```go
union := Union([]int{1, 2, 3}, []int{3, 4, 5})
fmt.Println(union) // Output: [1 2 3 4 5]
```

#### GetMapKeys
`GetMapKeys` returns a slice of keys from a map.

```go
func GetMapKeys[K comparable, V any](m map[K]V) []K
```

Example:

```go
keys := GetMapKeys(map[string]int{"a": 1, "b": 2})
fmt.Println(keys) // Output: [a b]
```


#### GetMapValues
`GetMapValues` returns a slice of values from a map.

```go
func GetMapValues[K comparable, V any](m map[K]V, removeDuplicates bool) []V
```

Example:

```go
values := GetMapValues(map[string]int{"a": 1, "b": 2, "c": 1}, false)
fmt.Println(values) // Output: [1 2 1]

values = GetMapValues(map[string]int{"a": 1, "b": 2, "c": 1}, true)
fmt.Println(values) // Output: [1 2]
```

## Contributions

Contributions to this library are welcome. Please open an issue to discuss the enhancement or feature you would like to add, or just make a pull request.

## License

Devtoolkit is licensed under the MIT License. Please see the [LICENSE](LICENSE) file for details.