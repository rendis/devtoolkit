package main

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

{{- range .Fields }}
// Get{{.FieldNameUpperCamel}} returns the value of {{$typeName}}.{{.OriginalName}}
func (w *{{$wrapperName}}) Get{{.FieldNameUpperCamel}}() {{.FieldType}} {
    return w.{{$typeName}}.{{.OriginalName}}
}

// Set{{.FieldNameUpperCamel}} sets the value of {{$typeName}}.{{.OriginalName}}
func (w *{{$wrapperName}}) Set{{.FieldNameUpperCamel}}(value {{.FieldType}}) {
    w.{{$typeName}}.{{.OriginalName}} = value
    w.changes.{{.FieldNameLowerCamel}}Changed = true
}

// Get{{.FieldNameUpperCamel}}WithChange returns the value of {{$typeName}}.{{.OriginalName}} and a boolean indicating if the value has changed
func (w *{{$wrapperName}}) Get{{.FieldNameUpperCamel}}WithChange() ({{.FieldType}}, bool) {
    return w.{{$typeName}}.{{.OriginalName}}, w.changes.{{.FieldNameLowerCamel}}Changed
}
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
    b.wrapper.{{.OriginalName}} = value
    return b
}
{{ end }}

// New{{$wrapperName}} returns a new {{$wrapperName}}
func New{{$wrapperName}}() *{{$wrapperName}} {
    return &{{$wrapperName}}{}
}
`