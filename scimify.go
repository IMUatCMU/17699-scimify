package main

import (
	"github.com/go-scim/scimify/service"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", service.Mux())
}
