package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type POI struct {
	Name      string
	Latitude  float64
	Longitude float64
}

type LocationMessage struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
	IsNear    bool    `json:"is_near"`
}

var pois = []POI{
	{"Jakarta1", -6.2087634, 106.845599},
	{"Yogyakarta", -7.7956, 110.3695},
	{"Bandung", -6.9175, 107.6191},
	{"Jakarta2", -6.1751, 106.8227},
}

var vehicleIDs = []string{"B1234ABC", "B5678DEF", "D9012GHI", "B3456JKL", "E7890MNO"}

func main() {
	// MQTT broker address (adjust as needed)
	broker := "tcp://localhost:1883"
	if envBroker := os.Getenv("MQTT_BROKER"); envBroker != "" {
		broker = envBroker
	}
	clientID := "location-sync-worker"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetAutoReconnect(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	defer client.Disconnect(250)

	log.Println("Connected to MQTT broker")

	// Seed random number generator
	random := rand.New(rand.NewSource(time.Hour.Nanoseconds()))

	// Worker loop
	for {
		// Select random POI
		poi := pois[random.Intn(len(pois))]

		// Decide if near or far (50% chance)
		isNear := rand.Float32() < 0.5

		var latOffset, lonOffset float64
		if isNear {
			// Near: within ~1km (0.01 degrees)
			latOffset = (random.Float64() - 0.5) * 0.02
			lonOffset = (random.Float64() - 0.5) * 0.02
		} else {
			// Far: within ~50km (0.5 degrees)
			latOffset = (random.Float64() - 0.5) * 1.0
			lonOffset = (random.Float64() - 0.5) * 1.0
		}

		latitude := poi.Latitude + latOffset
		longitude := poi.Longitude + lonOffset

		// Select random vehicle ID
		vehicleID := vehicleIDs[random.Intn(len(vehicleIDs))]

		message := LocationMessage{
			VehicleID: vehicleID,
			Latitude:  latitude,
			Longitude: longitude,
			Timestamp: time.Now().Unix(),
			IsNear:    isNear,
		}

		payload, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			continue
		}

		topic := fmt.Sprintf("fleet/vehicle/%s/location", vehicleID)
		token := client.Publish(topic, 0, false, payload)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Failed to publish message: %v", token.Error())
		} else {
			log.Printf("Published location for %s: lat=%.6f, lon=%.6f, near=%v", vehicleID, latitude, longitude, isNear)
		}

		// Wait 5 seconds before next update
		time.Sleep(2 * time.Second)
	}
}
