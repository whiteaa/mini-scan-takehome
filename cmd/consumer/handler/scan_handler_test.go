package handler

import (
	"context"
	"encoding/json"
	"mini-scan-takehome/pkg/scanning"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDB struct{}

func (m *MockDB) InsertScan(context.Context, scanning.ReceivedScan, string) error {
	return nil
}

func TestProcessScan(t *testing.T) {
	tests := []struct {
		name    string
		version int
		data    interface{}
		want    error
	}{
		{
			"V1 Data OK",
			scanning.V1,
			scanning.V1Data{ResponseBytesUtf8: []byte("v1")},
			nil,
		},
		{
			"V2 Data OK",
			scanning.V2,
			scanning.V2Data{ResponseStr: "v2"},
			nil,
		},
		{
			"Unsupported Error",
			-1,
			scanning.V1Data{},
			ErrUnsupportedData,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, _ := json.Marshal(test.data)
			scan := scanning.ReceivedScan{DataVersion: test.version, Data: data}
			assert.Equal(t, test.want, ProcessScan(context.Background(), scan, &MockDB{}))
		})
	}
}
