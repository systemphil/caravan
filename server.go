package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
)

var REQUIRED_VARS = []string{"PORT", "GCP_PRIMARY_BUCKET_NAME", "GCP_SECONDARY_BUCKET_NAME"}

type StorageClient interface {
	Bucket(name string) *storage.BucketHandle
}

type ReadObjectRequest struct {
	Object string `json:"object"`
}

type WriteObjectRequest struct {
	Object string `json:"object"`
}

func main() {
	/*
		Load and validate environment variables
	*/
	err := godotenv.Load()
	if err != nil {
		log.Printf("[Server] No env file found, attempting to use system environment variables.")
	}
	var missingVars []string
	for _, varName := range REQUIRED_VARS {
		value := os.Getenv(varName)
		if value == "" {
			missingVars = append(missingVars, varName)
		}
	}
	if len(missingVars) > 0 {
		errorMessage := fmt.Sprintf("[Server] [ERROR] Missing required environment variables: %s\nPlease set these variables in your .env file or system environment.", strings.Join(missingVars, ", "))
		log.Fatal(errorMessage)
	}
	port := os.Getenv("PORT")
	bucket := os.Getenv("GCP_PRIMARY_BUCKET_NAME")
	/*
		Initialize Google Cloud Storage client
	*/
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		errorMessage := fmt.Errorf("storage.NewClient: %w", err)
		log.Fatal(errorMessage)
	}
	defer client.Close()
	/*
		Define server routes
	*/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World! From Caravan."))
	})
	http.HandleFunc("/read-object", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Unmarshal the JSON data into a Go struct
		var requestBody ReadObjectRequest
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
			return
		}
		object := requestBody.Object
		generateV4GetObjectSignedURL(w, bucket, object, client, ctx)
	})
	http.HandleFunc("/write-object", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Unmarshal the JSON data into a Go struct
		var requestBody WriteObjectRequest
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
			return
		}
		object := requestBody.Object
		generateV4PutObjectSignedURL(w, bucket, object, client)
	})
	http.HandleFunc("/delete-object", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Unmarshal the JSON data into a Go struct
		var requestBody WriteObjectRequest
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
			return
		}
		object := requestBody.Object
		deleteObject(w, bucket, object, client, ctx)
	})
	/*
		Start server
	*/
	fmt.Println("[Server] Listening on port 8080")
	http.ListenAndServe(":"+port, nil)
}

func generateV4GetObjectSignedURL(w http.ResponseWriter, bucket, object string, client StorageClient, ctx context.Context) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}
	blob := client.Bucket(bucket).Object(object)

	// Check if the object exists
	_, err := blob.Attrs(ctx)
	if err != nil {
		fmt.Fprintf(w, "No object found: %s", http.StatusText(http.StatusNotFound))
		return "", nil
	}

	u, err := client.Bucket(bucket).SignedURL(object, opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "Error creating signed URL", fmt.Errorf("Bucket(%q).SignedURL: %w", bucket, err)
	}

	response := struct {
		URL string `json:"url"`
	}{
		URL: u,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	return u, nil
}

func generateV4PutObjectSignedURL(w http.ResponseWriter, bucket, object string, client StorageClient) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "PUT",
		Expires: time.Now().Add(15 * time.Minute),
	}

	u, err := client.Bucket(bucket).SignedURL(object, opts)
	if err != nil {
		return "", fmt.Errorf("Bucket(%q).SignedURL: %w", bucket, err)
	}

	response := struct {
		URL string `json:"url"`
	}{
		URL: u,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	return u, nil
}

func deleteObject(w http.ResponseWriter, bucket, object string, client StorageClient, ctx context.Context) (string, error) {
	err := client.Bucket(bucket).Object(object).Delete(ctx)
	if err != nil {
		return "", fmt.Errorf("Bucket(%q).SignedURL: %w", bucket, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Object deleted successfully"))

	return "", nil
}
