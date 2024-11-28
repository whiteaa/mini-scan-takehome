package db

import (
	"context"
	"sync"

	"mini-scan-takehome/pkg/scanning"

	"github.com/jackc/pgx/v5"
)

type ScansDB struct {
	client *pgx.Conn
}

type ScansDBInterface interface {
	InsertScan(context.Context, scanning.ReceivedScan, string) error
}

var once sync.Once
var database ScansDBInterface

func GetClient() ScansDBInterface {
	return database
}

func InitClient(ctx context.Context) {
	once.Do(func() {
		connString := "user=postgres password=password host=postgres port=5432 dbname=postgres"
		conn, err := pgx.Connect(ctx, connString)
		if err != nil {
			panic(err)
		}
		database = &ScansDB{client: conn}
	})
}

// InsertScan inserts a scan into the db if ip, port, service are unique.
// Else updates scanned_at, data_version, and data for the existing record.
func (s *ScansDB) InsertScan(ctx context.Context, scan scanning.ReceivedScan, data string) error {
	query := `INSERT INTO scan (ip, port, service, scanned_at, data_version, data) 
VALUES (@ip, @port, @service, @scanned_at, @data_version, @data) 
ON CONFLICT (ip, port, service) 
DO UPDATE SET 
scanned_at=@scanned_at, 
data_version=@data_version, 
data=@data`
	args := pgx.NamedArgs{
		"ip":           scan.Ip,
		"port":         scan.Port,
		"service":      scan.Service,
		"scanned_at":   scan.Timestamp,
		"data_version": scan.DataVersion,
		"data":         data,
	}
	_, err := s.client.Exec(ctx, query, args)
	return err
}
