package tool

import (
	"cdt/internal"
)

type PyTest struct {
	internal.ExecutableTool
}

// DetectPyTest create a tool for pytest
func DetectPyTest(environment internal.Environment) *PyTest {
	return NewPyTest(func() *internal.Executable {
		return environment.FindExecutable("pytest")
	})
}

// NewPyTest creates a pytest tool from a custom executable
func NewPyTest(detect func() *internal.Executable) *PyTest {
	return &PyTest{
		ExecutableTool: internal.MakeExecutableTool(
			"pytest",
			"PyTest",
			"The pytest framework makes it easy to write small, readable tests, and can scale to support complex "+
				"functional testing for applications and libraries.",
			detect,
		),
	}
}

func (p *PyTest) TestAll(project internal.Project, args []string) error {
	return p.RunForProject(project, args)
}

func (p *PyTest) Test(project internal.Project, pattern string, args []string) error {
	return p.RunForProject(project, append(args, pattern))
}
