package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"
)

func main() {

	loadGenProp()

	// exclude files to map
	var excludeFilesMap = make(map[string]bool)
	for _, file := range confProps.ExcludeFilesToScan {
		file = filepath.Clean(file)
		excludeFilesMap[file] = true
	}

	// extract to scan
	var filesToScanMap = make(map[string]map[string]struct{})
	for _, path := range confProps.ToScan {
		if isDirectory(path) {
			dir := filepath.Clean(path)
			files, err := listGoFiles(dir)
			if err != nil {
				panic(err)
			}
			for _, file := range files {
				fileName := filepath.Base(file)
				if excludeFilesMap[file] || filepath.Ext(file) != ".go" || fileName == *confProps.GeneratedFileName {
					continue
				}
				if filesToScanMap[dir] == nil {
					filesToScanMap[dir] = make(map[string]struct{})
				}
				filesToScanMap[dir][file] = struct{}{}
			}
		} else {
			file := filepath.Clean(path)
			fileName := filepath.Base(file)
			if excludeFilesMap[file] || filepath.Ext(file) != ".go" || fileName == *confProps.GeneratedFileName {
				continue
			}

			dir := filepath.Dir(file)
			if filesToScanMap[dir] == nil {
				filesToScanMap[dir] = make(map[string]struct{})
			}
			filesToScanMap[dir][file] = struct{}{}
		}

	}

	// process files
	for dir, files := range filesToScanMap {
		genCodeFile := filepath.Join(dir, *confProps.GeneratedFileName)
		removeFile(genCodeFile)

		filesArr := make([]string, 0, len(files))
		for file := range files {
			filesArr = append(filesArr, file)
		}
		code := genCode(filesArr)
		saveFile(genCodeFile, code)
	}
}

func genCode(files []string) string {
	analysis, err := extractStructsFromFilesInSamePackage(files)
	if err != nil {
		panic(err)
	}

	var codes string

	for _, structMap := range analysis.structs {
		for k, v := range structMap {
			wrapperName := *confProps.GeneratedStructPrefix + k + *confProps.GeneratedStructPostfix
			t := template.Must(template.New(wrapperName).Parse(wrapperStructTemplate))
			var b bytes.Buffer

			err := t.Execute(&b, struct {
				TypeName    string
				WrapperName string
				Fields      []map[string]string
			}{
				TypeName:    k,
				WrapperName: wrapperName,
				Fields:      v,
			})

			if err != nil {
				panic(err)
			}

			codes = fmt.Sprintf("%s\n%s", codes, b.String())
		}
	}

	// generate the header
	t := template.Must(template.New("header").Parse(wrapperHeaderTemplate))
	var b bytes.Buffer
	var imports []string
	for k := range analysis.imports {
		imports = append(imports, k)
	}
	err = t.Execute(&b, struct {
		PackageName string
		Imports     []string
		Content     string
	}{
		PackageName: analysis.packageName,
		Imports:     imports,
		Content:     codes,
	})

	if err != nil {
		panic(err)
	}

	return b.String()
}
