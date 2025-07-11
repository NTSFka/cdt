package tool

import (
	"cdt/internal"
	"reflect"
)

type SupportedTools struct {
	ClangFormat  *ClangFormat
	ClangTidy    *ClangTidy
	CMake        *CMake
	CTest        *CTest
	Go           *Go
	GolangCILint *GolangCILint
}

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
