package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
)

// bucket := "symposia-dev-bucket"
//
//	object := "cluvqhyly0007uwfdmg2hn33a/VID_20200103_135115.mp4"
//
// generateV4GetObjectSignedURL generates object signed URL with GET method.
func generateV4GetObjectSignedURL(w io.Writer, bucket, object string) (string, error) {
	// bucket := "bucket-name"
	// object := "object-name"

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

	u, err := client.Bucket(bucket).SignedURL(object, opts)
	if err != nil {
		return "", fmt.Errorf("Bucket(%q).SignedURL: %w", bucket, err)
	}

	fmt.Fprintln(w, "Generated GET signed URL:")
	fmt.Fprintf(w, "%q\n", u)
	fmt.Fprintln(w, "You can use this URL with any user agent, for example:")
	fmt.Fprintf(w, "curl %q\n", u)
	return u, nil
}
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bucket := "symposia-dev-bucket"
		object := "video/cluvqhyly0007uwfdmg2hn33a/VID_20200103_135115.mp4"
		generateV4GetObjectSignedURL(w, bucket, object)
	})
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
