package slack

import "testing"

func TestErrorResponse_Error(t *testing.T) {
	tests := []struct {
		name   string
		fields ErrorResponse
		want   string
	}{
		{
			name:   "basic test",
			fields: ErrorResponse{ErrorMessage: "something sent wrong"},
			want:   "Slack API Error: [something sent wrong]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.fields.Error(); got != tt.want {
				t.Errorf("ErrorResponse.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
