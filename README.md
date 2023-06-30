# Devtoolkit

Devtoolkit is a powerful and ever-expanding toolkit designed to streamline daily tasks in Golang software development. 
Within this library, you'll find an array of features, such as tools for working with yml or json prop files, slices, handling generic objects, managing concurrency, and more. 
As Devtoolkit continues to evolve, it will encompass even more functionalities to cater to a variety of programming needs.

## Installation

```
go get github.com/rendis/devtoolkit
```

## Usage

### Working with Concurrency

For handling concurrency, devtoolkit provides a `ConcurrentExec` type that allows to execute a series of functions concurrently.
The library offers a convenient way to manage concurrent execution, allowing to cancel execution, retrieve results and errors, and check when execution is done.

#### Creating and running concurrent functions

Here's an example of how to create and run concurrent functions with `ConcurrentExec`.

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

---

### Load Configuration Propertieswith Environment Variables

Utility functions for loading configuration properties from JSON or YAML files.
This functionality supports the injection of environment variables directly into the configuration properties.

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
    DBConfig `json:"dbConfig" yaml:"dbConfig"`
}

type DBConfig struct {
    Host string `json:"host" yaml:"host"`
    Port int `json:"port" yaml:"port"`
    Username string `json:"username" yaml:"username"`
    Password string `json:"password" yaml:"password"`
    Description string `json:"description" yaml:"description"`
}

cfg := &Config{}
err := devtoolkit.LoadPropFile("config.json", cfg)
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

---

### Working with Slices

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

##### LastIndexOfWithPredicate

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

## Map

`Map`

applies a transformation function to all items in a slice and returns a new slice containing the results.

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

---

## Contributions

Contributions to this library are welcome. Please open an issue to discuss the enhancement or feature you would like to add, or just make a pull request.

## License

Devtoolkit is licensed under the MIT License. Please see the [LICENSE](LICENSE) file for details.