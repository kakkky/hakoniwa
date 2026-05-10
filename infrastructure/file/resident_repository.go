package file

import (
	"encoding/json"
	"os"

	"github.com/kakkky/hakoniwa/domain"
)

type FileResidentRepository struct {
	filePath string
}

var _ domain.ResidentRepository = (*FileResidentRepository)(nil)

func NewFileResidentRepository(filePath string) *FileResidentRepository {
	return &FileResidentRepository{
		filePath: filePath,
	}
}

func (repo *FileResidentRepository) Save(resident *domain.Resident) error {
	// 既存の住民データを読み込む
	residents := make([]*domain.Resident, 0)
	data, err := os.ReadFile(repo.filePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &residents); err != nil {
		return err
	}
	residentsMap := make(map[string]*domain.Resident, len(residents))
	for _, r := range residents {
		residentsMap[r.ID] = r
	}

	// 新たな住民を追加 or 更新
	residentsMap[resident.ID] = resident
	residents = make([]*domain.Resident, 0, len(residentsMap))
	for _, r := range residentsMap {
		residents = append(residents, r)
	}

	data, err = json.Marshal(residents)
	if err != nil {
		return err
	}
	if err := os.WriteFile(repo.filePath, data, 0o644); err != nil {
		return err
	}
	return nil
}
