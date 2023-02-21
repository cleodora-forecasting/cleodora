package cleosrv

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/cleodora-forecasting/cleodora/cleosrv/dbmodel"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/generated"
	"github.com/cleodora-forecasting/cleodora/cleoutils"
)

type App struct {
	Out    io.Writer
	Err    io.Writer
	Config *Config
}

func NewApp() *App {
	c := &Config{}
	return &App{
		Out:    os.Stdout,
		Err:    os.Stderr,
		Config: c,
	}
}

func (a *App) Version() error {
	if _, err := fmt.Fprintf(a.Out, "%v\n", cleoutils.Version); err != nil {
		return err
	}
	return nil
}

func (a *App) Start() error {
	fmt.Printf(`cleosrv (Cleodora server) - Visit https://cleodora.org for more information.
Version: %s
Database: %s
Listening on: %s

`,
		cleoutils.Version,
		a.Config.Database,
		a.Config.Address,
	)

	err := os.MkdirAll(filepath.Dir(a.Config.Database), 0770)
	if err != nil {
		return fmt.Errorf("error making directories for database %v: %w", a.Config.Database, err)
	}

	router := chi.NewRouter()

	db, err := a.InitDB()
	if err != nil {
		return err
	}

	resolver := graph.NewResolver(db)
	err = resolver.AddDummyData()
	if err != nil {
		return err
	}

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: resolver}),
	)

	configureCORS(router, srv)

	router.Handle("/playground/",
		playground.Handler("GraphQL playground", "/query"),
	)
	router.Handle("/query", srv)

	serveFrontend(router, a.Config.Frontend.FooterText)

	return http.ListenAndServe(a.Config.Address, router)
}

func (a *App) InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(
		sqlite.Open(a.Config.Database),
		&gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		},
	)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&dbmodel.Forecast{},
		&dbmodel.Outcome{},
		&dbmodel.Probability{},
		&dbmodel.Estimate{},
	)
	if err != nil {
		return db, err
	}

	return db, nil
}
