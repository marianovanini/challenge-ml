package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("app.env")
	if err != nil {
		fmt.Println("Error loading app.env")
	}

	http.HandleFunc("/send-info", handleSendInfo)
	log.Println("Agent program is listening in port " + os.Getenv("AGENT_PORT") + "..")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("AGENT_PORT"), nil))

}

func handleSendInfo(w http.ResponseWriter, r *http.Request) {

	processorInfo := getProcessorInfo()
	processes := getRunningProcesses()
	users := getUsersWithSession()
	osName := runtime.GOOS
	osVersion := getOSVersion()

	csvData := [][]string{
		{"Processor", processorInfo},
		{"Running Processes", strings.Join(processes, ",")},
		{"Users", strings.Join(users, ",")},
		{"OS Name", osName},
		{"OS Version", osVersion},
	}

	// Create CSV buffer
	csvBuffer := &bytes.Buffer{}
	writer := csv.NewWriter(csvBuffer)

	// Write data to the buffer
	for _, data := range csvData {
		err := writer.Write(data)
		if err != nil {
			log.Fatal(err)
		}
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}

	sendDataToAPI(csvBuffer.Bytes())
}

func getProcessorInfo() string {
	cmd := exec.Command("uname", "-p")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(output))
}

func getRunningProcesses() []string {
	cmd := exec.Command("ps", "-e", "-o", "comm=")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n")
}

func getUsersWithSession() []string {
	cmd := exec.Command("who", "-q")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) > 1 {
		return strings.Split(lines[1], " ")
	}
	return []string{}
}

func getOSVersion() string {
	cmd := exec.Command("uname", "-v")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(output))
}

func sendDataToAPI(data []byte) {
	req, err := http.NewRequest(http.MethodPost, os.Getenv("API_URL")+"/system-info", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/csv")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API request failed with status: %s", resp.Status)
	}

	fmt.Println("System information sent to API successfully.")
}
