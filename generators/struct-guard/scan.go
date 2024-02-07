package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
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
		//structs = append(structs, structMap)
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
					fieldType, ok := getFieldTypeFromExpr(field.Type, "")
					if !ok {
						continue
					}

					fields = append(fields, map[string]string{
						"OriginalName":        fieldName.Name,
						"FieldNameLowerCamel": firstToLower(fieldName.Name),
						"FieldNameUpperCamel": firstToUpper(fieldName.Name),
						"FieldType":           fieldType,
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

func getFieldTypeFromExpr(typ ast.Expr, prefix string) (string, bool) {
	switch typ.(type) {
	case *ast.Ident:
		fieldType := prefix + typ.(*ast.Ident).Name
		return fieldType, true
	case *ast.StarExpr:
		return getFieldTypeFromExpr(typ.(*ast.StarExpr).X, prefix+"*")
	case *ast.SelectorExpr:
		se := typ.(*ast.SelectorExpr)
		fieldType := se.X.(*ast.Ident).Name + "." + se.Sel.Name
		fieldType = prefix + fieldType
		return fieldType, true
	case *ast.ArrayType:
		at := typ.(*ast.ArrayType)
		return getFieldTypeFromExpr(at.Elt, prefix+"[]")
	case *ast.MapType:
		mt := typ.(*ast.MapType)
		keyType, ok := getFieldTypeFromExpr(mt.Key, "")
		if !ok {
			return "", false
		}

		valueType, ok := getFieldTypeFromExpr(mt.Value, "")
		if !ok {
			return "", false
		}
		fieldType := prefix + "map[" + keyType + "]" + valueType
		return fieldType, true
	}

	return "", false
}
