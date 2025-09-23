package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JuanPabloJimenez0250013/orbital-net/pkg/discovery/consul"
	discovery "github.com/JuanPabloJimenez0250013/orbital-net/pkg/registry"
	"github.com/JuanPabloJimenez0250013/orbital-net/sim-service/handler"
	"github.com/JuanPabloJimenez0250013/orbital-net/sim-service/model"
)

const (
	SERVICE_NAME = "sim-service"
	PORT         = 8081
	CONSUL_HOST  = "consul"
)

func main() {
	model.CreateSimObjects()

	// Start HTTP server
	http.HandleFunc("/", handler.NodesHandler)
	go func() {
		addr := fmt.Sprintf(":%d", PORT)
		fmt.Printf("🌐 Server listening on %s%s\n", SERVICE_NAME, addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Register service with Consul
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	registry, instanceID := registerWithConsul(ctx)

	// Capture shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Simulation loop in a goroutine
	go func() {
		for {
			for _, node := range model.Nodes {
				node.Move()
			}
			time.Sleep(time.Millisecond * model.TICK_DELAY)
		}
	}()

	// Block until signal received
	<-sigCh
	fmt.Println("🛑 Shutting down...")

	// Deregister service before exit
	if err := registry.Deregister(ctx, instanceID, SERVICE_NAME); err != nil {
		log.Printf("⚠️ Failed to deregister service: %v", err)
	} else {
		fmt.Println("✅ Service deregistered from Consul")
	}
}

func registerWithConsul(ctx context.Context) (*consul.Registry, string) {
	registry, err := consul.NewRegistry(CONSUL_HOST + ":8500")
	if err != nil {
		log.Fatalf("❌ Failed to connect to Consul: %v", err)
	}

	instanceID := discovery.GenerateInstanceID(SERVICE_NAME)
	serviceAddr := fmt.Sprintf("%s:%d", SERVICE_NAME, PORT) // use service name instead of hostname

	if err := registry.Register(ctx, instanceID, SERVICE_NAME, serviceAddr); err != nil {
		log.Fatalf("❌ Failed to register in Consul: %v", err)
	}
	fmt.Println("✅ Service registered in Consul")

	// Keep reporting healthy state
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := registry.ReportHealthyState(instanceID, SERVICE_NAME); err != nil {
					log.Println("⚠️ Failed to report healthy state:", err)
				}
				time.Sleep(2 * time.Second)
			}
		}
	}()

	return registry, instanceID
}
