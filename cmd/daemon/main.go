package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	influxDBToken = "XE8XaPvQFxG-7EpoM0ECqX006NGsasOMDE9nGAjh7s0G3T19p-fDt7GXP6Npwn5CjAGFwYt3yVHBbx96QOtQ2A=="
	influxDBURL   = "http://localhost:8086"
	bucket        = "your_bucket"
	org           = "your_org"
)

func main() {
	// Set up InfluxDB client
	client := influxdb2.NewClient(influxDBURL, influxDBToken)
	defer client.Close()

	// TCP listener for incoming data
	PORT := ":8080"
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Error creating TCP listener:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Listening on port", PORT)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, client)
	}
}

func handleConnection(conn net.Conn, client influxdb2.Client) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		data := scanner.Text()
		// Process and write data to InfluxDB
		writeDataToInfluxDB(client, data)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from connection:", err)
	}
}

func writeDataToInfluxDB(client influxdb2.Client, data string) {
	// Split the data into its components
	parts := strings.Split(data, ", ")
	if len(parts) != 3 {
		fmt.Println("Invalid data format received:", data)
		return
	}

	// Extract data
	directory, command, durationStr := parts[0], parts[1], parts[2]

	// Convert duration to integer (milliseconds)
	var duration int
	_, err := fmt.Sscanf(durationStr, "%d", &duration)
	if err != nil {
		fmt.Println("Error parsing duration:", err)
		return
	}

	// Create a new point with the command data
	p := influxdb2.NewPointWithMeasurement("command_usage").
		AddTag("directory", directory).
		AddTag("command", command).
		AddField("duration_ms", duration).
		SetTime(time.Now())

	// Write the point
	writeAPI := client.WriteAPI(org, bucket)
	writeAPI.WritePoint(p)
	writeAPI.Flush()
}
