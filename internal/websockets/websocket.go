package websockets

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	pb "sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartWebSocket() error {
	url := "https://api.opticodds.com/api/v2/stream/odds?sportsbooks=bwin&league=UEFA%20-%20Champions%20League&key=88f9bd7f-463c-44ca-b938-fd5bf2704e52"

	// Create an HTTP client
	client := &http.Client{}

	// Send a GET request
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Use bufio.Scanner to read the stream line by line or in chunks
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text() // Get the line of data as string
		if !strings.HasPrefix(line, "data: ") {
			continue // Skip lines that do not start with "data: "
		}

		// Trim the prefix "data: " to get the JSON string
		jsonStr := strings.TrimPrefix(line, "data: ")

		var oddsData types.OddsStream
		// Unmarshal JSON from the JSON string
		if err := json.Unmarshal([]byte(jsonStr), &oddsData); err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		// Process the data
		log.Printf("Received odds data: %+v", oddsData)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from stream: %v", err)
	}
	return nil
}

func unmarshal(data []byte, liveOddsData *pb.LiveOddsData) {
	// Implement the conversion of binary data from WebSocket to LiveOddsData struct
	// You could use protobuf's Unmarshal here if the WebSocket sends protobuf binary
	// Or implement a custom parser if it's a JSON or some other format
}
