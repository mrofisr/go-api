package repository

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPersonRepository(t *testing.T) {
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

	p := NewPersonRepository(conn)

	personName := "Alice"
	personAge := 30
	newPersonName := "Bob"
	newPersonAge := 31

	t.Run("insert-data", func(t *testing.T) {
		err = p.Create(context.Background(), personName, personAge)
		if err != nil {
			t.Fatalf("Unable to insert data: %v", err)
		}
	})

	t.Run("read-data", func(t *testing.T) {
		person, err := p.FindById(context.Background(), 1)
		if err != nil {
			t.Fatalf("Unable to read data: %v", err)
		}
		if person.Name != personName || person.Age != personAge {
			t.Fatalf("Expected name: %s, age: %d. Got name: %s, age: %d", personName, personAge, person.Name, person.Age)
		}
	})

	t.Run("read-all-data", func(t *testing.T) {
		persons, err := p.FindAll(context.Background())
		if err != nil {
			t.Fatalf("Unable to read all data: %v", err)
		}
		if len(persons) != 1 || persons[0].Name != personName || persons[0].Age != personAge {
			t.Fatalf("Expected name: %s, age: %d. Got name: %s, age: %d", personName, personAge, persons[0].Name, persons[0].Age)
		}
	})

	t.Run("update-data", func(t *testing.T) {
		err = p.Update(context.Background(), 1, newPersonName, newPersonAge)
		if err != nil {
			t.Fatalf("Unable to update data: %v", err)
		}
	})

	t.Run("read-updated-data", func(t *testing.T) {
		person, err := p.FindById(context.Background(), 1)
		if err != nil {
			t.Fatalf("Unable to read updated data: %v", err)
		}
		if person.Name != newPersonName || person.Age != newPersonAge {
			t.Fatalf("Expected name: %s, age: %d. Got name: %s, age: %d", newPersonName, newPersonAge, person.Name, person.Age)
		}

	})

	t.Run("delete-data", func(t *testing.T) {
		err = p.Delete(context.Background(), 1)
		if err != nil {
			t.Fatalf("Unable to delete data: %v", err)
		}

	})

	t.Run("count-data", func(t *testing.T) {
		count, err := p.Count(context.Background())
		if err != nil {
			t.Fatalf("Unable to count data: %v", err)
		}
		if count != 0 {
			t.Fatalf("Expected count: %d. Got count: %d", 0, count)
		}
	})
}
