package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", NewRouter())
}
