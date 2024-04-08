package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/mrofisr/go-testing/internal/model"
	"log"
)

type PersonRepository interface {
	Create(ctx context.Context, name string, age int) error
	FindAll(ctx context.Context) ([]model.Person, error)
	FindById(ctx context.Context, id int) (model.Person, error)
	Update(ctx context.Context, id int, name string, age int) error
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
}

type personRepository struct {
	pgx       *pgx.Conn
	tableName string
}

func (p personRepository) Create(ctx context.Context, name string, age int) error {
	_, err := p.pgx.Exec(ctx, fmt.Sprintf("INSERT INTO %s (name, age) VALUES ($1, $2)", p.tableName), name, age)
	return err
}

func (p personRepository) FindAll(ctx context.Context) ([]model.Person, error) {
	query, err := p.pgx.Query(ctx, fmt.Sprintf("SELECT * FROM %s", p.tableName))
	if err != nil {
		return nil, err
	}

	var persons []model.Person

	for query.Next() {
		var person model.Person
		err := query.Scan(&person.ID, &person.Name, &person.Age) // Assuming these are the columns of your Person table
		if err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}

	if err := query.Err(); err != nil {
		return nil, err
	}

	return persons, nil
}

func (p personRepository) FindById(ctx context.Context, id int) (model.Person, error) {
	var person model.Person
	person.ID = id
	err := p.pgx.QueryRow(ctx, fmt.Sprintf("SELECT name, age FROM %s WHERE id = $1", p.tableName), 1).Scan(&person.Name, &person.Age)
	if err != nil {
		return person, err
	}
	return person, err
}

func (p personRepository) Update(ctx context.Context, id int, name string, age int) error {
	_, err := p.pgx.Exec(context.Background(), fmt.Sprintf("UPDATE %s SET name = $1, age = $2 WHERE id = $3", p.tableName), name, age, id)
	return err

}

func (p personRepository) Delete(ctx context.Context, id int) error {
	_, err := p.pgx.Exec(context.Background(), fmt.Sprintf("DELETE FROM %s WHERE id = $1", p.tableName), id)
	return err
}

func (p personRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := p.pgx.QueryRow(context.Background(), fmt.Sprintf("SELECT COUNT(*) FROM %s", p.tableName)).Scan(&count)
	return count, err
}

func NewPersonRepository(pgx *pgx.Conn) PersonRepository {

	// Create table
	tableName := "users"
	_, err := pgx.Exec(context.Background(), fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id SERIAL PRIMARY KEY, name VARCHAR(50), age INT)", tableName))
	if err != nil {
		log.Fatalf("Unable to create table %s: %v", tableName, err)
	}

	return &personRepository{
		pgx:       pgx,
		tableName: tableName,
	}
}
