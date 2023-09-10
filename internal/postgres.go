package internal

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultUser          = "rotabot"
	defaultPassword      = "password"
	defaultPostgresImage = "docker.io/postgres:14-alpine"
)

type PostgresContainer struct {
	testcontainers.Container
	dbName   string
	user     string
	password string
}

// ConnectionString returns the connection string for the postgres container, using the default 5432 port, and
// obtaining the host and exposed port from the container. It also accepts a variadic list of extra arguments
// which will be appended to the connection string. The format of the extra arguments is the same as the
// connection string format, e.g. "connect_timeout=10" or "application_name=myapp"
func (c *PostgresContainer) ConnectionString(ctx context.Context, args ...string) (string, error) {
	containerPort, err := c.Container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return "", err
	}

	host, err := c.Container.Host(ctx)
	if err != nil {
		return "", err
	}

	extraArgs := strings.Join(args, "&")
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?%s", c.user, c.password, net.JoinHostPort(host, containerPort.Port()), c.dbName, extraArgs)
	return connStr, nil
}

func (c *PostgresContainer) Terminate(ctx context.Context) error {
	return c.Container.Terminate(ctx)
}

func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*PostgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Image: defaultPostgresImage,
		Env: map[string]string{
			"POSTGRES_USER":     defaultUser,
			"POSTGRES_PASSWORD": defaultPassword,
			"POSTGRES_DB":       defaultUser,
		},
		ExposedPorts: []string{"5432/tcp"},
		Cmd:          []string{"postgres", "-c", "fsync=off"},
		WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
			return fmt.Sprintf("postgres://%s:%s@%s/%s?%s", "rotabot", "password", net.JoinHostPort(host, port.Port()), "rotabot", "sslmode=disable")
		}).WithStartupTimeout(time.Second * 5),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	for _, opt := range opts {
		opt.Customize(&genericContainerReq)
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, err
	}

	user := req.Env["POSTGRES_USER"]
	password := req.Env["POSTGRES_PASSWORD"]
	dbName := req.Env["POSTGRES_DB"]

	return &PostgresContainer{Container: container, dbName: dbName, password: password, user: user}, nil
}
