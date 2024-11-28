package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"mini-scan-takehome/cmd/consumer/db"
	"mini-scan-takehome/pkg/scanning"

	"cloud.google.com/go/pubsub"
)

var ErrUnsupportedData = errors.New("unsupported data type")

// ProcessScan processes the scan based on the data version and inserts the scan.
func ProcessScan(ctx context.Context, scan scanning.ReceivedScan, scansDB db.ScansDBInterface) error {
	var data string
	switch scan.DataVersion {
	case scanning.V1:
		var v1Data scanning.V1Data
		if err := json.Unmarshal(scan.Data, &v1Data); err != nil {
			return err
		}
		data = string(v1Data.ResponseBytesUtf8)
	case scanning.V2:
		var v2Data scanning.V2Data
		if err := json.Unmarshal(scan.Data, &v2Data); err != nil {
			return err
		}
		data = v2Data.ResponseStr
	default:
		return ErrUnsupportedData
	}

	return scansDB.InsertScan(ctx, scan, data)
}

func ReceiveScan(ctx context.Context, m *pubsub.Message) error {
	var scan scanning.ReceivedScan
	if err := json.Unmarshal(m.Data, &scan); err != nil {
		log.Printf("Failed to receive message. %v\n", err)
	}
	return ProcessScan(ctx, scan, db.GetClient())
}
