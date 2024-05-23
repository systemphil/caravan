package main

import (
	"context"
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

func generateV4GetObjectSignedURL(w io.Writer, bucket, object string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	// Signing a URL requires credentials authorized to sign a URL. You can pass
	// these in through SignedURLOptions with one of the following options:
	//    a. a Google service account private key, obtainable from the Google Developers Console
	//    b. a Google Access ID with iam.serviceAccounts.signBlob permissions
	//    c. a SignBytes function implementing custom signing.
	// In this example, none of these options are used, which means the SignedURL
	// function attempts to use the same authentication that was used to instantiate
	// the Storage client. This authentication must include a private key or have
	// iam.serviceAccounts.signBlob permissions.
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}

	blob := client.Bucket(bucket).Object(object)

	// Check if the object exists
	_, err = blob.Attrs(ctx)

	if err != nil {
		fmt.Fprintf(w, "No object found: %s", http.StatusText(http.StatusNotFound))
		return "", nil
	}

	u, err := client.Bucket(bucket).SignedURL(object, opts)
	if err != nil {
		return "Error creating signed URL", fmt.Errorf("Bucket(%q).SignedURL: %w", bucket, err)
	}

	fmt.Fprintln(w, "Generated GET signed URL:")
	fmt.Fprintf(w, "%q\n", u)
	fmt.Fprintln(w, "You can use this URL with any user agent, for example:")
	fmt.Fprintf(w, "curl %q\n", u)
	return u, nil
}
func main() {
	/*
		Validate and load environment variables
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
		Start server
	*/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		object := r.URL.Path[1:] // Get the object from the request URL path
		if object == "" {
			http.Error(w, "Object not found.", http.StatusNotFound)
			return
		}
		generateV4GetObjectSignedURL(w, bucket, object)
	})
	fmt.Println("[Server] Listening on port 8080")
	http.ListenAndServe(":"+port, nil)
}
