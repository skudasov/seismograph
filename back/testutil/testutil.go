package testutil

import (
	"fmt"

	"github.com/f4hrenh9it/seismograph/back/config"
)

func PgConnectionString(cfg config.DBConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.User,
		cfg.Postgres.Pwd,
		cfg.Postgres.DBName,
		cfg.Postgres.Port,
	)
}
