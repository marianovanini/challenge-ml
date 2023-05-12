package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

const (
	saveDirectory = "./data"
)

type SystemInfo struct {
	Processor        *string `json:"processor"`
	RunningProcesses *string `json:"running_processes"`
	Users            *string `json:"users"`
	OSName           *string `json:"os_name"`
	OSVersion        *string `json:"os_version"`
}

func main() {

	err := godotenv.Load("app.env")
	if err != nil {
		fmt.Println("Error loading app.env")
	}

	// Save directory
	err = os.MkdirAll(saveDirectory, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Endpoints
	http.HandleFunc("/system-info", handleSystemInfo)

	log.Println("API server started in port " + os.Getenv("API_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("API_PORT"), nil))
}

func handleSystemInfo(w http.ResponseWriter, r *http.Request) {
	// Filename
	timestamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s_%s", r.RemoteAddr, timestamp)

	// Parse the data
	reader := csv.NewReader(r.Body)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Println("Failed to read system info:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Convert data to SystemInfo struct
	sysInfo := parseSystemInfo(rows)

	// Create the CSV file
	csvFilePath := filepath.Join(saveDirectory, filename+".csv")
	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		log.Println("Failed to create CSV file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer csvFile.Close()

	// Write CSV data to the file
	csvWriter := csv.NewWriter(csvFile)
	err = csvWriter.WriteAll(rows)
	if err != nil {
		log.Println("Failed to write CSV data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	csvWriter.Flush()

	// Create the JSON file
	jsonFilePath := filepath.Join(saveDirectory, filename+".json")
	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		log.Println("Failed to create JSON file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer jsonFile.Close()

	// Convert SystemInfo struct to JSON
	jsonData, err := json.MarshalIndent(sysInfo, "", "\t")
	if err != nil {
		log.Println("Failed to convert system info to JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write JSON data to the file
	_, err = jsonFile.Write(jsonData)
	if err != nil {
		log.Println("Failed to save system info as JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("System info saved:", filename)
}

func parseSystemInfo(rows [][]string) *SystemInfo {
	if len(rows) != 5 {
		return nil
	}

	sysInfo := &SystemInfo{
		Processor:        &rows[0][1],
		RunningProcesses: &rows[1][1],
		Users:            &rows[2][1],
		OSName:           &rows[3][1],
		OSVersion:        &rows[4][1],
	}

	return sysInfo
}
