package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		confName   string
		wantConfig *Config
		wantErr    bool
	}{
		{
			name: "local config",
			path: "./",
			wantConfig: &Config{
				serverAddress: ":8080",
				baseURL:       "http://localhost:8080/",
				filePath:      "../urls.log",
				dbAddress:     "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfig, err := LoadConfig(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotConfig, tt.wantConfig) {
				t.Errorf("LoadConfig() gotConfig = %v, want %v", gotConfig, tt.wantConfig)
			}
		})
	}
}
