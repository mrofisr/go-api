package main

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestCRUD is a simple test to test CRUD operations
func TestCRUD(t *testing.T) {
	ctx := context.Background()

	dbName := "postgres"
	dbUser := "user"
	dbPassword := "password"

	// 1. Start the postgres container and run any migrations on it
	container, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Run any migrations on the database
	_, _, err = container.Exec(ctx, []string{"psql", "-U", dbUser, "-d", dbName, "-c", "CREATE TABLE users (id SERIAL, name TEXT NOT NULL, age INT NOT NULL)"})
	if err != nil {
		t.Fatal(err)
	}

	// 2. Create a snapshot of the database to restore later
	err = container.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	dbURL, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("all", func(t *testing.T) {

		t.Cleanup(func() {
			// 3. In each test, reset the DB to its snapshot state.
			err = container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		conn, err := pgx.Connect(context.Background(), dbURL)
		if err != nil {
			t.Fatal(err)
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
	})
}
