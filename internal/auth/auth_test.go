package auth

import (
	"errors"
	"net/http/httptest"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		headers     map[string]string
		want        string
		wantErr     error
		expectError bool
		errMsg      string // Добавляем поле для ожидаемого сообщения об ошибке
	}{
		{
			name:        "1. no authorization header",
			headers:     nil,
			want:        "",
			wantErr:     ErrNoAuthHeaderIncluded,
			expectError: true,
			errMsg:      "no authorization header included",
		},
		{
			name: "2. malformed authorization header - no ApiKey prefix",
			headers: map[string]string{
				"Authorization": "Bearer token",
			},
			want:        "",
			wantErr:     nil,
			expectError: true,
			errMsg:      "malformed authorization header",
		},
		{
			name: "3. malformed authorization header - empty",
			headers: map[string]string{
				"Authorization": "",
			},
			want:        "",
			wantErr:     ErrNoAuthHeaderIncluded, // Изменяем ожидаемую ошибку
			expectError: true,
			errMsg:      "no authorization header included",
		},
		{
			name: "4. valid api key",
			headers: map[string]string{
				"Authorization": "ApiKey valid-key-123",
			},
			want:        "valid-key-123",
			wantErr:     nil,
			expectError: false,
			errMsg:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.headers != nil {
				for k, v := range tt.headers {
					req.Header.Set(k, v)
				}
			}

			got, err := GetAPIKey(req.Header)

			if (err != nil) != tt.expectError {
				t.Errorf("GetAPIKey() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if tt.expectError {
				if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
					t.Errorf("GetAPIKey() wantErr = %v, got %v", tt.wantErr, err)
				}
				if err != nil && err.Error() != tt.errMsg {
					t.Errorf("Unexpected error message: got %q, want %q", err.Error(), tt.errMsg)
				}
			}

			if got != tt.want {
				t.Errorf("GetAPIKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
