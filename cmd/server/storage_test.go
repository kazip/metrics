package main

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_SetGauge(t *testing.T) {
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test ok",
			args: args{
				name:  "test",
				value: 0.0,
			},
			wantErr: false,
		},
		{
			name: "test empty name",
			args: args{
				name:  "",
				value: 1.0,
			},
			wantErr: true,
		},
		{
			name: "test -Inf",
			args: args{
				name:  "test2",
				value: math.Inf(-1),
			},
			wantErr: false,
		},
		{
			name: "test +Inf",
			args: args{
				name:  "test3",
				value: math.Inf(1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemStorage()
			err := m.SetGauge(tt.args.name, tt.args.value)
			if tt.wantErr {
				require.NotNil(t, err)
				return
			}

			if assert.Contains(t, m.Gauges, tt.args.name) {
				assert.Equal(t, m.Gauges[tt.args.name], tt.args.value)
			}

		})
	}
}

func TestMemStorage_SetCounter(t *testing.T) {
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test ok",
			args: args{
				name:  "test",
				value: 0,
			},
			wantErr: false,
		},
		{
			name: "test empty name",
			args: args{
				name:  "",
				value: 1,
			},
			wantErr: true,
		},
		{
			name: "test maxInt64",
			args: args{
				name:  "test2",
				value: math.MaxInt64,
			},
			wantErr: false,
		},		
		{
			name: "test minInt64",
			args: args{
				name:  "test3",
				value: math.MinInt64,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemStorage()

			if err := m.SetCounter(tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.SetCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
