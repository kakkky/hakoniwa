package file

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/kakkky/hakoniwa/domain"
)

var ErrBuildingManagerNotFound = errors.New("building manager not found")

type FileBuildingManagerRepository struct {
	filePath string
}

var _ domain.BuildingManagerRepository = (*FileBuildingManagerRepository)(nil)

func NewFileBuildingManagerRepository(paths *FilePaths) *FileBuildingManagerRepository {
	return &FileBuildingManagerRepository{
		filePath: paths.DataFilePaths.BuildingManagerFilePath,
	}
}

func (repo *FileBuildingManagerRepository) Save(bm *domain.BuildingManager) error {
	data, err := json.Marshal(bm)
	if err != nil {
		return err
	}
	return os.WriteFile(repo.filePath, data, 0o600)
}

func (repo *FileBuildingManagerRepository) Get() (domain.BuildingManager, error) {
	data, err := os.ReadFile(repo.filePath)
	if err != nil {
		return domain.BuildingManager{}, err
	}
	if len(data) == 0 {
		return domain.BuildingManager{}, ErrBuildingManagerNotFound
	}
	var bm domain.BuildingManager
	if err := json.Unmarshal(data, &bm); err != nil {
		return domain.BuildingManager{}, err
	}
	return bm, nil
}
