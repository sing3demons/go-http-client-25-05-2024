package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	url := "https://official-joke-api.appspot.com/random_joke"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Read and use the response data here
	fmt.Println("Response status:", resp.Status)
	fmt.Println("Response headers:", resp.Header)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Response body:", string(data))
}
