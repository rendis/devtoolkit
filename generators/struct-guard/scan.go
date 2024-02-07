package main

import (
	"github.com/rendis/devtoolkit"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type fieldComposedType int32

const (
	fieldComposedTypeComposed fieldComposedType = iota
	fieldComposedTypeArray
	fieldComposedTypeMap
)

type structsAnalysis struct {
	packageName string
	imports     map[string]struct{}
	structs     []map[string][]map[string]string
}

func extractStructsFromFilesInSamePackage(filesPath []string) (*structsAnalysis, error) {
	var structs = &structsAnalysis{
		imports: make(map[string]struct{}),
	}
	for _, filePath := range filesPath {
		pqName, imports, structMap, err := extractStructsFromFile(filePath)
		if err != nil {
			return nil, err
		}
		if structs.packageName == "" {
			structs.packageName = pqName
		}

		structs.structs = append(structs.structs, structMap)
		for k := range imports {
			structs.imports[k] = struct{}{}
		}
	}
	return structs, nil
}

func extractStructsFromFile(filePath string) (string, map[string]bool, map[string][]map[string]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return "", nil, nil, err
	}

	var structs = make(map[string][]map[string]string)

	var imports = make(map[string]bool)

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)

		if !ok {
			continue
		}

		if genDecl.Tok == token.IMPORT {
			for _, spec := range genDecl.Specs {
				importSpec, ok := spec.(*ast.ImportSpec)
				if !ok {
					continue
				}
				importPath := strings.Trim(importSpec.Path.Value, "\"")
				imports[importPath] = false
			}
			continue
		}

		if genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			var fields []map[string]string
			for _, field := range structType.Fields.List {
				for _, fieldName := range field.Names {
					fieldTypeStr, composedTyp, composedTypDesc1, composedTypDesc2, ok := getFieldTypeFromExpr(field.Type, "")
					if !ok {
						continue
					}

					fields = append(fields, map[string]string{
						"OriginalName":        fieldName.Name,
						"FieldNameLowerCamel": firstToLower(fieldName.Name),
						"FieldNameUpperCamel": firstToUpper(fieldName.Name),
						"FieldType":           fieldTypeStr,
						"IsArray":             devtoolkit.IfThenElse(composedTyp == fieldComposedTypeArray, "true", "false"),
						"IsMap":               devtoolkit.IfThenElse(composedTyp == fieldComposedTypeMap, "true", "false"),
						"ComposedTypeDesc1":   composedTypDesc1,
						"ComposedTypeDesc2":   composedTypDesc2,
					})
				}
			}

			structs[typeSpec.Name.Name] = fields
		}
	}

	var packageName = node.Name.Name
	return packageName, imports, structs, nil
}

func firstToLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func firstToUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func getFieldTypeFromExpr(expr ast.Expr, prefix string) (string, fieldComposedType, string, string, bool) {
	switch expr.(type) {
	case *ast.Ident:
		typ := prefix + expr.(*ast.Ident).Name
		return typ, fieldComposedTypeComposed, "", "", true
	case *ast.StarExpr:
		return getFieldTypeFromExpr(expr.(*ast.StarExpr).X, prefix+"*")
	case *ast.SelectorExpr:
		se := expr.(*ast.SelectorExpr)
		typ := se.X.(*ast.Ident).Name + "." + se.Sel.Name
		typ = prefix + typ
		return typ, fieldComposedTypeComposed, "", "", true
	case *ast.ArrayType:
		at := expr.(*ast.ArrayType)
		//typ, _, ok := getFieldTypeFromExpr(at.Elt, prefix+"[]")
		withoutPrefixTyp, _, _, _, ok := getFieldTypeFromExpr(at.Elt, prefix)
		typ := "[]" + withoutPrefixTyp
		return typ, fieldComposedTypeArray, withoutPrefixTyp, "", ok
	case *ast.MapType:
		mt := expr.(*ast.MapType)
		keyType, _, _, _, ok := getFieldTypeFromExpr(mt.Key, "")
		if !ok {
			return "", fieldComposedTypeComposed, "", "", false
		}

		valueType, _, _, _, ok := getFieldTypeFromExpr(mt.Value, "")
		if !ok {
			return "", fieldComposedTypeComposed, "", "", false
		}
		typ := prefix + "map[" + keyType + "]" + valueType
		return typ, fieldComposedTypeMap, keyType, valueType, true
	}

	return "", fieldComposedTypeComposed, "", "", false
}
