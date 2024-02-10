package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_router(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	type args struct {
		method string
		uri    string
	}
	tests := []struct {
		name string
		want want
		args args
	}{
		{
			name: "test empty metric list",
			want: want{
				code:        http.StatusOK,
				response:    "<html><head><title>Metric list</title><head><body><h1>All metrics</h1><h2>Gauges</h2><table><tr><th>Metric</th><th>Value</th></tr></table><h2>Counters</h2><table><tr><th>Metric</th><th>Value</th></tr></table></body></html>",
				contentType: "text/html; charset=utf-8",
			},
		},
		{
			name: "test invalid method options",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "",
			},
			args: args{
				method: http.MethodOptions,
				uri:    "/",
			},
		},
		{
			name: "test invalid method patch",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "",
			},
			args: args{
				method: http.MethodPatch,
				uri:    "/",
			},
		},
		{
			name: "test invalid method put",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "",
			},
			args: args{
				method: http.MethodPut,
				uri:    "/",
			},
		},
		{
			name: "test counter metric not found",
			want: want{
				code:        http.StatusBadRequest,
				response:    "Bad request\n",
				contentType: "text/plain; charset=utf-8",
			},
			args: args{
				method: http.MethodPost,
				uri:    "/update/unknown/testCounter/100",
			},
		},
		{
			name: "test counter metric not found",
			want: want{
				code:        http.StatusNotFound,
				response:    "Invalid metrics\n",
				contentType: "text/plain; charset=utf-8",
			},
			args: args{
				method: http.MethodPost,
				uri:    "/update/counter/",
			},
		},
		{
			name: "test counter invalid request",
			want: want{
				code:        http.StatusBadRequest,
				response:    "Invalid request\n",
				contentType: "text/plain; charset=utf-8",
			},
			args: args{
				method: http.MethodPost,
				uri:    "/update/counter/100",
			},
		},
		{
			name: "test set counter ok",
			want: want{
				code:        http.StatusOK,
				response:    "",
				contentType: "",
			},
			args: args{
				method: http.MethodPost,
				uri:    "/update/counter/100/1",
			},
		},
		{
			name: "test get counter not found",
			want: want{
				code:        http.StatusNotFound,
				response:    "unknown counter metric",
				contentType: "text/plain; charset=utf-8",
			},
			args: args{
				method: http.MethodGet,
				uri:    "/value/counter/1001",
			},
		},
		{
			name: "test get counter ok",
			want: want{
				code:        http.StatusOK,
				response:    "1",
				contentType: "text/plain; charset=utf-8",
			},
			args: args{
				method: http.MethodGet,
				uri:    "/value/counter/100",
			},
		},
		{
			name: "test gauge metric not found",
			want: want{
				code:        http.StatusNotFound,
				response:    "Invalid metrics\n",
				contentType: "text/plain; charset=utf-8",
			},
			args: args{
				method: http.MethodPost,
				uri:    "/update/gauge/",
			},
		},
		{
			name: "test gauge bad request",
			want: want{
				code:        http.StatusBadRequest,
				response:    "Invalid request\n",
				contentType: "text/plain; charset=utf-8",
			},
			args: args{
				method: http.MethodPost,
				uri:    "/update/gauge/100",
			},
		},
		{
			name: "test set gauge ok",
			want: want{
				code:        http.StatusOK,
				response:    "",
				contentType: "",
			},
			args: args{
				method: http.MethodPost,
				uri:    "/update/gauge/100/1.123",
			},
		},
		{
			name: "test get gauge not found",
			want: want{
				code:        http.StatusNotFound,
				response:    "unknown gauge metric",
				contentType: "text/plain; charset=utf-8",
			},
			args: args{
				method: http.MethodGet,
				uri:    "/value/gauge/1001",
			},
		},
		{
			name: "test get gauge ok",
			want: want{
				code:        http.StatusOK,
				response:    "1.123",
				contentType: "text/plain; charset=utf-8",
			},
			args: args{
				method: http.MethodGet,
				uri:    "/value/gauge/100",
			},
		},
	}

	storage := NewMemStorage()
	testServer := httptest.NewServer(Router(storage))
	defer testServer.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request, _ := http.NewRequest(tt.args.method, testServer.URL+tt.args.uri, nil)
			t.Log(testServer.URL + tt.args.uri)

			res, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, string(resBody))
		})
	}
}
