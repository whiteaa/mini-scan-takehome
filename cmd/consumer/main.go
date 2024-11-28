package main

import (
	"context"
	"flag"
	"mini-scan-takehome/cmd/consumer/db"
	"mini-scan-takehome/cmd/consumer/handler"
	sub "mini-scan-takehome/cmd/consumer/subscriber"
)

func main() {
	projectId := flag.String("project", "test-project", "GCP Project ID")
	subId := flag.String("subscription", "scan-sub", "GCP PubSub Subscription ID")

	ctx := context.Background()

	// Initialize DB client
	db.InitClient(ctx)

	// Initialize Subscriber client
	sub.InitClient(ctx, *projectId)
	subscriber := sub.GetClient()

	// Subscribe to a GCP topic and process messages
	subscriber.Subscribe(ctx, *subId, handler.ReceiveScan)
}
