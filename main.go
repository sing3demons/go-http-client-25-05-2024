package main

import (
	"fmt"

	"github.com/sing3demons/go-http-client/client"
)

func main() {
	url := "https://official-joke-api.appspot.com/random_joke"
	// client := &http.Client{
	// 	Timeout: 10 * time.Second,
	// }

	c := client.NewHttpClient()
	body, err := c.Get(url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(body)
}
