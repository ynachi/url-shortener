package server

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/go-redis/redis/v8"
)

// newFirestoreClient instantiates a new firestore client. This client should
// ideally be closed when done with a defer statement.
func newFirestoreClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		Logger.Error("unable to create firestore client", err, "project_id", projectID)
		return nil, ErrFirestoreClientCreate
	}
	// Close client when done with
	// defer client.Close()
	return client, nil
}

// newRedisClient NewRedisPool creates a new redis connexion pool. Initialize the
// pool at package level to maintain a single pool other on the whole application.
// Init done in server.go.
func newRedisClient(redisAddr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		Logger.Error("unable to connect to redis server", err, "redis_server", redisAddr)
		return nil, ErrRedisClientCreate
	}
	Logger.Info("connexion to redis server was successful", "redis_server", redisAddr)
	// Close client when done with
	// defer client.Close()
	return client, nil
}

// persistURL save the long url along with it's shortID in the Firestore database.
func persistURL(ctx context.Context, longURL, shortID string, firestoreClient *firestore.Client) error {
	_, err := firestoreClient.Collection("urls").Doc(shortID).Set(ctx, map[string]interface{}{
		"long_url": longURL,
	})
	if err != nil {
		Logger.Error("failed to save url", err, "long_url", longURL, "short_url_id", shortID)
		return err
	}
	return nil
}

// getFromCache retrieves long url matching the given ID from the caching servers
// Returns redis.Nil error type if the key does not exist
func getFromCache(ctx context.Context, shortID string, redisClient *redis.Client) (string, error) {
	longURL, err := redisClient.Get(ctx, shortID).Result()
	if err != nil {
		Logger.Error("failed to get data to cache server", err)
		return "", err
	}
	return longURL, nil
}

// getFromStorage retrieves long url matching the given ID from the storage servers
func getFromStorage(ctx context.Context, shortID string, firestoreClient *firestore.Client) (string, error) {
	data, err := firestoreClient.Collection("urls").Doc(shortID).Get(ctx)
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

// SaveToCache saves long url matching the given ID to the caching servers
func saveToCache(ctx context.Context, shortID string, longURL string, redisClient *redis.Client) error {
	err := redisClient.Set(ctx, shortID, longURL, CacheDuration).Err()
	if err != nil {
		Logger.Error("failed to save data to cache server", err)
		return err
	}
	return nil
}
