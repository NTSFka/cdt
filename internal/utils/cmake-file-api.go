package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CmakeFileApi struct {
	BuildDirectory string
}

// MakeCmakeFileApi create a struct for working with CMake File API.
func MakeCmakeFileApi(directory string) CmakeFileApi {
	return CmakeFileApi{BuildDirectory: directory}
}

func (c *CmakeFileApi) GetApiDirectory() string {
	return filepath.Join(c.BuildDirectory, ".cmake/api/v1/")
}

func (c *CmakeFileApi) GetQueryDirectory() string {
	return filepath.Join(c.GetApiDirectory(), "query")
}

func (c *CmakeFileApi) GetReplyDirectory() string {
	return filepath.Join(c.GetApiDirectory(), "reply")
}

// The Query creates a query to CMake.
func (c *CmakeFileApi) Query(kind string, version int) error {
	path := filepath.Join(c.GetQueryDirectory(), fmt.Sprintf("%s-v%d", kind, version))

	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return err
	}

	if _, err := os.Create(path); err != nil {
		return err
	}

	return nil
}

const (
	TargetUnsupported string = "unsupported"
	TargetExecutable  string = "executable"
	TargetLibrary     string = "library"
)

type ReplyTarget struct {
	Name     string
	Type     string
	Path     string
	External bool
	Files    []string
}

type Reply struct {
	Targets []ReplyTarget
}

func loadJSONFile(path string, target any) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, target); err != nil {
		return err
	}

	return nil
}

func convertTargetType(targetType string) string {
	if targetType == "EXECUTABLE" {
		return TargetExecutable
	}

	if targetType == "SHARED_LIBRARY" || targetType == "STATIC_LIBRARY" {
		return TargetLibrary
	}

	return TargetUnsupported
}

// Reply returns query reply after cmake run.
func (c *CmakeFileApi) Reply() (*Reply, error) {
	indices, err := filepath.Glob(filepath.Join(c.GetReplyDirectory(), "index-*.json"))

	if err != nil {
		return nil, err
	}

	if len(indices) == 0 {
		return nil, errors.New("no reply index file found")
	}

	var dataIndex cmakeFileApiIndex
	if err := loadJSONFile(indices[0], &dataIndex); err != nil {
		return nil, err
	}

	codeModelPath := filepath.Join(c.GetReplyDirectory(), dataIndex.Reply.CodemodelV2.JSONFile)

	var dataCodemodel cmakeFileApiCodemodelV2
	if err := loadJSONFile(codeModelPath, &dataCodemodel); err != nil {
		return nil, err
	}

	var targets []ReplyTarget

	for _, configuration := range dataCodemodel.Configurations {
		for _, target := range configuration.Targets {
			if replyTarget, err := c.processTarget(target); replyTarget != nil {
				targets = append(targets, *replyTarget)
			} else if err != nil {
				return nil, err
			}
		}
	}

	return &Reply{
		Targets: targets,
	}, nil
}

type cmakeFileApiIndex struct {
	Reply struct {
		CodemodelV2 struct {
			JSONFile string `json:"jsonFile"`
		} `json:"codemodel-v2"`
	} `json:"reply"`
}

type cmakeFileApiCodemodelV2 struct {
	Configurations []cmakeFileApiCodemodelV2Configuration `json:"configurations"`
}

type cmakeFileApiCodemodelV2Configuration struct {
	Targets []cmakeFileApiCodemodelV2ConfigurationTarget `json:"targets"`
}

type cmakeFileApiCodemodelV2ConfigurationTarget struct {
	Name     string `json:"name"`
	JSONFile string `json:"jsonFile"`
}

type cmakeFileApiTarget struct {
	Sources []struct {
		Path        string `json:"path"`
		IsGenerated bool   `json:"isGenerated"`
	} `json:"sources"`
	Type       string `json:"type"`
	NameOnDisk string `json:"nameOnDisk"`
	Paths      struct {
		Source string `json:"source"`
		Build  string `json:"build"`
	} `json:"paths"`
}

func (c *CmakeFileApi) processTarget(
	target cmakeFileApiCodemodelV2ConfigurationTarget,
) (*ReplyTarget, error) {
	path := filepath.Join(c.GetReplyDirectory(), target.JSONFile)

	var targetInfo cmakeFileApiTarget
	if err := loadJSONFile(path, &targetInfo); err != nil {
		return nil, err
	}

	// Ignore UTILITY targets
	if targetInfo.Type == "UTILITY" {
		return nil, nil // nolint: nilnil
	}

	// If the target source is in the build directory, it's a dependency
	dependency := strings.HasPrefix(targetInfo.Paths.Source, c.BuildDirectory)

	var files []string

	if !dependency {
		for _, source := range targetInfo.Sources {
			if source.IsGenerated {
				continue
			}

			files = append(files, source.Path)
		}
	}

	return &ReplyTarget{
		Name:     target.Name,
		Type:     convertTargetType(targetInfo.Type),
		External: dependency,
		Path:     targetInfo.NameOnDisk,
		Files:    files,
	}, nil
}
