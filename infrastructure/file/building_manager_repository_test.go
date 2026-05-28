package file_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kakkky/hakoniwa/config"
	"github.com/kakkky/hakoniwa/domain"
	"github.com/kakkky/hakoniwa/infrastructure/file"
	"github.com/stretchr/testify/assert"
)

func TestFileBuildingManagerRepository_Save(t *testing.T) {
	now := time.Date(2026, 5, 17, 0, 0, 0, 0, time.UTC)
	bm := domain.NewBuildingManager("管理人A", 50, now)

	tests := []struct {
		name   string
		seed   string
		toSave *domain.BuildingManager
	}{
		{
			name:   "空ファイルに保存できる",
			seed:   "",
			toSave: bm,
		},
		{
			name:   "既存内容を上書きできる",
			seed:   `{"Name":"先","Age":1,"AppointedAt":"2020-01-01T00:00:00Z"}`,
			toSave: bm,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			paths, err := file.NewFilePaths(&config.Config{XdgDataHome: dir, XdgStateHome: dir})
			if err != nil {
				t.Fatalf("setup: %v", err)
			}
			if err := os.WriteFile(paths.DataFilePaths.BuildingManagerFilePath, []byte(tt.seed), 0o600); err != nil {
				t.Fatalf("seed: %v", err)
			}
			repo := file.NewFileBuildingManagerRepository(paths)

			if err := repo.Save(tt.toSave); err != nil {
				t.Fatalf("Save: %v", err)
			}

			// 同じ repo で Get して書き込まれた内容が一致することを確認
			got, err := repo.Get()
			if err != nil {
				t.Fatalf("Get after Save: %v", err)
			}
			assert.Equal(t, *tt.toSave, got)
		})
	}
}

func TestFileBuildingManagerRepository_Save_Error(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) *file.FilePaths
	}{
		{
			name: "親ディレクトリが存在しないと WriteFile でエラー",
			setup: func(t *testing.T) *file.FilePaths {
				return &file.FilePaths{
					DataFilePaths: file.DataFilePaths{
						BuildingManagerFilePath: filepath.Join(t.TempDir(), "no_such_dir", "bm.json"),
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := tt.setup(t)
			repo := file.NewFileBuildingManagerRepository(paths)
			if err := repo.Save(domain.NewBuildingManager("x", 1, time.Now())); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestFileBuildingManagerRepository_Get(t *testing.T) {
	tests := []struct {
		name string
		seed string
		want domain.BuildingManager
	}{
		{
			name: "保存済み JSON が読み込める",
			seed: `{"Name":"管理人A","Age":50,"AppointedAt":"2026-05-17T00:00:00Z"}`,
			want: domain.BuildingManager{
				Name:        "管理人A",
				Age:         50,
				AppointedAt: time.Date(2026, 5, 17, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			paths, err := file.NewFilePaths(&config.Config{XdgDataHome: dir, XdgStateHome: dir})
			if err != nil {
				t.Fatalf("setup: %v", err)
			}
			if err := os.WriteFile(paths.DataFilePaths.BuildingManagerFilePath, []byte(tt.seed), 0o600); err != nil {
				t.Fatalf("seed: %v", err)
			}
			repo := file.NewFileBuildingManagerRepository(paths)

			got, err := repo.Get()
			if err != nil {
				t.Fatalf("Get: %v", err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFileBuildingManagerRepository_Get_Error(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T) *file.FilePaths
		wantErrIs error // nil なら error 発生のみ確認
	}{
		{
			name: "ファイルが空なら ErrBuildingManagerNotFound",
			setup: func(t *testing.T) *file.FilePaths {
				dir := t.TempDir()
				paths, err := file.NewFilePaths(&config.Config{XdgDataHome: dir, XdgStateHome: dir})
				if err != nil {
					t.Fatalf("setup: %v", err)
				}
				// NewFilePaths が作成したファイルは空のままなのでそのまま使う
				return paths
			},
			wantErrIs: domain.NotfoundErr,
		},
		{
			name: "ファイルが存在しないと ReadFile でエラー",
			setup: func(t *testing.T) *file.FilePaths {
				return &file.FilePaths{
					DataFilePaths: file.DataFilePaths{
						BuildingManagerFilePath: filepath.Join(t.TempDir(), "nonexistent.json"),
					},
				}
			},
		},
		{
			name: "不正な JSON だと Unmarshal エラー",
			setup: func(t *testing.T) *file.FilePaths {
				dir := t.TempDir()
				paths, err := file.NewFilePaths(&config.Config{XdgDataHome: dir, XdgStateHome: dir})
				if err != nil {
					t.Fatalf("setup: %v", err)
				}
				if err := os.WriteFile(paths.DataFilePaths.BuildingManagerFilePath, []byte("invalid json"), 0o600); err != nil {
					t.Fatalf("seed: %v", err)
				}
				return paths
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := tt.setup(t)
			repo := file.NewFileBuildingManagerRepository(paths)
			_, err := repo.Get()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Errorf("errors.Is mismatch: got=%v want=%v", err, tt.wantErrIs)
			}
		})
	}
}
