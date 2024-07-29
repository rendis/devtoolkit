package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type fieldComposedType int32

const (
	fieldComposedTypeNotComposed fieldComposedType = iota
	fieldComposedTypeArray
	fieldComposedTypeMap
)

type structsAnalysis struct {
	packageName string
	imports     map[string]struct{}
	structs     []map[string][]map[string]string
}

type fieldTypeInfo struct {
	fieldTypeStr     string
	composedTyp      fieldComposedType
	isArray          bool
	isMap            bool
	isPtr            bool
	ptrFieldTypeStr  string
	composedTypDesc1 string
	composedTypDesc2 string
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

func extractStructsFromFile(filePath string) (string, map[string]struct{}, map[string][]map[string]string, error) {
	fSet := token.NewFileSet()
	node, err := parser.ParseFile(fSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return "", nil, nil, err
	}

	var structs = make(map[string][]map[string]string)

	var imports = make(map[string]struct{})

	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			processGenDecl(genDecl, imports, structs)
		}
	}

	var packageName = node.Name.Name
	return packageName, imports, structs, nil
}

func processGenDecl(genDecl *ast.GenDecl, imports map[string]struct{}, structs map[string][]map[string]string) {
	if genDecl.Tok == token.IMPORT {
		for _, spec := range genDecl.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}
			var alias string
			if importSpec.Name != nil {
				alias = importSpec.Name.Name + " "
			}
			importPath := alias + importSpec.Path.Value
			imports[importPath] = struct{}{}
		}
		return
	}

	if genDecl.Tok != token.TYPE {
		return
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
				fieldInfo := getFieldTypeFromExpr(field.Type)
				if fieldInfo == nil {
					continue
				}

				fields = append(fields, map[string]string{
					"OriginalName":        fieldName.Name,
					"FieldNameLowerCamel": firstToLower(fieldName.Name),
					"FieldNameUpperCamel": firstToUpper(fieldName.Name),
					"FieldType":           fieldInfo.fieldTypeStr,
					"IsArray":             fmt.Sprintf("%t", fieldInfo.isArray),
					"IsMap":               fmt.Sprintf("%t", fieldInfo.isMap),
					"IsPtr":               fmt.Sprintf("%t", fieldInfo.isPtr),
					"PtrFieldType":        fieldInfo.ptrFieldTypeStr,
					"ComposedTypeDesc1":   fieldInfo.composedTypDesc1,
					"ComposedTypeDesc2":   fieldInfo.composedTypDesc2,
				})
			}
		}

		structs[typeSpec.Name.Name] = fields
	}
}

func getFieldTypeFromExpr(expr ast.Expr) *fieldTypeInfo {
	switch expr.(type) {
	case *ast.Ident:
		return &fieldTypeInfo{
			fieldTypeStr: expr.(*ast.Ident).Name,
			composedTyp:  fieldComposedTypeNotComposed,
		}
	case *ast.StarExpr:
		typeInfo := getFieldTypeFromExpr(expr.(*ast.StarExpr).X)
		return &fieldTypeInfo{
			fieldTypeStr:    "*" + typeInfo.fieldTypeStr,
			ptrFieldTypeStr: typeInfo.fieldTypeStr,
			composedTyp:     fieldComposedTypeNotComposed,
			isPtr:           true,
		}
	case *ast.SelectorExpr:
		se := expr.(*ast.SelectorExpr)
		typ := se.X.(*ast.Ident).Name + "." + se.Sel.Name
		return &fieldTypeInfo{
			fieldTypeStr: typ,
			composedTyp:  fieldComposedTypeNotComposed,
		}
	case *ast.ArrayType:
		at := expr.(*ast.ArrayType)
		typeInfo := getFieldTypeFromExpr(at.Elt)
		return &fieldTypeInfo{
			fieldTypeStr:     "[]" + typeInfo.fieldTypeStr,
			composedTyp:      fieldComposedTypeArray,
			isArray:          true,
			composedTypDesc1: typeInfo.fieldTypeStr,
		}
	case *ast.MapType:
		mt := expr.(*ast.MapType)

		keyInfo := getFieldTypeFromExpr(mt.Key)
		valueInfo := getFieldTypeFromExpr(mt.Value)

		return &fieldTypeInfo{
			fieldTypeStr:     "map[" + keyInfo.fieldTypeStr + "]" + valueInfo.fieldTypeStr,
			composedTyp:      fieldComposedTypeMap,
			isMap:            true,
			composedTypDesc1: keyInfo.fieldTypeStr,
			composedTypDesc2: valueInfo.fieldTypeStr,
		}
	}

	return nil
}
