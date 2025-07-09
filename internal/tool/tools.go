package tool

import "cdt/internal"

type SupportedTools struct {
	ClangFormat  *ClangFormat
	ClangTidy    *ClangTidy
	CMake        *CMake
	CTest        *CTest
	Go           *Go
	GolangCILint *GolangCILint
}

func (t *SupportedTools) ToTools() internal.Tools {
	return internal.Tools{
		t.ClangFormat,
		t.ClangTidy,
		t.CMake,
		t.CTest,
		t.Go,
		t.GolangCILint,
	}
}
