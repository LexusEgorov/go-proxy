package client

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/LexusEgorov/go-proxy/internal/config"
)

func TestClient_Request(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, header := range r.Header {

			if len(header) == 1 {
				w.Header().Set(key, header[0])
				continue
			}

			for _, value := range header {
				w.Header().Add(key, value)
			}
		}

		body, err := io.ReadAll(r.Body)

		if err != nil {
			t.Fatalf("mockServer.Read: %v", err)
		}

		_, err = w.Write(body)

		if err != nil {
			t.Fatalf("mockServer.Write: %v", err)
		}
	}))
	defer testServer.Close()

	testClient := New(config.ClientConfig{
		URL: testServer.URL,
	})

	type args struct {
		method  string
		url     string
		body    io.Reader
		headers http.Header
	}
	type response struct {
		code    int
		body    []byte
		headers http.Header
	}
	tests := []struct {
		name    string
		c       Client
		args    args
		want    response
		wantErr bool
	}{
		{
			name: "simple request",
			c:    *testClient,
			args: args{
				url:     "/users",
				method:  http.MethodGet,
				body:    nil,
				headers: nil,
			},
			want: response{
				code:    http.StatusOK,
				body:    []byte{},
				headers: http.Header{},
			},
			wantErr: false,
		},
		{
			name: "body request",
			c:    *testClient,
			args: args{
				url:     "/users",
				method:  http.MethodPost,
				body:    bytes.NewBuffer([]byte("test")),
				headers: nil,
			},
			want: response{
				code:    http.StatusOK,
				body:    []byte("test"),
				headers: http.Header{},
			},
			wantErr: false,
		},
		{
			name: "header request",
			c:    *testClient,
			args: args{
				url:    "/users",
				method: http.MethodGet,
				body:   nil,
				headers: http.Header{
					"Test-Header": []string{"testVal"},
				},
			},
			want: response{
				code: http.StatusOK,
				body: []byte{},
				headers: http.Header{
					"Test-Header": []string{"testVal"},
				},
			},
			wantErr: false,
		},
		{
			name: "mult header request",
			c:    *testClient,
			args: args{
				url:    "/users",
				method: http.MethodGet,
				body:   nil,
				headers: http.Header{
					"Test-Header": []string{"1", "2"},
				},
			},
			want: response{
				code: http.StatusOK,
				body: []byte{},
				headers: http.Header{
					"Test-Header": []string{"1", "2"},
				},
			},
			wantErr: false,
		},
		{
			name: "empty url",
			c:    *New(config.ClientConfig{}),
			args: args{
				url:     "/users",
				method:  http.MethodGet,
				body:    nil,
				headers: nil,
			},
			want: response{
				code:    http.StatusOK,
				body:    []byte{},
				headers: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Request(tt.args.method, tt.args.url, tt.args.body, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got.StatusCode() != tt.want.code {
				t.Errorf("Client.Code = %v, want %v", got.StatusCode(), tt.want.code)
			}

			body := got.Body()
			if !reflect.DeepEqual(tt.want.body, body) {
				t.Errorf("Client.Body = %v, want %v", body, tt.want.body)
			}

			for header, value := range tt.want.headers {
				gotHeader := got.Header()[header]
				if !reflect.DeepEqual(gotHeader, value) {
					t.Errorf("Client.Header('%s') = %v, want %v", header, gotHeader, value)
				}
			}
		})
	}
}
