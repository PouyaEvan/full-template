package domain

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr error
	}{
		{
			name:    "Valid Phone",
			phone:   "09123456789",
			wantErr: nil,
		},
		{
			name:    "Invalid Phone - Short",
			phone:   "0912",
			wantErr: ErrInvalidPhone,
		},
		{
			name:    "Invalid Phone - Letters",
			phone:   "0912abc3456",
			wantErr: ErrInvalidPhone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUser(tt.phone)
			if err != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
