package tool

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"cdt/internal"
	"cdt/internal/utils"
)

const IdCMake = "cmake"

type CMake struct {
	internal.ExecutableTool
}

// NewCMake creates a cmake tool from a custom executable.
func NewCMake(detect internal.ExecutableToolDetectFunc) *CMake {
	return &CMake{
		internal.MakeExecutableTool(
			IdCMake,
			"CMake",
			"A Powerful Software Build System",
			internal.Tags{
				internal.ToolTagC,
				internal.ToolTagCpp,
				internal.ToolTagConfigure,
				internal.ToolTagBuild,
			},
			detect,
		),
	}
}

// DetectCMake create cmake tool can be used in the project.
func DetectCMake(
	ctx context.Context,
	options DetectOptions,
) *CMake {
	return NewCMake(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdCMake, "cmake"))
	})
}

func (c *CMake) Structure(
	ctx context.Context,
	info internal.ProjectInfo,
) (*internal.ProjectStructure, error) {
	if err := c.configureIfNeeded(ctx, internal.ProjectConfiguratorOptions{ProjectInfo: info}); err != nil {
		return nil, err
	}

	if info.OutputDirectory == nil {
		return nil, internal.ErrNoOutputDirectory
	}

	fileApi := utils.MakeCmakeFileApi(*info.OutputDirectory)

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
	if options.OutputDirectory == nil {
		return internal.ErrNoOutputDirectory
	}

	fileApi := utils.MakeCmakeFileApi(*options.OutputDirectory)

	if err := fileApi.Query("codemodel", 2); err != nil {
		return err
	}

	var callArgs []string
	callArgs = append(callArgs, "-DCMAKE_EXPORT_COMPILE_COMMANDS=ON")
	callArgs = append(callArgs, "-S", ".")
	callArgs = append(callArgs, "-B", *options.OutputDirectory)
	callArgs = append(callArgs, options.ExtraArgs...)

	return c.RunForProject(ctx, options.ProjectInfo, callArgs)
}

func (c *CMake) BuildTargets(
	ctx context.Context,
	options internal.ProjectBuilderOptions,
) error {
	if err := c.configureIfNeeded(
		ctx,
		internal.ProjectConfiguratorOptions{ProjectInfo: options.ProjectInfo},
	); err != nil {
		return err
	}

	if options.OutputDirectory == nil {
		return internal.ErrNoOutputDirectory
	}

	args := []string{"--build", *options.OutputDirectory}

	if options.Targets != nil && len(*options.Targets) > 0 {
		args = append(args, "--target")
		args = append(args, *options.Targets...)
	}

	args = append(args, options.ExtraArgs...)

	return c.RunForProject(ctx, options.ProjectInfo, args)
}

func (c *CMake) RunTarget(
	ctx context.Context,
	options internal.ProjectRunnerOptions,
	target string,
) error {
	buildOptions := internal.ProjectBuilderOptions{
		ProjectInfo: options.ProjectInfo,
		Targets:     &[]string{target},
	}

	if err := c.BuildTargets(ctx, buildOptions); err != nil {
		return err
	}

	if options.OutputDirectory == nil {
		return internal.ErrNoOutputDirectory
	}

	fileApi := utils.MakeCmakeFileApi(*options.OutputDirectory)

	reply, err := fileApi.Reply()
	if err != nil {
		return err
	}

	for _, t := range reply.Targets {
		if t.Name == target && t.Type == utils.TargetExecutable {
			return runTarget(ctx, options, t)
		}
	}

	return fmt.Errorf("target '%s' not found", target)
}

func (c *CMake) configureIfNeeded(
	ctx context.Context,
	options internal.ProjectConfiguratorOptions,
) error {
	if options.OutputDirectory != nil {
		if !internal.PathExists(filepath.Join(*options.OutputDirectory, "CMakeCache.txt")) {
			return c.Configure(ctx, options)
		}
	}

	return nil
}

func runTarget(
	ctx context.Context,
	options internal.ProjectRunnerOptions,
	target utils.ReplyTarget,
) error {
	executable := internal.Executable{
		Path:    filepath.Join(*options.OutputDirectory, target.Name),
		Runtime: options.Runtime,
	}

	return executable.Run(
		ctx,
		internal.RunOptions{
			Directory: options.Directory,
			Output:    os.Stdout,
			Error:     os.Stderr,
			Input:     os.Stdin,
		},
		options.ExtraArgs,
	)
}
