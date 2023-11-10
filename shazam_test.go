package shazam_test

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/aditya109/shazam"
)

func TestBoom(t *testing.T) {
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "testing GET request with no headers",
			args: args{
				request: &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme: "https",
						Host:   "reqres.in",
						Path:   "/api/users?page=2",
					},
					Body: nil,
				},
			},
			want:    "curl --location --request GET 'https://reqres.in/api/users?page=2'",
			wantErr: false,
		},
		{
			name: "testing GET request with one header",
			args: args{
				request: &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme: "https",
						Host:   "reqres.in",
						Path:   "/api/users",
					},
					Header: map[string][]string{
						"x-panem-token": {
							"BUM99779r42aUZsadasdasdUZB8Z95YLK",
						},
					},
					Body: nil,
				},
			},
			want:    "curl --location --request GET 'https://reqres.in/api/users' --header 'x-panem-token: BUM99779r42aUZsadasdasdUZB8Z95YLK'",
			wantErr: false,
		},
		{
			name: "testing GET request with one header and multiple query parameters",
			args: args{
				request: &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme:   "https",
						Host:     "reqres.in",
						Path:     "/api/users",
						RawQuery: "err=reoted&page=2",
					},
					Header: map[string][]string{
						"x-panem-token": {
							"BUM99779r42aUZsadasdasdUZB8Z95YLK",
						},
					},
					Body: nil,
				},
			},
			want:    "curl --location --request GET 'https://reqres.in/api/users?err=reoted&page=2' --header 'x-panem-token: BUM99779r42aUZsadasdasdUZB8Z95YLK'",
			wantErr: false,
		},
		{
			name: "testing POST request with no body",
			args: args{
				request: &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme: "https",
						Host:   "reqres.in",
						Path:   "/api/users",
					},
					Body: nil,
				},
			},
			want:    "curl --location --request POST 'https://reqres.in/api/users'",
			wantErr: false,
		},
		{
			name: "testing POST request with valid body",
			args: args{
				request: getSampleRequestForCreateUserPayload(),
			},
			want:    "curl --location --request POST 'https://reqres.in/api/users' --data-raw '{\"name\": \"morpheus\", \"job\": \"leader\"}'",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shazam.Boom(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Boom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Boom() = %v, want %v", got, tt.want)
			}
			if !checkIfFailSafeWorks(tt.args.request) {
				t.Errorf("failSafe() did not work as expected. %v", *tt.args.request)
			}
		})
	}
}

func getSampleRequestForCreateUserPayload() *http.Request {
	request, _ := http.NewRequest(
		http.MethodPost,
		"https://reqres.in/api/users",
		bytes.NewReader([]byte(`{"name": "morpheus", "job": "leader"}`)),
	)
	return request
}

func checkIfFailSafeWorks(request *http.Request) bool {
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return false
	}
	if resp.StatusCode < 400 {
		return true
	}
	return false
}
