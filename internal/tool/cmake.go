package tool

import (
	"cdt/internal"
	"cdt/internal/utils"
	"context"
	"fmt"
	"path/filepath"
)

type CMake struct {
	internal.ExecutableTool
}

// NewCMake creates a cmake tool from a custom executable
func NewCMake(detect func() *internal.Executable) *CMake {
	return &CMake{
		internal.MakeExecutableTool(
			"cmake",
			"CMake",
			"A Powerful Software Build System",
			internal.Tags{internal.ToolTagC, internal.ToolTagCpp, internal.ToolTagConfigure, internal.ToolTagBuild},
			detect,
		),
	}
}

// DetectCMake create cmake tool can be used in the project
func DetectCMake(environment internal.Environment) *CMake {
	return NewCMake(func() *internal.Executable {
		return environment.FindExecutable("cmake")
	})
}

func (c *CMake) Structure(info internal.ProjectInfo) (*internal.ProjectStructure, error) {
	if err := c.Configure(info, []string{}); err != nil {
		return nil, err
	}

	if info.IntermediateDirectory == nil {
		return nil, internal.ErrNoIntermediateDirectory
	}

	fileApi := utils.NewCmakeFileApi(*info.IntermediateDirectory)

	structure := internal.ProjectStructure{
		Targets: make(map[string]internal.ProjectTarget),
	}

	if reply, err := fileApi.Reply(); err == nil {
		for _, target := range reply.Targets {
			structure.Targets[target.Name] = internal.ProjectTarget{
				Files:      target.Files,
				Dependency: target.External,
			}
		}
	}

	return &structure, nil
}

func (c *CMake) Configure(info internal.ProjectInfo, args []string) error {
	if info.IntermediateDirectory == nil {
		return internal.ErrNoIntermediateDirectory
	}

	fileApi := utils.NewCmakeFileApi(*info.IntermediateDirectory)

	if err := fileApi.Query("codemodel", 2); err != nil {
		return err
	}

	callArgs := args
	callArgs = append(callArgs, "-DCMAKE_EXPORT_COMPILE_COMMANDS=ON")
	callArgs = append(callArgs, "-S", ".")
	callArgs = append(callArgs, "-B", *info.IntermediateDirectory)

	return c.RunForProject(info, callArgs)
}

func (c *CMake) BuildAll(info internal.ProjectInfo, args []string) error {
	if err := c.Configure(info, []string{}); err != nil {
		return err
	}

	if info.IntermediateDirectory == nil {
		return internal.ErrNoIntermediateDirectory
	}

	callArgs := args
	callArgs = append(callArgs, "--build", *info.IntermediateDirectory)

	return c.RunForProject(info, callArgs)
}

func (c *CMake) BuildTargets(info internal.ProjectInfo, targets []string, args []string) error {
	if err := c.Configure(info, []string{}); err != nil {
		return err
	}

	if info.IntermediateDirectory == nil {
		return internal.ErrNoIntermediateDirectory
	}

	callArgs := args
	callArgs = append(callArgs, "--build", *info.IntermediateDirectory)
	callArgs = append(callArgs, "--target")
	callArgs = append(callArgs, targets...)

	return c.RunForProject(info, callArgs)
}

func (c *CMake) RunTarget(info internal.ProjectInfo, target string, args []string) error {
	if err := c.BuildTargets(info, []string{target}, []string{}); err != nil {
		return err
	}

	if info.IntermediateDirectory == nil {
		return internal.ErrNoIntermediateDirectory
	}

	fileApi := utils.NewCmakeFileApi(*info.IntermediateDirectory)

	reply, err := fileApi.Reply()
	if err != nil {
		return err
	}

	for _, t := range reply.Targets {
		if t.Name == target && t.Type == utils.TargetExecutable {
			// TODO: run environment?
			executable := internal.Executable{
				Path:    filepath.Join(*info.IntermediateDirectory, t.Name),
				Runtime: internal.SystemEnvironment,
			}

			return executable.Run(context.Background(), internal.RunOptions{Directory: info.Directory}, args)
		}
	}

	return fmt.Errorf("target '%s' not found", target)
}
