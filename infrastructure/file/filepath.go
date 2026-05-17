package file

import (
	"os"
	"path/filepath"

	"github.com/kakkky/hakoniwa/config"
)

const AppName = "hakoniwa"

const (
	DefaultDataRel  = ".local/share"
	DefaultStateRel = ".local/state"
)

const (
	ResidentsFileName       = "residents.json"
	BuildingManagerFileName = "building_manager.json"
	LogFileName             = "hakoniwa.log"
)

const (
	dirMode  os.FileMode = 0o700
	fileMode os.FileMode = 0o600
)

type FilePaths struct {
	DataFilePaths DataFilePaths
	LogFilePath   string
}

type DataFilePaths struct {
	ResidentsFilePath       string
	BuildingManagerFilePath string
}

func NewFilePaths(cfg *config.Config) (*FilePaths, error) {
	dataDir, stateDir, err := resolveDirs(cfg)
	if err != nil {
		return nil, err
	}

	for _, d := range []string{dataDir, stateDir} {
		if err := os.MkdirAll(d, dirMode); err != nil {
			return nil, err
		}
	}

	residentsFilePath := filepath.Join(dataDir, ResidentsFileName)
	buildingManagerFilePath := filepath.Join(dataDir, BuildingManagerFileName)
	logFilePath := filepath.Join(stateDir, LogFileName)

	for _, p := range []string{residentsFilePath, buildingManagerFilePath, logFilePath} {
		f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, fileMode)
		if err != nil {
			return nil, err
		}
		if err := f.Close(); err != nil {
			return nil, err
		}
	}

	return &FilePaths{
		DataFilePaths: DataFilePaths{
			ResidentsFilePath:       residentsFilePath,
			BuildingManagerFilePath: buildingManagerFilePath,
		},
		LogFilePath: logFilePath,
	}, nil
}

func resolveDirs(cfg *config.Config) (dataDir, stateDir string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	dataDir = resolveXDGPath(cfg.XdgDataHome, home, DefaultDataRel)
	stateDir = resolveXDGPath(cfg.XdgStateHome, home, DefaultStateRel)
	return dataDir, stateDir, nil
}

func resolveXDGPath(xdgValue, home, defaultRel string) string {
	if xdgValue != "" {
		return filepath.Join(xdgValue, AppName)
	}
	return filepath.Join(home, defaultRel, AppName)
}
