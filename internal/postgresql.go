package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	testcontainers.Container
	dbName   string
	user     string
	password string
	dsn      string
}

func (c *PostgresContainer) ConnectionString() string {
	return c.dsn
}

func RunContainer(ctx context.Context) (*PostgresContainer, error) {
	var connStr string

	req := testcontainers.ContainerRequest{
		Image: "postgres:15",
		Env: map[string]string{
			"POSTGRES_USER":     "rotabot",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "rotabot",
		},
		ExposedPorts: []string{"5432/tcp"},
		Cmd:          []string{"postgres", "-c", "fsync=off"},
		WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
			connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", "rotabot", "password", host, port.Port(), "rotabot", "sslmode=disable")
			return connStr
		}).WithStartupTimeout(time.Second * 5),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, err
	}

	user := req.Env["POSTGRES_USER"]
	password := req.Env["POSTGRES_PASSWORD"]
	dbName := req.Env["POSTGRES_DB"]

	return &PostgresContainer{Container: container, dbName: dbName, password: password, user: user, dsn: connStr}, nil
}
