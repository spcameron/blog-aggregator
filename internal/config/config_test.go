package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name        string
		contents    []byte
		expectedCfg Config
		wantErr     bool
	}{
		{
			name:     "happy",
			contents: []byte(`{"db_url":"testing_url","current_user_name":"testing_user_name"}`),
			expectedCfg: Config{
				URL:             "testing_url",
				CurrentUserName: "testing_user_name",
			},
			wantErr: false,
		},
		{
			name:        "missing_file",
			contents:    nil,
			expectedCfg: Config{},
			wantErr:     true,
		},
		{
			name:        "malformed_json",
			contents:    []byte("garbage"),
			expectedCfg: Config{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempHome := t.TempDir()
			t.Setenv("HOME", tempHome)
			tempFilepath := filepath.Join(tempHome, configFileName)

			if tt.contents != nil {
				err := os.WriteFile(tempFilepath, tt.contents, 0644)
				if err != nil {
					t.Fatalf("os.WriteFile() returned non-nil error, %v", err)
				}
			}

			cfg, err := Read()
			if err != nil && !tt.wantErr {
				t.Fatalf("Read() returned non-nil error, %v", err)
			} else if err == nil && tt.wantErr {
				t.Errorf("Expected error, but did not receive one")
			}

			if !tt.wantErr && cfg != tt.expectedCfg {
				t.Errorf("cfg does not match expected: want %+v, got %+v", tt.expectedCfg, cfg)
			}
		})
	}
}

func TestSetUser(t *testing.T) {
	tests := []struct {
		name        string
		newUser     string
		startingCfg Config
		expectedCfg Config
		wantErr     bool
	}{
		{
			name:    "updates_username",
			newUser: "new_name",
			startingCfg: Config{
				URL:             "testing_url",
				CurrentUserName: "old_name",
			},
			expectedCfg: Config{
				URL:             "testing_url",
				CurrentUserName: "new_name",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempHome := t.TempDir()
			t.Setenv("HOME", tempHome)
			tempFilepath := filepath.Join(tempHome, configFileName)

			file, err := os.Create(tempFilepath)
			if err != nil {
				t.Fatalf("Failed to create temp file, %v", err)
			}

			if err := json.NewEncoder(file).Encode(tt.startingCfg); err != nil {
				t.Fatalf("json.Encode() returned non-nil error, %v", err)
			}

			file.Close()

			if err := tt.startingCfg.SetUser(tt.newUser); err != nil && !tt.wantErr {
				t.Fatalf("SetUser() returned non-nil error, %v", err)
			}

			cfg, err := Read()
			if err != nil && !tt.wantErr {
				t.Fatalf("Read() returned non-nil error, %v", err)
			} else if err == nil && tt.wantErr {
				t.Errorf("Expected error, but did not receive one")
			}

			if !tt.wantErr && cfg != tt.expectedCfg {
				t.Errorf("cfg does not match expected: want %+v, got %+v", tt.expectedCfg, cfg)
			}
		})
	}
}
