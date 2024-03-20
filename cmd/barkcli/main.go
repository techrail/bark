package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func main() {
	// Parse command-line arguments
	logLevel := flag.String("logLevel", "", "Log level to filter by")
	serviceName := flag.String("serviceName", "", "Service name to filter by")
	sessionName := flag.String("sessionName", "", "Session name to filter by")
	startDate := flag.String("startDate", "", "Start date for logs")
	endDate := flag.String("endDate", "", "End date for logs")
	flag.Parse()

	// Make a GET request to the /fetchLogs endpoint with the provided parameters
	logs := makeGetRequestToFetchLogs(*logLevel, *serviceName, *sessionName, *startDate, *endDate)

	// Display the fetched logs
	fmt.Println(logs)
}

func makeGetRequestToFetchLogs(logLevel, serviceName, sessionName, startDate, endDate string) string {
	serverURL := "http://localhost:8080/fetchLogs" // Change this to your server's URL
	params := url.Values{}
	params.Add("logLevel", logLevel)
	params.Add("serviceName", serviceName)
	params.Add("sessionName", sessionName)
	params.Add("startDate", startDate)
	params.Add("endDate", endDate)

	resp, err := http.Get(fmt.Sprintf("%s?%s", serverURL, params.Encode()))
	if err != nil {
		return fmt.Sprintf("Error making request: %s", err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error reading response: %s", err.Error())
	}
	return string(body)
}
