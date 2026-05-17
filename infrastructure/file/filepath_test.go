package file_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kakkky/hakoniwa/config"
	"github.com/kakkky/hakoniwa/infrastructure/file"
)

func TestNewFilePaths(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) (cfg *config.Config, wantDataDir, wantStateDir string)
	}{
		{
			name: "XDG_DATA_HOME / XDG_STATE_HOME が指定されているとそこに作成される",
			setup: func(t *testing.T) (*config.Config, string, string) {
				dataHome := t.TempDir()
				stateHome := t.TempDir()
				return &config.Config{
						XdgDataHome:  dataHome,
						XdgStateHome: stateHome,
					},
					filepath.Join(dataHome, file.AppName),
					filepath.Join(stateHome, file.AppName)
			},
		},
		{
			name: "XDG_* が空なら HOME 配下のデフォルトパスに作成される",
			setup: func(t *testing.T) (*config.Config, string, string) {
				home := t.TempDir()
				t.Setenv("HOME", home)
				return &config.Config{},
					filepath.Join(home, file.DefaultDataRel, file.AppName),
					filepath.Join(home, file.DefaultStateRel, file.AppName)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, wantDataDir, wantStateDir := tt.setup(t)

			paths, err := file.NewFilePaths(cfg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if paths == nil {
				t.Fatal("expected non-nil paths")
			}

			wantResidents := filepath.Join(wantDataDir, file.ResidentsFileName)
			wantBuilding := filepath.Join(wantDataDir, file.BuildingManagerFileName)
			wantLog := filepath.Join(wantStateDir, file.LogFileName)

			if paths.DataFilePaths.ResidentsFilePath != wantResidents {
				t.Errorf("ResidentsFilePath: got=%q want=%q", paths.DataFilePaths.ResidentsFilePath, wantResidents)
			}
			if paths.DataFilePaths.BuildingManagerFilePath != wantBuilding {
				t.Errorf("BuildingManagerFilePath: got=%q want=%q", paths.DataFilePaths.BuildingManagerFilePath, wantBuilding)
			}
			if paths.LogFilePath != wantLog {
				t.Errorf("LogFilePath: got=%q want=%q", paths.LogFilePath, wantLog)
			}

			for _, p := range []string{wantResidents, wantBuilding, wantLog} {
				if _, err := os.Stat(p); err != nil {
					t.Errorf("expected file to exist at %s: %v", p, err)
				}
			}
		})
	}
}

func TestNewFilePaths_Error(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) *config.Config
	}{
		{
			name: "XdgDataHome の親が「ファイル」だと MkdirAll がエラーになる",
			setup: func(t *testing.T) *config.Config {
				// regular file を作って、その配下にディレクトリを掘ろうとする
				tmp := t.TempDir()
				blocker := filepath.Join(tmp, "blocker")
				if err := os.WriteFile(blocker, []byte("x"), 0o600); err != nil {
					t.Fatalf("setup: %v", err)
				}
				return &config.Config{
					XdgDataHome:  blocker,
					XdgStateHome: t.TempDir(),
				}
			},
		},
		{
			name: "ファイル作成先ディレクトリが書き込み不可だと OpenFile がエラーになる",
			setup: func(t *testing.T) *config.Config {
				dataHome := t.TempDir()
				appDir := filepath.Join(dataHome, file.AppName)
				if err := os.MkdirAll(appDir, 0o500); err != nil {
					t.Fatalf("setup: %v", err)
				}
				t.Cleanup(func() { _ = os.Chmod(appDir, 0o700) })
				return &config.Config{
					XdgDataHome:  dataHome,
					XdgStateHome: t.TempDir(),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.setup(t)
			paths, err := file.NewFilePaths(cfg)
			if err == nil {
				t.Fatalf("expected error, got paths=%+v", paths)
			}
			if paths != nil {
				t.Errorf("expected nil paths, got %+v", paths)
			}
		})
	}
}
