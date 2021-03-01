package server

import (
	"log"
	"os"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiber_logger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"

	"github.com/f4hrenh9it/seismograph/back/config"
	"github.com/f4hrenh9it/seismograph/back/testutil"
)

type SeismographService struct {
	Cfg *config.Config
	PG  *gorm.DB
	M   *minio.Client
}

const (
	URLHandleTest                = "/project/:project/test"
	URLHandleDeleteTest          = "/project/:project/test/:id"
	URLHandleTests               = "/project/:project/tests"
	URLHandleGetTestDataId       = "/project/:project/test/:id"
	URLHandleCreateProject       = "/project"
	URLHandleDeleteProject       = "/project/:project"
	URLHandleProjects            = "/projects"
	URLHandleCreateAttackCluster = "/attack_cluster"
	URLHandleGetAttackCluster    = "/attack_cluster/:id"
	URLHandleGetAttackClusterVMS = "/attack_cluster/:id/vms"
	URLHandleGetAttackClusters   = "/attack_clusters"
	URLHandleDeleteAttackCluster = "/attack_cluster/:id"
	URLHandleCreateAttackVM      = "/vm"
	URLHandleGetAttackVM         = "/vm/:id"
	URLHandleDeleteAttackVM      = "/vm/:id"
)

func NewSeismographService(cfg *config.Config) {
	dsn := testutil.PgConnectionString(cfg.DB)
	db, err := gorm.Open(postgres.Open(dsn),
		&gorm.Config{Logger: gorm_logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gorm_logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      gorm_logger.Info,
				Colorful:      true,
			},
		)})
	if err != nil {
		log.Fatal(err)
	}
	mi, err := minio.New(cfg.DB.Minio.Url, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.DB.Minio.AccessKey, cfg.DB.Minio.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	s := &SeismographService{
		Cfg: cfg,
		PG:  db,
		M:   mi,
	}

	app := fiber.New(fiber.Config{
		BodyLimit: cfg.Server.BodyLimit,
	})
	if cfg.Server.Pprof {
		app.Use(pprof.New())
	}
	if cfg.Server.Prometheus {
		prometheus := fiberprometheus.New("seismographd")
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}
	logCfg := fiber_logger.Config{
		Next:         nil,
		Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat:   "15:04:05",
		TimeZone:     "Local",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stdout,
	}
	app.Use(fiber_logger.New(logCfg))
	app.Use(cors.New())

	// Test
	app.Post(URLHandleTest, s.HandleTest)
	app.Get(URLHandleGetTestDataId, s.HandleGetTestDataId)
	app.Delete(URLHandleDeleteTest, s.HandleDeleteTest)
	app.Get(URLHandleTests, s.HandleTestsForProject)

	// Project
	app.Post(URLHandleCreateProject, s.HandleCreateProject)
	app.Get(URLHandleProjects, s.HandleProjects)
	app.Delete(URLHandleDeleteProject, s.HandleDeleteProjectWithTests)

	// AttackCluster
	app.Post(URLHandleCreateAttackCluster, s.HandleAttackCluster)
	app.Get(URLHandleGetAttackCluster, s.HandleGetAttackCluster)
	app.Get(URLHandleGetAttackClusters, s.HandleGetAttackClusters)
	app.Get(URLHandleGetAttackClusterVMS, s.HandleGetAttackClusterVMS)
	app.Delete(URLHandleDeleteAttackCluster, s.HandleDeleteAttackCluster)

	// AttackVM
	app.Post(URLHandleCreateAttackVM, s.HandleAttackVM)
	app.Get(URLHandleGetAttackVM, s.HandleGetAttackVM)
	app.Delete(URLHandleDeleteAttackVM, s.HandleDeleteAttackVM)

	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
