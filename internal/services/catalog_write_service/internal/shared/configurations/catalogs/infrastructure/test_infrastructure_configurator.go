package infrastructure

//
//import (
//	"context"
//	"testing"
//
//	"emperror.dev/errors"
//	"github.com/go-playground/validator"
//	"gorm.io/gorm"
//
//	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer/json"
//	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
//	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
//	otelMetrics "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics"
//	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
//	gorm2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/gorm"
//
//	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/config"
//	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/mocks/testData"
//	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/models"
//	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/contracts"
//)
//
//type testInfrastructureConfigurator struct {
//	log logger.Logger
//	cfg *config.Config
//	t   *testing.T
//}
//
//func NewTestInfrastructureConfigurator(
//	t *testing.T,
//	log logger.Logger,
//	cfg *config.Config,
//) contracts.InfrastructureConfigurator {
//	return &testInfrastructureConfigurator{log: log, cfg: cfg, t: t}
//}
//
//func (ic *testInfrastructureConfigurator) ConfigInfrastructures(
//	ctx context.Context,
//) (*contracts.InfrastructureConfigurations, func(), error) {
//	infrastructure := &contracts.InfrastructureConfigurations{
//		Cfg:       ic.cfg,
//		Log:       ic.log,
//		Validator: validator.New(),
//	}
//
//	meter, err := otelMetrics.NewOtelMetrics(ctx, ic.cfg.OTelMetricsConfig, ic.log)
//	if err != nil {
//		return nil, nil, err
//	}
//	infrastructure.Metrics = meter
//
//	var cleanup []func()
//
//	grpcClient, err := grpc.NewGrpcClient(ic.cfg.GRPC)
//	if err != nil {
//		return nil, nil, err
//	}
//	cleanup = append(cleanup, func() {
//		_ = grpcClient.Close()
//	})
//	infrastructure.GrpcClient = grpcClient
//
//	gorm, err := gorm2.NewGormTestContainers().Start(ctx, ic.t)
//	if err != nil {
//		return nil, nil, err
//	}
//	infrastructure.Gorm = gorm
//
//	err = seedGormAndMigration(gorm)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	traceProvider, err := tracing.NewOtelTracing(ic.cfg.OTel)
//	if err != nil {
//		return nil, nil, err
//	}
//	cleanup = append(cleanup, func() {
//		_ = traceProvider.Shutdown(ctx)
//	})
//
//	infrastructure.EventSerializer = json.NewEventSerializer()
//
//	return infrastructure, func() {
//		for _, c := range cleanup {
//			c()
//		}
//	}, nil
//}
//
//func seedGormAndMigration(gormDB *gorm.DB) error {
//	// migration
//	err := gormDB.AutoMigrate(models.Product{})
//	if err != nil {
//		return errors.Wrap(err, "error in seed database")
//	}
//
//	// seed data
//	err = gormDB.CreateInBatches(testData.Products, len(testData.Products)).Error
//	if err != nil {
//		return errors.Wrap(err, "error in seed database")
//	}
//	return nil
//}
