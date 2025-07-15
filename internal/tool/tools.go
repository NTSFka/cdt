package tool

import (
	"cdt/internal"
	"reflect"
)

// SupportedTools stores all supported tools
type SupportedTools struct {
	ClangFormat  *ClangFormat
	ClangTidy    *ClangTidy
	CMake        *CMake
	CTest        *CTest
	Go           *Go
	GolangCILint *GolangCILint
}

// NewSupportedTools initializes all supported tools for a given environment
func NewSupportedTools(environment internal.Environment) SupportedTools {
	return SupportedTools{
		ClangFormat:  DetectClangFormat(environment, nil),
		ClangTidy:    DetectClangTidy(environment, nil),
		CMake:        DetectCMake(environment),
		CTest:        DetectCTest(environment),
		Go:           DetectGo(environment),
		GolangCILint: DetectGolangCILint(environment),
	}
}

// ToTools convert supported tools to a list of generic tools
func (t *SupportedTools) ToTools() (result internal.Tools) {
	add := func(tool internal.Tool) {
		if !reflect.ValueOf(tool).IsNil() {
			result = append(result, tool)
		}
	}

	add(t.ClangFormat)
	add(t.ClangTidy)
	add(t.CMake)
	add(t.CTest)
	add(t.Go)
	add(t.GolangCILint)

	return
}
