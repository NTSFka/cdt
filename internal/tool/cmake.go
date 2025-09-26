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

// NewCMake creates a cmake tool from a custom executable.
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

// DetectCMake create cmake tool can be used in the project.
func DetectCMake(ctx context.Context, environment internal.Environment) *CMake {
	return NewCMake(func() *internal.Executable {
		return environment.FindExecutable(ctx, "cmake")
	})
}

func (c *CMake) Structure(ctx context.Context, info internal.ProjectInfo) (*internal.ProjectStructure, error) {
	if err := c.Configure(ctx, internal.ProjectConfiguratorOptions{ProjectInfo: info}); err != nil {
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

func (c *CMake) Configure(ctx context.Context, options internal.ProjectConfiguratorOptions) error {
	if options.IntermediateDirectory == nil {
		return internal.ErrNoIntermediateDirectory
	}

	fileApi := utils.NewCmakeFileApi(*options.IntermediateDirectory)

	if err := fileApi.Query("codemodel", 2); err != nil {
		return err
	}

	var callArgs []string
	callArgs = append(callArgs, "-DCMAKE_EXPORT_COMPILE_COMMANDS=ON")
	callArgs = append(callArgs, "-S", ".")
	callArgs = append(callArgs, "-B", *options.IntermediateDirectory)
	callArgs = append(callArgs, options.ExtraArgs...)

	return c.RunForProject(ctx, options.ProjectInfo, callArgs)
}

func (c *CMake) BuildAll(ctx context.Context, options internal.ProjectBuilderOptions) error {
	if err := c.Configure(ctx, internal.ProjectConfiguratorOptions{ProjectInfo: options.ProjectInfo}); err != nil {
		return err
	}

	if options.IntermediateDirectory == nil {
		return internal.ErrNoIntermediateDirectory
	}

	callArgs := []string{"--build", *options.IntermediateDirectory}
	callArgs = append(callArgs, options.ExtraArgs...)

	return c.RunForProject(ctx, options.ProjectInfo, callArgs)
}

func (c *CMake) BuildTargets(ctx context.Context, options internal.ProjectBuilderOptions, targets []string) error {
	if err := c.Configure(ctx, internal.ProjectConfiguratorOptions{ProjectInfo: options.ProjectInfo}); err != nil {
		return err
	}

	if options.IntermediateDirectory == nil {
		return internal.ErrNoIntermediateDirectory
	}

	callArgs := []string{"--build", *options.IntermediateDirectory}
	callArgs = append(callArgs, "--target")
	callArgs = append(callArgs, targets...)
	callArgs = append(callArgs, options.ExtraArgs...)

	return c.RunForProject(ctx, options.ProjectInfo, callArgs)
}

func (c *CMake) RunTarget(ctx context.Context, options internal.ProjectRunnerOptions, target string) error {
	if err := c.BuildTargets(ctx, internal.ProjectBuilderOptions{ProjectInfo: options.ProjectInfo}, []string{target}); err != nil {
		return err
	}

	if options.IntermediateDirectory == nil {
		return internal.ErrNoIntermediateDirectory
	}

	fileApi := utils.NewCmakeFileApi(*options.IntermediateDirectory)

	reply, err := fileApi.Reply()
	if err != nil {
		return err
	}

	for _, t := range reply.Targets {
		if t.Name == target && t.Type == utils.TargetExecutable {
			// TODO: run environment?
			executable := internal.Executable{
				Path:    filepath.Join(*options.IntermediateDirectory, t.Name),
				Runtime: internal.SystemEnvironment,
			}

			return executable.Run(ctx, internal.RunOptions{Directory: options.Directory}, options.ExtraArgs)
		}
	}

	return fmt.Errorf("target '%s' not found", target)
}
