package suithelp

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ContainerData struct {
	Port             uint16
	Host             string
	ConnectionString string
}

func NewPostgresContainer(ctx context.Context) (*tcpostgres.PostgresContainer, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	pgc, err := tcpostgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16.2"),
		tcpostgres.WithDatabase("gophermart"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("123"),
		tcpostgres.WithInitScripts(),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)

	return pgc, err
}

func GetPostgresqlContainerData(ctx context.Context, pgc *tcpostgres.PostgresContainer) (*ContainerData, error) {
	host, err := pgc.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := pgc.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}
	portDigit := uint16(port.Int())

	cd := &ContainerData{
		Port: portDigit,
		Host: host,
		ConnectionString: fmt.Sprintf(
			"host=%s user=postgres password=123 dbname=gophermart port=%d sslmode=disable",
			host, portDigit),
	}
	return cd, nil
}
