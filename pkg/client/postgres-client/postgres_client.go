package postgres_client

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type PostgresClient struct {
	db *sql.DB
}

func GetPostgresClient(connectionString string) *PostgresClient {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal().Msgf("client: could not configure postgres client %v", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().Msgf("client: could not verify connection of postgres client %v", err.Error())
	}
	return &PostgresClient{
		db: db,
	}
}

func (client *PostgresClient) Close() {
	if client.db != nil {
		err := client.db.Close()
		if err != nil {
			log.Error().Msgf("client: error occured while attempting to close postgres connection %v", err.Error())
			return
		}
	}
}

/* Query with resulted rows, i.e GET */
func (client *PostgresClient) ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
    rows, err := client.db.Query(query, args...)
    if err != nil {
		log.Error().Msgf("client: error occured while executing postgres query %v", err.Error())
        return nil, err
    }
    return rows, nil
}

/* Query with non resulted rows, i.e INSERT, UPDATE, DELETE */
func (client *PostgresClient) ExecuteNonQuery(query string, args ...interface{}) (sql.Result, error) {
    result, err := client.db.Exec(query, args...)
    if err != nil {
		log.Error().Msgf("client: error occured while executing postgres non query %v", err.Error())
        return nil, err
    }
    return result, nil
}
