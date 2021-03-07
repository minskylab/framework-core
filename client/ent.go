package client

import (
	"context"
	"fmt"
	"opencensus/core/ent"
	"opencensus/core/ent/migrate"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type pgConfig struct {
	hostname string
	port     string
	dbName   string
	user     string
	password string
}

func getPGConfig() pgConfig {
	_ = godotenv.Load() // load envs from .env file (create your own).

	hostname := os.Getenv("PG_HOSTNAME")
	port := os.Getenv("PG_PORT")
	dbName := os.Getenv("PG_DATABASE")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")

	return pgConfig{
		hostname: hostname,
		port:     port,
		dbName:   dbName,
		user:     user,
		password: password,
	}
}

// NewClient returns a new ent client.
func NewClient() (*ent.Client, error) {
	pgConfig := getPGConfig()

	key := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		pgConfig.hostname,
		pgConfig.port,
		pgConfig.user,
		pgConfig.dbName,
		pgConfig.password,
	)

	fmt.Println(key)

	client, _ := ent.Open("postgres", key)
	defer client.Close()

	ctx := context.Background()

	err := client.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)

	// TODO: Implement entgo auto migration and return ready client
	return client, err
}
