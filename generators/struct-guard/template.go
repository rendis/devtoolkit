package main

const wrapperHeaderTemplate = `// Code generated by 'devtoolkit/generators/struct-guard'. DO NOT EDIT.
// Any changes made to this file will be lost when the file is regenerated

package {{.PackageName}}

{{- range .Imports }}
import "{{.}}"
{{- end }}

{{- .Content }}
`

const wrapperStructTemplate = `
{{- $typeName := .TypeName }}
{{- $wrapperName := .WrapperName }}
// {{$wrapperName}} wraps {{$typeName}} with changes tracking
type {{$wrapperName}} struct {
    {{$typeName}}
    changes {{$typeName}}Changes
}

// {{$typeName}}Changes is a struct to track changes in {{$typeName}}
type {{$typeName}}Changes struct {
    {{- range .Fields }}
    {{.FieldNameLowerCamel}}Changed bool
    {{- end }}
}

// ResetChanges resets the changes in {{$typeName}}
func (w *{{$wrapperName}}) ResetChanges() {
	w.changes = {{$typeName}}Changes{}
}

{{- range .Fields }}
// Get{{.FieldNameUpperCamel}} returns the value of {{$typeName}}.{{.OriginalName}}
func (w *{{$wrapperName}}) Get{{.FieldNameUpperCamel}}() {{.FieldType}} {
    return w.{{$typeName}}.{{.OriginalName}}
}

// Get{{.FieldNameUpperCamel}}WithChange returns the value of {{$typeName}}.{{.OriginalName}} and a boolean indicating if the value has changed
func (w *{{$wrapperName}}) Get{{.FieldNameUpperCamel}}WithChange() ({{.FieldType}}, bool) {
    return w.{{$typeName}}.{{.OriginalName}}, w.changes.{{.FieldNameLowerCamel}}Changed
}

// Set{{.FieldNameUpperCamel}} sets the value of {{$typeName}}.{{.OriginalName}}
func (w *{{$wrapperName}}) Set{{.FieldNameUpperCamel}}(value {{.FieldType}}) {
    w.{{$typeName}}.{{.OriginalName}} = value
    w.changes.{{.FieldNameLowerCamel}}Changed = true
}

{{- if eq .IsArray "true" }}
// GetLast{{.FieldNameUpperCamel}} returns the last value of {{$typeName}}.{{.OriginalName}}
func (w *{{$wrapperName}}) GetLast{{.FieldNameUpperCamel}}() ({{.ComposedTypeDesc1}}, bool) {
	if len(w.{{$typeName}}.{{.OriginalName}}) == 0 {
		var zero {{.ComposedTypeDesc1}}
		return zero, false
	}
	return w.{{$typeName}}.{{.OriginalName}}[len(w.{{$typeName}}.{{.OriginalName}})-1], true
}

// GetLast{{.FieldNameUpperCamel}}WithChange returns the last value of {{$typeName}}.{{.OriginalName}} and a boolean indicating if the value has changed
func (w *{{$wrapperName}}) GetLast{{.FieldNameUpperCamel}}WithChange() ({{.ComposedTypeDesc1}}, bool) {
	if len(w.{{$typeName}}.{{.OriginalName}}) == 0 {
		var zero {{.ComposedTypeDesc1}}
		return zero, w.changes.{{.FieldNameLowerCamel}}Changed
	}
	return w.{{$typeName}}.{{.OriginalName}}[len(w.{{$typeName}}.{{.OriginalName}})-1], w.changes.{{.FieldNameLowerCamel}}Changed
}

// AppendTo{{.FieldNameUpperCamel}} appends a value to {{$typeName}}.{{.OriginalName}}
func (w *{{$wrapperName}}) AppendTo{{.FieldNameUpperCamel}}(value {{.ComposedTypeDesc1}}) {
	w.{{$typeName}}.{{.OriginalName}} = append(w.{{$typeName}}.{{.OriginalName}}, value)
	w.changes.{{.FieldNameLowerCamel}}Changed = true
}
{{ end }}

{{- if eq .IsMap "true" }}
// AddTo{{.FieldNameUpperCamel}} adds a value to {{$typeName}}.{{.OriginalName}}
func (w *{{$wrapperName}}) AddTo{{.FieldNameUpperCamel}}(key {{.ComposedTypeDesc1}}, value {{.ComposedTypeDesc2}}) {
	if w.{{$typeName}}.{{.OriginalName}} == nil {
		w.{{$typeName}}.{{.OriginalName}} = make({{.FieldType}})
	}
	w.{{$typeName}}.{{.OriginalName}}[key] = value
	w.changes.{{.FieldNameLowerCamel}}Changed = true
}

// RemoveFrom{{.FieldNameUpperCamel}} removes a value from {{$typeName}}.{{.OriginalName}}
func (w *{{$wrapperName}}) RemoveFrom{{.FieldNameUpperCamel}}(key {{.ComposedTypeDesc1}}) {
	if w.{{$typeName}}.{{.OriginalName}} == nil {
		return
	}
	delete(w.{{$typeName}}.{{.OriginalName}}, key)
	w.changes.{{.FieldNameLowerCamel}}Changed = true
}

// Get{{.FieldNameUpperCamel}}Value returns the value of {{$typeName}}.{{.OriginalName}} for the given key
func (w *{{$wrapperName}}) Get{{.FieldNameUpperCamel}}Value(key {{.ComposedTypeDesc1}}) ({{.ComposedTypeDesc2}}, bool) {
	if w.{{$typeName}}.{{.OriginalName}} == nil {
		var zero {{.ComposedTypeDesc2}}
		return zero, false
	}
	value, ok := w.{{$typeName}}.{{.OriginalName}}[key]
	return value, ok
}
{{ end }}

{{- if eq .IsPtr "true" }}
// Is{{.FieldNameUpperCamel}}Nil returns true if {{$typeName}}.{{.OriginalName}} is nil
func (w *{{$wrapperName}}) Is{{.FieldNameUpperCamel}}Nil() bool {
	return w.{{$typeName}}.{{.OriginalName}} == nil
}

// Get{{.FieldNameUpperCamel}}Value returns the value of {{$typeName}}.{{.OriginalName}} and a boolean indicating if the value is not nil
func (w *{{$wrapperName}}) Get{{.FieldNameUpperCamel}}Value() ({{.PtrFieldType}}, bool) {
	if w.{{$typeName}}.{{.OriginalName}} == nil {
		var zero {{.PtrFieldType}}
		return zero, false
	}
	return *w.{{$typeName}}.{{.OriginalName}}, true
}

// Get{{.FieldNameUpperCamel}}OrZeroValue returns the value of {{$typeName}}.{{.OriginalName}} and a zero value if the value is nil
func (w *{{$wrapperName}}) Get{{.FieldNameUpperCamel}}OrZeroValue() {{.PtrFieldType}} {
	if w.{{$typeName}}.{{.OriginalName}} == nil {
		var zero {{.PtrFieldType}}
		return zero
	}
	return *w.{{$typeName}}.{{.OriginalName}}
}
{{ end }}

{{ end }}

// ToBuilder returns a builder for {{$wrapperName}}
func (w *{{$wrapperName}}) ToBuilder() *{{$wrapperName}}Builder {
	return &{{$wrapperName}}Builder{wrapper: w}
}

// {{$wrapperName}}Builder is a builder for {{$wrapperName}}
type {{$wrapperName}}Builder struct {
    wrapper *{{$wrapperName}}
}

// New{{$wrapperName}}Builder returns a new {{$wrapperName}}Builder
func New{{$wrapperName}}Builder() *{{$wrapperName}}Builder {
    return &{{$wrapperName}}Builder{wrapper: New{{$wrapperName}}()}
}

// Build returns the built {{$wrapperName}}
func (b *{{$wrapperName}}Builder) Build() *{{$wrapperName}} {
    return b.wrapper
}

{{- range .Fields }}
// With{{.FieldNameUpperCamel}} sets the value of {{$typeName}}.{{.OriginalName}} and returns the builder
// This method only sets the value of {{$typeName}}.{{.OriginalName}} and does not track changes
func (b *{{$wrapperName}}Builder) With{{.FieldNameUpperCamel}}(value {{.FieldType}}) *{{$wrapperName}}Builder {
    b.wrapper.{{$typeName}}.{{.OriginalName}} = value
    return b
}
{{ end }}

// New{{$wrapperName}} returns a new {{$wrapperName}}
func New{{$wrapperName}}() *{{$wrapperName}} {
    return &{{$wrapperName}}{}
}

// New{{$wrapperName}}From returns a new {{$wrapperName}} with the given {{$typeName}}
func New{{$wrapperName}}From({{$typeName}} {{$typeName}}) *{{$wrapperName}} {
	return &{{$wrapperName}}{
		{{$typeName}}: {{$typeName}},
	}
}
`
