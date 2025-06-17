package tool

import (
	. "cdt/internal"
	"cdt/internal/utils"
	"fmt"
	"path/filepath"
)

type CMake struct {
	ExecutableTool
}

// NewCMake creates a cmake tool from a custom executable
func NewCMake(executable *Executable) *CMake {
	return &CMake{
		MakeExecutableTool(
			"cmake",
			"CMake",
			"A Powerful Software Build System",
			executable,
		),
	}
}

// DetectCMake create cmake tool can be used in the project
func DetectCMake(environment Environment) *CMake {
	return NewCMake(environment.FindExecutable("cmake"))
}

func (c *CMake) Structure(project Project) (*ProjectStructure, error) {
	if err := c.Configure(project, []string{}); err != nil {
		return nil, err
	}

	fileApi := utils.NewCmakeFileApi(project.BuildDirectory())

	info := ProjectStructure{
		Targets: make(map[string]ProjectTarget),
	}

	if reply, err := fileApi.Reply(); err == nil {
		for _, target := range reply.Targets {
			info.Targets[target.Name] = ProjectTarget{
				Files:      target.Files,
				Dependency: target.External,
			}
		}
	}

	return &info, nil
}

func (c *CMake) Configure(project Project, args []string) error {
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

func (c *CMake) BuildAll(project Project, args []string) error {
	if err := c.Configure(project, []string{}); err != nil {
		return err
	}

	callArgs := args
	callArgs = append(callArgs, "--build", project.BuildDirectory())

	return c.Run(project, callArgs)
}

func (c *CMake) BuildTargets(project Project, targets []string, args []string) error {
	if err := c.Configure(project, []string{}); err != nil {
		return err
	}

	callArgs := args
	callArgs = append(callArgs, "--build", project.BuildDirectory())
	callArgs = append(callArgs, "--target")
	callArgs = append(callArgs, targets...)

	return c.Run(project, callArgs)
}

func (c *CMake) RunTarget(project Project, target string, args []string) error {
	if err := c.BuildAll(project, []string{}); err != nil {
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
			executable := Executable{
				Path:    filepath.Join(project.BuildDirectory(), t.Name),
				RunFunc: SystemEnvironment.RunExecutable,
			}

			return executable.Run(args)
		}
	}

	return fmt.Errorf("target '%s' not found", target)
}
