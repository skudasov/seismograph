package testutil

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func RunMinio() (*dockertest.Pool, *dockertest.Resource, string) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "minio/minio",
		Tag:        "RELEASE.2020-12-10T01-54-29Z",
		PortBindings: map[docker.Port][]docker.PortBinding{
			"9000/tcp": {{HostPort: "9000"}},
		},
		Cmd: []string{"server", "/data/minio"},
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	mappedPort := resource.GetPort("9000/tcp")
	endpoint := fmt.Sprintf("localhost:%s", mappedPort)
	if err := pool.Retry(func() error {
		url := fmt.Sprintf("http://%s/minio/health/live", endpoint)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status code not OK")
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return pool, resource, mappedPort
}

func RunPGSQL() (*dockertest.Pool, *dockertest.Resource, string) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "13.0", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=testdb"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	mappedPort := resource.GetPort("5432/tcp")

	if err := pool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable", mappedPort, "testdb"))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return pool, resource, mappedPort
}
