package configmapper

import (
	"os"
	"reflect"
	"testing"
)

type Config struct {
	EnvConfig EnvConfig `mapstructure:"env"`
}
type EnvConfig struct {
	Env         string `validate:"required"`
	ServiceName string `mapstructure:"service_name" validate:"required"`
	ProjectID   string `mapstructure:"project_id" validate:"required"`
	LogLevel    string `mapstructure:"log_level" validate:"required"`
}

func TestInitialize(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	os.Setenv(FileEnvKey, "config_test.toml")
	type args struct {
		config interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "",
			args: args{
				config: Config{},
			},
			want: Config{
				EnvConfig: EnvConfig{
					Env:         "1",
					ServiceName: "2",
					ProjectID:   "3",
					LogLevel:    "4",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Initialize(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Initialize() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializeManifest(t *testing.T) {
	os.Setenv("LOCAL", "1")
	type args struct {
		config interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "",
			args: args{
				config: Config{},
			},
			want: Config{
				EnvConfig: EnvConfig{
					Env:         "1",
					ServiceName: "2",
					ProjectID:   "3",
					LogLevel:    "4",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Initialize(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Initialize() got = %v, want %v", got, tt.want)
			}
		})
	}
}
