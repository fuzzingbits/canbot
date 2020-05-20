package slack

import (
	"net/http"
	"testing"
)

type mockableTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m mockableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

var transport = &mockableTransport{}
var testHTTPClient = &http.Client{
	Transport: transport,
}
var testService = restService{}
var testClient = &Client{}

func Test_slackRestService_ErrorCheck(t *testing.T) {
	type args struct {
		responseBytes []byte
	}
	tests := []struct {
		name    string
		fields  restService
		args    args
		wantErr bool
	}{
		{
			name:   "malformed json",
			fields: restService{},
			args: args{
				responseBytes: []byte(`{"ok" false}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := restService{
				token: tt.fields.token,
			}
			if err := s.ErrorCheck(tt.args.responseBytes); (err != nil) != tt.wantErr {
				t.Errorf("slackRestService.ErrorCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
