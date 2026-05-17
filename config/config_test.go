package config_test

import (
	"os"
	"testing"

	"github.com/kakkky/hakoniwa/config"
)

// unsetEnv removes an env var for the duration of the test (restored on cleanup).
// t.Setenv だけでは空文字を「設定済み」と扱われ required を弾けないため。
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	// t.Setenv はクリーンアップで元の値に戻してくれるので、まず一旦これで登録してから unset する
	t.Setenv(key, "placeholder")
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("failed to unset %s: %v", key, err)
	}
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name         string
		envDataHome  string
		envStateHome string
		envAPIKey    string
	}{
		{
			name:         "全 env が揃っていれば値が反映される",
			envDataHome:  "/tmp/data",
			envStateHome: "/tmp/state",
			envAPIKey:    "test-key",
		},
		{
			name:         "XDG_DATA_HOME / XDG_STATE_HOME は空でもよい (required ではない)",
			envDataHome:  "",
			envStateHome: "",
			envAPIKey:    "another-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("XDG_DATA_HOME", tt.envDataHome)
			t.Setenv("XDG_STATE_HOME", tt.envStateHome)
			t.Setenv("GEMINI_API_KEY", tt.envAPIKey)

			cfg, err := config.NewConfig()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg == nil {
				t.Fatal("expected non-nil config")
			}
			if cfg.XdgDataHome != tt.envDataHome {
				t.Errorf("XdgDataHome: got=%q want=%q", cfg.XdgDataHome, tt.envDataHome)
			}
			if cfg.XdgStateHome != tt.envStateHome {
				t.Errorf("XdgStateHome: got=%q want=%q", cfg.XdgStateHome, tt.envStateHome)
			}
			if cfg.GeminiAPIKey != tt.envAPIKey {
				t.Errorf("GeminiAPIKey: got=%q want=%q", cfg.GeminiAPIKey, tt.envAPIKey)
			}
		})
	}
}

func TestNewConfig_Error(t *testing.T) {
	tests := []struct {
		name     string
		setupEnv func(t *testing.T)
	}{
		{
			name: "GEMINI_API_KEY が未設定ならエラー",
			setupEnv: func(t *testing.T) {
				unsetEnv(t, "GEMINI_API_KEY")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv(t)
			cfg, err := config.NewConfig()
			if err == nil {
				t.Fatalf("expected error, got cfg=%+v", cfg)
			}
			if cfg != nil {
				t.Errorf("expected nil cfg, got %+v", cfg)
			}
		})
	}
}
