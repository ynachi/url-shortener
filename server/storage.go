package server

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
)

// CreateFirestoreClient instanciates a new firestore client. This client should ideally be closed when done with a defer statement.
func CreateFirestoreClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		Logger.Error("unable to create firestore client", err, "project_id", projectID)
		return nil, err
	}
	// Close client when done with
	// defer client.Close()
	return client, nil
}

func SaveURL(longURL, shortID string) error {
	ctx := context.Background()
	projectID := os.Getenv("URL_SHORTNER_PROJECT_ID")
	firestoreClient, err := CreateFirestoreClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer firestoreClient.Close()
	// set document ID to url short ID
	_, err = firestoreClient.Collection("urls").Doc(shortID).Set(ctx, map[string]interface{}{
		"short_id": shortID,
		"long_url": longURL,
	})
	if err != nil {
		Logger.Error("failed to save url", err, "long_url", longURL, "short_url_id", shortID)
		return err
	}
	return nil
}
