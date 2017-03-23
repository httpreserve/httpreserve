package main

import (
	"fmt"
	"net/http"
)

func handleHttpreserve(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Return server httpreserve analysis here.")
}
