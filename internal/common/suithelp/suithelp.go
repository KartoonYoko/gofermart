package suithelp

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresqlContainerData struct {
	Port             uint16
	Host             string
	ConnectionString string
}

type AccrualContainerData struct {
	Port uint16
	Host string
	Tc   testcontainers.Container
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

func NewAccrualContainer(ctx context.Context) (*AccrualContainerData, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	// TODO не находит dockerfile по указанному пути

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       "../../../../cmd/accrual/.",
			Dockerfile:    "../../../../docker/Dockerfile.accrual",
			PrintBuildLog: true,
			KeepImage:     false,
		},
	}

	gReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	container, err := testcontainers.GenericContainer(ctx, gReq)
	if err != nil {
		return nil, err
	}
	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := container.MappedPort(ctx, "8080")
	if err != nil {
		return nil, err
	}

	portDigit := uint16(port.Int())
	return &AccrualContainerData{
		Host: host,
		Port: portDigit,
		Tc:   container,
	}, err
}

func GetPostgresqlContainerData(ctx context.Context, pgc *tcpostgres.PostgresContainer) (*PostgresqlContainerData, error) {
	host, err := pgc.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := pgc.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}
	portDigit := uint16(port.Int())

	cd := &PostgresqlContainerData{
		Port: portDigit,
		Host: host,
		ConnectionString: fmt.Sprintf(
			"host=%s user=postgres password=123 dbname=gophermart port=%d sslmode=disable",
			host, portDigit),
	}
	return cd, nil
}
