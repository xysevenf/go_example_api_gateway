package main

import (
	"api-gateway/config"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("config load error")
		return
	}

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", NewRouter(cfg))
}
