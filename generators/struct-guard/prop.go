package main

import (
	"github.com/rendis/devtoolkit"
	"log"
	"path/filepath"
)

const (
	propFilePath  = "codegen.yml"
	defaultGoFile = "codegen.go"
)

var confProps *GeneratorConfProp

type GeneratorConfProp struct {
	*GeneratorProp `yaml:"code-generator" validate:"required"`
}

type GeneratorProp struct {
	GeneratedFileName      *string  `yaml:"generated-file-name"`
	GeneratedStructPrefix  *string  `yaml:"generated-struct-prefix"`
	GeneratedStructPostfix *string  `yaml:"generated-struct-postfix"`
	ToScan                 []string `yaml:"to-scan"`
	ExcludeFilesToScan     []string `yaml:"exclude-files-to-scan"`
}

func (p *GeneratorConfProp) SetDefaults() {
	if p.GeneratorProp == nil {
		p.GeneratorProp = &GeneratorProp{}
	}

	if p.GeneratorProp.GeneratedFileName == nil {
		p.GeneratorProp.GeneratedFileName = devtoolkit.ToPtr(defaultGoFile)
	} else {
		if ext := filepath.Ext(*p.GeneratorProp.GeneratedFileName); ext != ".go" {
			*p.GeneratorProp.GeneratedFileName = *p.GeneratorProp.GeneratedFileName + ".go"
		}
	}

	if p.GeneratorProp.GeneratedStructPrefix == nil {
		p.GeneratorProp.GeneratedStructPrefix = devtoolkit.ToPtr("")
	}

	if p.GeneratorProp.GeneratedStructPostfix == nil {
		p.GeneratorProp.GeneratedStructPostfix = devtoolkit.ToPtr("Wrapper")
	}
}

func loadProp() {
	confProps = &GeneratorConfProp{}
	var props = []devtoolkit.ToolKitProp{
		confProps,
	}

	if err := devtoolkit.LoadPropFile(propFilePath, props); err != nil {
		log.Fatalf("failed to load prop file '%s'.\n%v", propFilePath, err)
	}
}
