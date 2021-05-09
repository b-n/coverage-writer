package coverage

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type CoverageMetric struct {
	Total   float64 `json:"total",firestore:"total"`
	Covered float64 `json:"covered",firestore:"covered"`
}

type CoverageData struct {
	Repo     string                    `json:"repo",firestore:"repo"`
	Branch   string                    `json:"branch",firestore:"branch"`
	Language string                    `json:"language",firestore:"language"`
	RunAt    *time.Time                `json:"runAt",firestore:"runAt"`
	Metrics  map[string]CoverageMetric `json:"metrics",firestore:"metrics"`
}

func (c *CoverageData) Validate() error {
	if c.Repo == "" {
		return errors.New("repo is required")
	}
	if c.Branch == "" {
		return errors.New("branch is required")
	}
	if c.Language == "" {
		return errors.New("language is required")
	}
	if c.RunAt == nil {
		return errors.New("runAt is required")
	}

	return nil
}

var (
	ctx context.Context
	db  *firestore.Client
)

func init() {
	ctx = context.Background()
	var err error
	db, err = firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		log.Fatal(err)
	}
}

func getCoverage(w http.ResponseWriter, r *http.Request) {
	//TODO we need to be clever and query, not return all
	docs := db.Collection("coverages").Doc("b-n").Collection("coverage")

	iter := docs.DocumentRefs(ctx)

	var refs []*firestore.DocumentRef

	for {
		ref, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			handleError(w, http.StatusInternalServerError, errors.New("Failed to retrieve documents"))
			return
		}
		refs = append(refs, ref)
	}

	data, err := db.GetAll(ctx, refs)
	if err != nil {
		handleError(w, http.StatusInternalServerError, errors.New("Failed to retrieve documents"))
		return
	}

	results := make([]CoverageData, len(refs))
	for i, v := range data {
		err := v.DataTo(&results[i])
		if err != nil {
			handleError(w, http.StatusInternalServerError, errors.New("Failed to retrieve documents"))
			return
		}
	}

	json.NewEncoder(w).Encode(results)
}

func createCoverage(w http.ResponseWriter, r *http.Request) {
	org := r.URL.Query()["org"][0]
	if org == "" {
		handleError(w, http.StatusBadRequest, errors.New("Requires an org query parameter"))
	}

	var body CoverageData
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		handleError(w, http.StatusBadRequest, errors.New("Invalid JSON"))
		return
	}

	if err := body.Validate(); err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}

	orgSnapshot, err := db.Collection("coverages").Doc(org).Get(ctx)
	if err != nil {
		handleError(w, http.StatusBadRequest, errors.New("Org is not being tracked"))
		return
	}

	doc, wr, err := orgSnapshot.Ref.Collection("coverage").Add(ctx, body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, errors.New("Failed to save coverage"))
		return
	}

	json.NewEncoder(w).Encode(doc)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		getCoverage(w, r)
	case r.Method == "POST":
		createCoverage(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
