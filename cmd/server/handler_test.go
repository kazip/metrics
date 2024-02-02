package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handleCounter(t *testing.T) {

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
			name: "test 1",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "Method not Allowed\n",
				contentType: "text/plain; charset=utf-8",
			},
			args: args{
				method: http.MethodGet,
				uri:    "/",
			},
		},
		{
			name: "test metric not found",
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
			name: "test invalid request",
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
			name: "test ok",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(tt.args.method, tt.args.uri, nil)
			storage := NewMemStorage()
			w := httptest.NewRecorder()

			handler := handleCounterFunc(storage)
			handler(w, request)
			res := w.Result()

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, string(resBody))
		})
	}

}

func Test_handleGauge(t *testing.T) {

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
			name: "test invalid method get",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "Method not Allowed\n",
				contentType: "text/plain; charset=utf-8",
			},
			args: args{
				method: http.MethodGet,
				uri:    "/",
			},
		},
		{
			name: "test invalid method options",
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "Method not Allowed\n",
				contentType: "text/plain; charset=utf-8",
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
				response:    "Method not Allowed\n",
				contentType: "text/plain; charset=utf-8",
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
				response:    "Method not Allowed\n",
				contentType: "text/plain; charset=utf-8",
			},
			args: args{
				method: http.MethodPut,
				uri:    "/",
			},
		},
		{
			name: "test metric not found",
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
			name: "test bad request",
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
			name: "test ok",
			want: want{
				code:        http.StatusOK,
				response:    "",
				contentType: "",
			},
			args: args{
				method: http.MethodPost,
				uri:    "/update/gauge/100/1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemStorage()
			request := httptest.NewRequest(tt.args.method, tt.args.uri, nil)

			w := httptest.NewRecorder()
			handler := handleGaugeFunc(storage)
			handler(w, request)
			res := w.Result()

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, string(resBody))
		})
	}
}
