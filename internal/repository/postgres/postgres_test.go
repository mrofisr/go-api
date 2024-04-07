package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgreSQLDB(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "postgres:13-alpine",
		Name:  "postgresql-test",
		Env: map[string]string{
			"POSTGRES_USER":     "username",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "database_name",
		},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		ExposedPorts: []string{"5432/tcp", "5432/tcp"},
		Networks:     []string{"bridge"},
	}
	// Start a PostgreSQL container
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to start PostgreSQL container: %v", err)
	}
	defer postgresContainer.Terminate(ctx)

	return postgresContainer, nil
}

// TestCRUD is a simple test to test CRUD operations
func TestCRUD(t *testing.T) {
	ctx := context.Background()
	postgresContainer, err := setupPostgreSQLDB(ctx)
	if err != nil {
		t.Fatalf("Unable to setup PostgreSQL database: %v", err)
	}
	containerHost, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Unable to get container host: %v", err)
	}
	containerPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Unable to get container port: %v", err)
	}
	fmt.Printf("Host: %s, Port: %s\n", containerHost, containerPort.Port())
	// Open a connection to the PostgreSQL database
	conn, err := pgx.Connect(context.Background(), "postgresql://username:password@"+containerHost+":"+containerPort.Port()+"/database_name")
	if err != nil {
		t.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Create table
	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name VARCHAR(50), age INT)")
	if err != nil {
		t.Fatalf("Unable to create table: %v", err)
	}

	// Insert data
	_, err = conn.Exec(context.Background(), "INSERT INTO users (name, age) VALUES ($1, $2)", "Alice", 30)
	if err != nil {
		t.Fatalf("Unable to insert data: %v", err)
	}

	// Read data
	var name string
	var age int
	err = conn.QueryRow(context.Background(), "SELECT name, age FROM users WHERE id = $1", 1).Scan(&name, &age)
	if err != nil {
		t.Fatalf("Unable to read data: %v", err)
	}
	if name != "Alice" || age != 30 {
		t.Fatalf("Expected name: %s, age: %d. Got name: %s, age: %d", "Alice", 30, name, age)
	}

	// Update data
	_, err = conn.Exec(context.Background(), "UPDATE users SET age = $1 WHERE id = $2", 31, 1)
	if err != nil {
		t.Fatalf("Unable to update data: %v", err)
	}

	// Read updated data
	err = conn.QueryRow(context.Background(), "SELECT name, age FROM users WHERE id = $1", 1).Scan(&name, &age)
	if err != nil {
		t.Fatalf("Unable to read updated data: %v", err)
	}
	if age != 31 {
		t.Fatalf("Expected age: %d. Got age: %d", 31, age)
	}

	// Delete data
	_, err = conn.Exec(context.Background(), "DELETE FROM users WHERE id = $1", 1)
	if err != nil {
		t.Fatalf("Unable to delete data: %v", err)
	}

	// Verify deletion
	var count int
	err = conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("Unable to verify deletion: %v", err)
	}
	if count != 0 {
		t.Fatalf("Expected count: %d. Got count: %d", 0, count)
	}
}
