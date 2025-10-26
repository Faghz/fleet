package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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
	{"National Monument (Monas)", -6.1753924, 106.8271528},
	{"Istiqlal Mosque", -6.169856, 106.830759},
	{"Borobudur Temple", -7.6079, 110.2038},
	{"Jakarta Cathedral", -6.1690, 106.8330},
	{"Taman Mini Indonesia Indah", -6.3024, 106.8952},
	{"Grand Indonesia Mall", -6.1951, 106.8227},
	{"Ancol Dreamland", -6.1173, 106.8584},
	{"Mount Bromo", -7.9425, 112.9533},
	{"Kuta Beach", -8.7203, 115.1671},
	{"Prambanan Temple", -7.7520, 110.4915},
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
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Worker loop
	for {
		// Select random POI
		poi := pois[random.Intn(len(pois))]

		// Decide if near or far (50% chance)
		isNear := random.Float32() < 0.5

		var latitude, longitude float64
		if isNear {
			// Generate random point within 50 meters radius
			distance := random.Float64() * 50
			angle := random.Float64() * 2 * math.Pi
			deltaLat := distance * math.Cos(angle) / 111320
			deltaLon := distance * math.Sin(angle) / (111320 * math.Cos(poi.Latitude*math.Pi/180))
			latitude = poi.Latitude + deltaLat
			longitude = poi.Longitude + deltaLon
		} else {
			poi1 := pois[random.Intn(len(pois))]
			poi2 := pois[random.Intn(len(pois))]
			fraction := random.Float64()
			latitude = poi1.Latitude + fraction*(poi2.Latitude-poi1.Latitude)
			longitude = poi1.Longitude + fraction*(poi2.Longitude-poi1.Longitude)
		}

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
