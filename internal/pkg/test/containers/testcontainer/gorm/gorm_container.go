package gorm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/contracts"
)

// https://github.com/testcontainers/testcontainers-go/issues/1359
// https://github.com/testcontainers/testcontainers-go/issues/1249
type gormTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.PostgresContainerOptions
}

func NewGormTestContainers() contracts.GormContainer {
	return &gormTestContainers{
		defaultOptions: &contracts.PostgresContainerOptions{
			Database:  "test_db",
			Port:      "5432/tcp",
			Host:      "localhost",
			UserName:  "testcontainers",
			Password:  "testcontainers",
			Tag:       "latest",
			ImageName: "postgres",
			Name:      "postgresql-testcontainer",
		},
	}
}

func (g *gormTestContainers) CreatingContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.PostgresContainerOptions,
) (*gormPostgres.GormOptions, error) {
	// https://github.com/testcontainers/testcontainers-go
	// https://dev.to/remast/go-integration-tests-using-testcontainers-9o5
	containerReq := g.getRunOptions(options...)

	// TODO: Using Parallel Container
	dbContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return nil, err
	}

	// get a free random host hostPort
	hostPort, err := dbContainer.MappedPort(ctx, nat.Port(g.defaultOptions.Port))
	if err != nil {
		return nil, err
	}
	g.defaultOptions.HostPort = hostPort.Int()

	host, err := dbContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	g.container = dbContainer

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := dbContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	gormOptions := &gormPostgres.GormOptions{
		Port:     g.defaultOptions.HostPort,
		Host:     host,
		Password: g.defaultOptions.Password,
		DBName:   g.defaultOptions.Database,
		SSLMode:  false,
		User:     g.defaultOptions.UserName,
	}
	return gormOptions, nil
}

func (g *gormTestContainers) Start(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.PostgresContainerOptions,
) (*gorm.DB, error) {
	gormOptions, err := g.CreatingContainerOptions(ctx, t, options...)
	if err != nil {
		return nil, err
	}

	db, err := gormPostgres.NewGorm(gormOptions)

	return db, nil
}

func (g *gormTestContainers) Cleanup(ctx context.Context) error {
	return g.container.Terminate(ctx)
}

func (g *gormTestContainers) getRunOptions(
	opts ...*contracts.PostgresContainerOptions,
) testcontainers.ContainerRequest {
	if len(opts) > 0 && opts[0] != nil {
		option := opts[0]
		if option.ImageName != "" {
			g.defaultOptions.ImageName = option.ImageName
		}
		if option.Host != "" {
			g.defaultOptions.Host = option.Host
		}
		if option.Port != "" {
			g.defaultOptions.Port = option.Port
		}
		if option.UserName != "" {
			g.defaultOptions.UserName = option.UserName
		}
		if option.Password != "" {
			g.defaultOptions.Password = option.Password
		}
		if option.Tag != "" {
			g.defaultOptions.Tag = option.Tag
		}
	}

	//hostFreePort, err := freeport.GetFreePort()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//g.defaultOptions.HostPort = hostFreePort

	containerReq := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", g.defaultOptions.ImageName, g.defaultOptions.Tag),
		ExposedPorts: []string{g.defaultOptions.Port},
		WaitingFor:   wait.ForListeningPort(nat.Port(g.defaultOptions.Port)).WithPollInterval(2 * time.Second),
		Hostname:     g.defaultOptions.Host,
		SkipReaper:   true,
		Env: map[string]string{
			"POSTGRES_DB":       g.defaultOptions.Database,
			"POSTGRES_PASSWORD": g.defaultOptions.Password,
			"POSTGRES_USER":     g.defaultOptions.UserName,
		},
	}

	return containerReq
}
