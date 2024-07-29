package main

import (
	"github.com/rendis/devtoolkit"
	"log"
	"path/filepath"
)

const (
	propFilePath  = "devtoolkit.yml"
	defaultGoFile = "codegen.go"
)

var generatorProp *StructGuardProp

type GeneratorsConfProp struct {
	*GeneratorsProp `yaml:"generators" validate:"required"`
}

type GeneratorsProp struct {
	StructGuard *StructGuardProp `yaml:"struct-guard" validate:"required"`
}

type StructGuardProp struct {
	// GeneratedFileName is the name of the generated file, defaults to 'codegen.go'
	GeneratedFileName *string `yaml:"generated-file-name"`

	// GeneratedStructPrefix is the prefix to be added to the generated struct name, defaults to ''
	GeneratedStructPrefix *string `yaml:"generated-struct-prefix"`

	// GeneratedStructPostfix is the postfix to be added to the generated struct name, defaults to 'Wrapper'
	GeneratedStructPostfix *string `yaml:"generated-struct-postfix"`

	// ToScan is the list of directories or files to scan for structs
	ToScan []string `yaml:"to-scan"`

	// ExcludeFilesToScan is the list of files to exclude from scanning
	ExcludeFilesToScan []string `yaml:"exclude-files-to-scan"`

	// ForceExport is a flag to force export of the generated struct, defaults to false (private)
	ForceExport bool `yaml:"force-export"`
}

func (p *GeneratorsConfProp) SetDefaults() {
	if p.GeneratorsProp == nil {
		p.GeneratorsProp = &GeneratorsProp{}
	}

	if p.GeneratorsProp.StructGuard == nil {
		p.GeneratorsProp.StructGuard = &StructGuardProp{}
	}

	// set defaults
	p.GeneratorsProp.StructGuard.SetDefaults()
}

func (p *StructGuardProp) SetDefaults() {
	if p.GeneratedFileName == nil {
		p.GeneratedFileName = devtoolkit.ToPtr(defaultGoFile)
	} else {
		if ext := filepath.Ext(*p.GeneratedFileName); ext != ".go" {
			*p.GeneratedFileName = *p.GeneratedFileName + ".go"
		}
	}

	if p.GeneratedStructPrefix == nil {
		p.GeneratedStructPrefix = devtoolkit.ToPtr("")
	}

	if p.GeneratedStructPostfix == nil {
		p.GeneratedStructPostfix = devtoolkit.ToPtr("Wrapper")
	}
}

func loadGenProp() {
	p := &GeneratorsConfProp{}
	var props = []devtoolkit.ToolKitProp{p}

	if err := devtoolkit.LoadPropFile(propFilePath, props); err != nil {
		log.Fatalf("failed to load prop file '%s'.\n%v", propFilePath, err)
	}

	generatorProp = p.GeneratorsProp.StructGuard
}
