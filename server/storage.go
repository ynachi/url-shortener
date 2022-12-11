package server

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/go-redis/redis/v8"
)

// NewFirestoreClient instantiates a new firestore client. This client should
// ideally be closed when done with a defer statement.
func NewFirestoreClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		Logger.Error("unable to create firestore client", err, "project_id", projectID)
		return nil, err
	}
	// Close client when done with
	// defer client.Close()
	return client, nil
}

// NewRedisClient NewRedisPool creates a new redis connexion pool. Initialize the
// pool at package level to maintain a single pool other on the whole application.
// Init done in server.go.
func NewRedisClient(redisAddr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		Logger.Error("unable to connect to redis server", err, "redis_server", redisAddr)
		return nil, err
	}
	Logger.Info("connexion to redis server was successfully", "redis_server", redisAddr)
	return client, nil
}

// PersistURL save the long url along with it's shortID in the Firestore database.
func PersistURL(longURL, shortID string) error {
	ctx := context.Background()
	projectID := os.Getenv("URL_SHORTNER_PROJECT_ID")
	firestoreClient, err := NewFirestoreClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer firestoreClient.Close()
	// set document ID to url short ID
	_, err = firestoreClient.Collection("urls").Doc(shortID).Set(ctx, map[string]interface{}{
		"long_url": longURL,
	})
	if err != nil {
		Logger.Error("failed to save url", err, "long_url", longURL, "short_url_id", shortID)
		return err
	}
	return nil
}

// GetFromCache retrieves long url matching the given ID from the caching servers
// Returns redis.Nil error type if the key does not exist
func GetFromCache(shortID string) (string, error) {
	longURL, err := redisClient.Get(Ctx, shortID).Result()
	if err != nil {
		Logger.Error("failed to get data to cache server", err)
		return "", err
	}
	return longURL, nil
}

// GetFromStorage retrieves long url matching the given ID from the storage servers
func GetFromStorage(shortID string) (string, error) {
	data, err := firestoreClient.Collection("urls").Doc(shortID).Get(Ctx)
	if err != nil {
		Logger.Error("failed to get Firestore document by ID", err, "short_url", shortID)
		return "", err
	}
	dataInterface, err := data.DataAt("long_url")
	if err != nil {
		Logger.Error("failed to get Firestore document long_url field", err, "short_url", shortID)
		return "", err
	}
	dataStr, ok := dataInterface.(string)
	if !ok {
		Logger.Error("failed to converts data field to string", err, "short_url", shortID)
		return "", err
	}
	return dataStr, nil
}

// SaveToCache saves long url matching the given ID from the caching servers
func SaveToCache(shortID string, longURL string) error {
	err := redisClient.Set(Ctx, shortID, longURL, CacheDuration).Err()
	if err != nil {
		Logger.Error("failed to save data to cache server", err)
		return err
	}
	return nil
}
