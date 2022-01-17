package lib

import (
	"testing"
)

func TestConfig_AddFile(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		config  *Config
		args    args
		wantErr bool
	}{
		{
			name:    "Add file",
			config:  &Config{Files: []string{}},
			args:    args{"file"},
			wantErr: false,
		},
		{
			name:    "Add already exist file",
			config:  &Config{Files: []string{"file"}},
			args:    args{"file"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.AddFile(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("Config.AddFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_RemoveFile(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		config  *Config
		args    args
		wantErr bool
	}{
		{name: "Remove file", config: &Config{Files: []string{"file"}}, wantErr: false, args: args{"file"}},
		{name: "Remove inexistent file", config: &Config{Files: []string{}}, wantErr: true, args: args{"file"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.RemoveFile(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("Config.RemoveFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
