package httpreserve

import (
	"fmt"
	"net/http"
)

func handleOptions(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Return server OPTIONS here.")
}
