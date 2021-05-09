package hello

import (
	"fmt"
	"net/http"
	"time"
)

type CoverageMetric struct {
	total   uint64
	covered uint64
}

type CoverageData struct {
	org      string
	repo     string
	branch   string
	language string
	runDate  *time.Time
	metrics  map[string]CoverageMetric
}

type QueryParams struct {
	org      string
	repo     string
	branch   string
	language string
	from     *time.Time
	to       *time.Time
}

func Hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "someone"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}
