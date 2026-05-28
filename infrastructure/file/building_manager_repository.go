package file

import (
	"encoding/json"
	"os"

	"github.com/kakkky/hakoniwa/domain"
)

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

func (repo *FileBuildingManagerRepository) Get() (*domain.BuildingManager, error) {
	data, err := os.ReadFile(repo.filePath)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, domain.NotfoundErr.With("building manager not found")
	}
	var bm domain.BuildingManager
	if err := json.Unmarshal(data, &bm); err != nil {
		return nil, err
	}
	return &bm, nil
}
