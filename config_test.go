package chaakoo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	validWindow := func(name string) *Window {
		return &Window{Name: name, Grid: "vim"}
	}

	tests := []struct {
		name    string
		config  *Config
		wantErr string
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: "config is nil",
		},
		{
			name:    "no sessions",
			config:  &Config{},
			wantErr: "at-least 1 session is required",
		},
		{
			name: "session without name",
			config: &Config{
				Sessions: []*Session{
					{Windows: []*Window{validWindow("w1")}},
				},
			},
			wantErr: "session name required at index 0",
		},
		{
			name: "session without windows",
			config: &Config{
				Sessions: []*Session{
					{Name: "s1"},
				},
			},
			wantErr: "atleast 1 window is required for session - s1",
		},
		{
			name: "second session without name",
			config: &Config{
				Sessions: []*Session{
					{Name: "s1", Windows: []*Window{validWindow("w1")}},
					{Windows: []*Window{validWindow("w2")}},
				},
			},
			wantErr: "session name required at index 1",
		},
		{
			name: "valid single session",
			config: &Config{
				Sessions: []*Session{
					{Name: "s1", Windows: []*Window{validWindow("w1")}},
				},
			},
		},
		{
			name: "valid multi session",
			config: &Config{
				Sessions: []*Session{
					{Name: "s1", Windows: []*Window{validWindow("w1")}},
					{Name: "s2", Windows: []*Window{validWindow("w2")}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
