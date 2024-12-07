package main

import (
	"database/sql"
    _ "github.com/lib/pq"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"time"
	"log"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option" 
	"github.com/censys/scan-takehome/pkg/scanning"
)

var db *sql.DB

func initDB() {
    var err error
    // Connection string example: "postgres://username:password@localhost/dbname?sslmode=disable"
    connStr := "postgres://username:password@localhost/dbname?sslmode=disable"
    
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Error connecting to the database: ", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Error pinging database: ", err)
    }

    log.Println("Connected to PostgreSQL")
}


func main() {
	initDB()
    emulatorHost := os.Getenv("PUBSUB_EMULATOR_HOST")
    log.Println("Connecting to Pub/Sub Emulator at:", emulatorHost)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := pubsub.NewClient(ctx, "test-project", option.WithEndpoint(emulatorHost+":8085"), option.WithoutAuthentication())
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer client.Close()

    sub := client.Subscription("scan-sub")
    sub.ReceiveSettings.Synchronous = true
    sub.ReceiveSettings.MaxOutstandingMessages = 10

    err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
        log.Println("Received message:", string(msg.Data))
        handleMessage(msg)
        msg.Ack()
    })

    if err != nil {
        log.Fatalf("Receive failed: %v", err)
    }
}



func handleMessage(msg *pubsub.Message) {
    var scan scanning.Scan
    if err := json.Unmarshal(msg.Data, &scan); err != nil {
        log.Printf("Error decoding message: %v", err)
        return
    }

    var response string
    switch scan.DataVersion {
    case scanning.V1:
        var data scanning.V1Data
        if err := json.Unmarshal(scan.Data, &data); err != nil {
            log.Printf("Error unmarshalling V1 data: %v", err)
            return
        }
        decoded, err := base64.StdEncoding.DecodeString(data.ResponseBytesUtf8)
        if err != nil {
            log.Printf("Error decoding base64: %v", err)
            return
        }
        response = string(decoded)
        log.Printf("Processed V1 data for IP %s: %s", scan.Ip, response)

    case scanning.V2:
        var data scanning.V2Data
        if err := json.Unmarshal(scan.Data, &data); err != nil {
            log.Printf("Error unmarshalling V2 data: %v", err)
            return
        }
        response = data.ResponseStr
        log.Printf("Processed V2 data for IP %s: %s", scan.Ip, response)

    default:
        log.Printf("Unknown data version %d", scan.DataVersion)
        return
    }

   
    err := storeData(scan.Ip, uint32(scan.Port), scan.Service, response)
    if err != nil {
        log.Printf("Error storing data for IP %s: %v", scan.Ip, err)
    }
}



func storeData(ip string, port uint32, service, response string) error {
    query := `INSERT INTO scan_records (ip, port, service, last_scanned, response)
          VALUES ($1, $2, $3, NOW(), $4)
          ON CONFLICT (ip, port, service) DO UPDATE
          SET last_scanned = excluded.last_scanned, response = excluded.response;`
	_, err := db.Exec(query, ip, port, service, response)

    if err != nil {
        log.Printf("Error storing data: %v", err)
        return err
    }
    log.Printf("Data stored/updated for %s:%d/%s", ip, port, service)
    return nil
}
