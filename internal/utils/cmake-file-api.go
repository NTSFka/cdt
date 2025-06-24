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

// NewCmakeFileApi create a struct for working with CMake File API
func NewCmakeFileApi(directory string) CmakeFileApi {
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

// The Query creates a query to CMake
func (c *CmakeFileApi) Query(kind string, version int) error {
	path := filepath.Join(c.GetQueryDirectory(), fmt.Sprintf("%s-v%d", kind, version))

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
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

func loadJsonFile(path string, target any) error {
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

// The Reply returns query reply after cmake run
func (c *CmakeFileApi) Reply() (*Reply, error) {
	indices, err := filepath.Glob(filepath.Join(c.GetReplyDirectory(), "index-*.json"))

	if err != nil {
		return nil, err
	}

	if len(indices) == 0 {
		return nil, errors.New("no reply index file found")
	}

	var dataIndex struct {
		Reply struct {
			CodemodelV2 struct {
				JsonFile string `json:"jsonFile"`
			} `json:"codemodel-v2"`
		} `json:"reply"`
	}

	if err := loadJsonFile(indices[0], &dataIndex); err != nil {
		return nil, err
	}

	var dataCodemodel struct {
		Configurations []struct {
			Targets []struct {
				Name     string `json:"name"`
				JsonFile string `json:"jsonFile"`
			} `json:"targets"`
		} `json:"configurations"`
	}

	codeModelPath := filepath.Join(c.GetReplyDirectory(), dataIndex.Reply.CodemodelV2.JsonFile)

	if err := loadJsonFile(codeModelPath, &dataCodemodel); err != nil {
		return nil, err
	}

	var targets []ReplyTarget

	for _, configuration := range dataCodemodel.Configurations {
		for _, target := range configuration.Targets {
			path := filepath.Join(c.GetReplyDirectory(), target.JsonFile)

			var targetInfo struct {
				Sources []struct {
					Path string `json:"path"`
				} `json:"sources"`
				Type       string `json:"type"`
				NameOnDisk string `json:"nameOnDisk"`
				Paths      struct {
					Source string `json:"source"`
					Build  string `json:"build"`
				} `json:"paths"`
			}

			if err := loadJsonFile(path, &targetInfo); err != nil {
				return nil, err
			}

			// If the target source is in build directory, it's dependency
			dependency := strings.HasPrefix(targetInfo.Paths.Source, c.BuildDirectory)

			var files []string

			if !dependency {
				for _, source := range targetInfo.Sources {
					files = append(files, source.Path)
				}
			}

			targets = append(targets, ReplyTarget{
				Name:     target.Name,
				Type:     convertTargetType(targetInfo.Type),
				External: dependency,
				Path:     targetInfo.NameOnDisk,
				Files:    files,
			})
		}
	}

	return &Reply{
		Targets: targets,
	}, nil
}
