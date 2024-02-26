package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	url := "https://official-joke-api.appspot.com/random_joke"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// req.Header.Add("Authorization", "Bearer API_KEY")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer res.Body.Close()

	fmt.Println("Status code:", res.StatusCode)
	fmt.Println("Status:", res.Status)
	fmt.Println("Header:", res.Header)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	fmt.Println("Body:", string(body))
}
