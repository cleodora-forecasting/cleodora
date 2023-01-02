package cleosrv

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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

func Start(address string, database string, frontendFooterText string) error {
	fmt.Printf(`cleosrv (Cleodora server) - Visit https://cleodora.org for more information.
Version: %s
Database: %s
Listening on: %s

`,
		cleoutils.Version,
		database,
		address,
	)

	err := os.MkdirAll(filepath.Dir(database), 0770)
	if err != nil {
		return fmt.Errorf("error making directories for database %v: %w", database, err)
	}

	router := chi.NewRouter()

	db, err := InitDB(database)
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

	serveFrontend(router, frontendFooterText)

	return http.ListenAndServe(address, router)
}

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
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
