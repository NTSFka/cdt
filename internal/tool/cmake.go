package tool

import (
	. "cdt/internal"
	"cdt/internal/utils"
	"errors"
	"fmt"
	"path/filepath"
)

type CMake struct {
	executable *Executable
}

// NewCMake creates a cmake tool from a custom executable
func NewCMake(executable *Executable) *CMake {
	return &CMake{
		executable: executable,
	}
}

// DetectCMake create cmake tool can be used in the project
func DetectCMake() *CMake {
	return NewCMake(FindExecutable("cmake"))
}

func (c *CMake) Id() string {
	return "cmake"
}

func (c *CMake) Name() string {
	return "CMake"
}

func (c *CMake) Info() string {
	return "A Powerful Software Build System"
}

func (c *CMake) ExecutablePath() *string {
	if c.executable != nil {
		return &c.executable.Path
	}

	return nil
}

func (c *CMake) IsAvailable() bool {
	return c.executable != nil
}

func (c *CMake) Run(_ Project, args []string) error {
	if c.executable == nil {
		return errors.New("CMake is not installed on the system")
	}

	return c.executable.Run(args)
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

	return c.executable.Run(callArgs)
}

func (c *CMake) BuildAll(project Project, args []string) error {
	if err := c.Configure(project, []string{}); err != nil {
		return err
	}

	callArgs := args
	callArgs = append(callArgs, "--build", project.BuildDirectory())

	return c.executable.Run(callArgs)
}

func (c *CMake) BuildTargets(project Project, targets []string, args []string) error {
	if err := c.Configure(project, []string{}); err != nil {
		return err
	}

	callArgs := args
	callArgs = append(callArgs, "--build", project.BuildDirectory())
	callArgs = append(callArgs, "--target")
	callArgs = append(callArgs, targets...)

	return c.executable.Run(callArgs)
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
			executable := Executable{Path: filepath.Join(project.BuildDirectory(), t.Name)}

			return executable.Run(args)
		}
	}

	return fmt.Errorf("target '%s' not found", target)
}
