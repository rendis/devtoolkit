package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func extractStructsFromFiles(filesPath []string) (string, []map[string][]map[string]string, error) {
	var structs []map[string][]map[string]string
	var packageName string
	for _, filePath := range filesPath {
		pqName, structMap, err := extractStructsFromFile(filePath)
		if err != nil {
			return "", nil, err
		}
		if packageName == "" {
			packageName = pqName
		}
		structs = append(structs, structMap)
	}
	return packageName, structs, nil
}

func extractStructsFromFile(filePath string) (string, map[string][]map[string]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return "", nil, err
	}

	var structs = make(map[string][]map[string]string)

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
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
					var fieldType string
					switch field.Type.(type) {
					case *ast.Ident:
						fieldType = field.Type.(*ast.Ident).Name
					case *ast.StarExpr:
						fieldType = field.Type.(*ast.StarExpr).X.(*ast.Ident).Name
						fieldType = "*" + fieldType
					default:
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
	return packageName, structs, nil
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

func getFileDirPath(filePath string) string {
	if filePath == "" {
		return ""
	}
	return filePath[:strings.LastIndex(filePath, "/")+1]
}
