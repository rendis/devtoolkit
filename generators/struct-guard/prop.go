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

var confProps *StructGuardProp

type GeneratorConfProp struct {
	*GeneratorsProp `yaml:"generators" validate:"required"`
}

type GeneratorsProp struct {
	StructGuard *StructGuardProp `yaml:"struct-guard" validate:"required"`
}

type StructGuardProp struct {
	GeneratedFileName      *string  `yaml:"generated-file-name"`
	GeneratedStructPrefix  *string  `yaml:"generated-struct-prefix"`
	GeneratedStructPostfix *string  `yaml:"generated-struct-postfix"`
	ToScan                 []string `yaml:"to-scan"`
	ExcludeFilesToScan     []string `yaml:"exclude-files-to-scan"`
}

func (p *GeneratorConfProp) SetDefaults() {
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
	p := &GeneratorConfProp{}
	var props = []devtoolkit.ToolKitProp{p}

	if err := devtoolkit.LoadPropFile(propFilePath, props); err != nil {
		log.Fatalf("failed to load prop file '%s'.\n%v", propFilePath, err)
	}

	confProps = p.GeneratorsProp.StructGuard
}
