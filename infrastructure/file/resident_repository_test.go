package file_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/kakkky/hakoniwa/config"
	"github.com/kakkky/hakoniwa/domain"
	"github.com/kakkky/hakoniwa/infrastructure/file"
	"github.com/stretchr/testify/assert"
)

func mustNewResident(t *testing.T, name string, age int, gender domain.Gender, traits []domain.Trait) *domain.Resident {
	t.Helper()
	r, err := domain.NewResident(name, age, gender, traits)
	if err != nil {
		t.Fatalf("fixture NewResident: %v", err)
	}
	return r
}

func TestFileResidentRepository_Save(t *testing.T) {
	yamada := mustNewResident(t, "山田", 30, domain.Male, []domain.Trait{"優しい"})
	hanako := mustNewResident(t, "花子", 25, domain.Female, []domain.Trait{"明るい"})
	yamadaUpdated := *yamada
	yamadaUpdated.Age = 99

	tests := []struct {
		name    string
		seed    string
		toSave  *domain.Resident
		wantIDs []domain.ResidentID
	}{
		{
			name:    "空配列ファイルに 1 件追加 → 1 件になる",
			seed:    "[]",
			toSave:  yamada,
			wantIDs: []domain.ResidentID{yamada.ID},
		},
		{
			name:    "既存 1 件 → 別 ID を追加で 2 件になる",
			seed:    fmt.Sprintf(`[{"id":%q,"name":"花子","age":25,"gender":2,"traits":["明るい"]}]`, hanako.ID),
			toSave:  yamada,
			wantIDs: []domain.ResidentID{hanako.ID, yamada.ID},
		},
		{
			name:    "既存 1 件 → 同じ ID なら上書きで件数 1 のまま",
			seed:    fmt.Sprintf(`[{"id":%q,"name":"山田","age":30,"gender":1,"traits":["優しい"]}]`, yamada.ID),
			toSave:  &yamadaUpdated,
			wantIDs: []domain.ResidentID{yamada.ID},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			paths, err := file.NewFilePaths(&config.Config{XdgDataHome: dir, XdgStateHome: dir})
			if err != nil {
				t.Fatalf("setup: %v", err)
			}
			if err := os.WriteFile(paths.DataFilePaths.ResidentsFilePath, []byte(tt.seed), 0o600); err != nil {
				t.Fatalf("seed: %v", err)
			}
			repo := file.NewFileResidentRepository(paths)

			if err := repo.Save(tt.toSave); err != nil {
				t.Fatalf("Save: %v", err)
			}

			data, err := os.ReadFile(paths.DataFilePaths.ResidentsFilePath)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			var got []*domain.Resident
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			gotIDs := make([]domain.ResidentID, 0, len(got))
			for _, r := range got {
				gotIDs = append(gotIDs, r.ID)
			}
			assert.ElementsMatch(t, tt.wantIDs, gotIDs)
		})
	}
}

func TestFileResidentRepository_Save_Error(t *testing.T) {
	yamada := mustNewResident(t, "山田", 30, domain.Male, []domain.Trait{"優しい"})

	tests := []struct {
		name  string
		setup func(t *testing.T) *file.FilePaths
	}{
		{
			name: "ファイルが存在しないと ReadFile でエラー",
			setup: func(t *testing.T) *file.FilePaths {
				return &file.FilePaths{
					DataFilePaths: file.DataFilePaths{
						ResidentsFilePath: filepath.Join(t.TempDir(), "nonexistent.json"),
					},
				}
			},
		},
		{
			name: "不正な JSON だと Unmarshal でエラー",
			setup: func(t *testing.T) *file.FilePaths {
				dir := t.TempDir()
				paths, err := file.NewFilePaths(&config.Config{XdgDataHome: dir, XdgStateHome: dir})
				if err != nil {
					t.Fatalf("setup: %v", err)
				}
				if err := os.WriteFile(paths.DataFilePaths.ResidentsFilePath, []byte("invalid json"), 0o600); err != nil {
					t.Fatalf("seed: %v", err)
				}
				return paths
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := tt.setup(t)
			repo := file.NewFileResidentRepository(paths)
			if err := repo.Save(yamada); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestFileResidentRepository_SaveAll(t *testing.T) {
	yamada := mustNewResident(t, "山田", 30, domain.Male, []domain.Trait{"優しい"})
	hanako := mustNewResident(t, "花子", 25, domain.Female, []domain.Trait{"明るい"})

	tests := []struct {
		name    string
		seed    string
		toSave  []*domain.Resident
		wantIDs []domain.ResidentID
	}{
		{
			name:    "空配列で上書きすると保存内容も空配列になる",
			seed:    fmt.Sprintf(`[{"id":%q,"name":"先","age":1,"gender":1,"traits":["x"]}]`, yamada.ID),
			toSave:  []*domain.Resident{},
			wantIDs: []domain.ResidentID{},
		},
		{
			name:    "複数件を渡すと全件保存される (既存は上書き)",
			seed:    `[{"id":"old","name":"OLD","age":0,"gender":0,"traits":["x"]}]`,
			toSave:  []*domain.Resident{yamada, hanako},
			wantIDs: []domain.ResidentID{yamada.ID, hanako.ID},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			paths, err := file.NewFilePaths(&config.Config{XdgDataHome: dir, XdgStateHome: dir})
			if err != nil {
				t.Fatalf("setup: %v", err)
			}
			if err := os.WriteFile(paths.DataFilePaths.ResidentsFilePath, []byte(tt.seed), 0o600); err != nil {
				t.Fatalf("seed: %v", err)
			}
			repo := file.NewFileResidentRepository(paths)

			if err := repo.SaveAll(tt.toSave); err != nil {
				t.Fatalf("SaveAll: %v", err)
			}

			data, err := os.ReadFile(paths.DataFilePaths.ResidentsFilePath)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			var got []*domain.Resident
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			gotIDs := make([]domain.ResidentID, 0, len(got))
			for _, r := range got {
				gotIDs = append(gotIDs, r.ID)
			}
			assert.ElementsMatch(t, tt.wantIDs, gotIDs)
		})
	}
}

func TestFileResidentRepository_SaveAll_Error(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) *file.FilePaths
	}{
		{
			name: "親ディレクトリが存在しないと WriteFile でエラー",
			setup: func(t *testing.T) *file.FilePaths {
				return &file.FilePaths{
					DataFilePaths: file.DataFilePaths{
						ResidentsFilePath: filepath.Join(t.TempDir(), "no_such_dir", "residents.json"),
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := tt.setup(t)
			repo := file.NewFileResidentRepository(paths)
			if err := repo.SaveAll([]*domain.Resident{}); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestFileResidentRepository_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		seed    string
		wantIDs []domain.ResidentID
	}{
		{
			name:    "空配列なら 0 件",
			seed:    "[]",
			wantIDs: []domain.ResidentID{},
		},
		{
			name:    "2 件 seed すると 2 件返る",
			seed:    `[{"id":"r-1","name":"山田","age":30,"gender":1,"traits":["優しい"]},{"id":"r-2","name":"花子","age":25,"gender":2,"traits":["明るい"]}]`,
			wantIDs: []domain.ResidentID{"r-1", "r-2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			paths, err := file.NewFilePaths(&config.Config{XdgDataHome: dir, XdgStateHome: dir})
			if err != nil {
				t.Fatalf("setup: %v", err)
			}
			if err := os.WriteFile(paths.DataFilePaths.ResidentsFilePath, []byte(tt.seed), 0o600); err != nil {
				t.Fatalf("seed: %v", err)
			}
			repo := file.NewFileResidentRepository(paths)

			got, err := repo.GetAll()
			if err != nil {
				t.Fatalf("GetAll: %v", err)
			}
			gotIDs := make([]domain.ResidentID, 0, len(got))
			for _, r := range got {
				gotIDs = append(gotIDs, r.ID)
			}
			assert.ElementsMatch(t, tt.wantIDs, gotIDs)
		})
	}
}

func TestFileResidentRepository_GetAll_Error(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) *file.FilePaths
	}{
		{
			name: "ファイルが存在しないと ReadFile でエラー",
			setup: func(t *testing.T) *file.FilePaths {
				return &file.FilePaths{
					DataFilePaths: file.DataFilePaths{
						ResidentsFilePath: filepath.Join(t.TempDir(), "nonexistent.json"),
					},
				}
			},
		},
		{
			name: "不正な JSON だと Unmarshal でエラー",
			setup: func(t *testing.T) *file.FilePaths {
				dir := t.TempDir()
				paths, err := file.NewFilePaths(&config.Config{XdgDataHome: dir, XdgStateHome: dir})
				if err != nil {
					t.Fatalf("setup: %v", err)
				}
				if err := os.WriteFile(paths.DataFilePaths.ResidentsFilePath, []byte("invalid json"), 0o600); err != nil {
					t.Fatalf("seed: %v", err)
				}
				return paths
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := tt.setup(t)
			repo := file.NewFileResidentRepository(paths)
			if _, err := repo.GetAll(); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}
