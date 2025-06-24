package tool

import (
	"cdt/internal"
	"cdt/internal/utils"
	"fmt"
	"path/filepath"
)

type CMake struct {
	internal.ExecutableTool
}

// NewCMake creates a cmake tool from a custom executable
func NewCMake(executable *internal.Executable) *CMake {
	return &CMake{
		internal.MakeExecutableTool(
			"cmake",
			"CMake",
			"A Powerful Software Build System",
			executable,
		),
	}
}

// DetectCMake create cmake tool can be used in the project
func DetectCMake(environment internal.Environment) *CMake {
	return NewCMake(environment.FindExecutable("cmake"))
}

func (c *CMake) Structure(project internal.Project) (*internal.ProjectStructure, error) {
	if err := c.Configure(project, []string{}); err != nil {
		return nil, err
	}

	fileApi := utils.NewCmakeFileApi(project.BuildDirectory())

	info := internal.ProjectStructure{
		Targets: make(map[string]internal.ProjectTarget),
	}

	if reply, err := fileApi.Reply(); err == nil {
		for _, target := range reply.Targets {
			info.Targets[target.Name] = internal.ProjectTarget{
				Files:      target.Files,
				Dependency: target.External,
			}
		}
	}

	return &info, nil
}

func (c *CMake) Configure(project internal.Project, args []string) error {
	fileApi := utils.NewCmakeFileApi(project.BuildDirectory())

	if err := fileApi.Query("codemodel", 2); err != nil {
		return err
	}

	callArgs := args
	callArgs = append(callArgs, "-DCMAKE_EXPORT_COMPILE_COMMANDS=ON")
	callArgs = append(callArgs, "-S", project.RootDirectory())
	callArgs = append(callArgs, "-B", project.BuildDirectory())

	return c.Run(project, callArgs)
}

func (c *CMake) BuildAll(project internal.Project, args []string) error {
	if err := c.Configure(project, []string{}); err != nil {
		return err
	}

	callArgs := args
	callArgs = append(callArgs, "--build", project.BuildDirectory())

	return c.Run(project, callArgs)
}

func (c *CMake) BuildTargets(project internal.Project, targets []string, args []string) error {
	if err := c.Configure(project, []string{}); err != nil {
		return err
	}

	callArgs := args
	callArgs = append(callArgs, "--build", project.BuildDirectory())
	callArgs = append(callArgs, "--target")
	callArgs = append(callArgs, targets...)

	return c.Run(project, callArgs)
}

func (c *CMake) RunTarget(project internal.Project, target string, args []string) error {
	if err := c.BuildTargets(project, []string{target}, []string{}); err != nil {
		return err
	}

	fileApi := utils.NewCmakeFileApi(project.BuildDirectory())

	reply, err := fileApi.Reply()
	if err != nil {
		return err
	}

	for _, t := range reply.Targets {
		if t.Name == target && t.Type == utils.TargetExecutable {
			// TODO: run environment?
			executable := internal.Executable{
				Path:    filepath.Join(project.BuildDirectory(), t.Name),
				RunFunc: internal.SystemEnvironment.RunExecutable,
			}

			return executable.Run(args)
		}
	}

	return fmt.Errorf("target '%s' not found", target)
}
