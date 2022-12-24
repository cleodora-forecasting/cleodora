package cleosrv

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/cleodora-forecasting/cleodora/cleosrv/dbmodel"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph/generated"
	"github.com/cleodora-forecasting/cleodora/cleoutils"
)

func Start(address string, frontendFooterText string) error {
	fmt.Printf(
		"Starting Cleodora (version: %s) http://%s\n",
		cleoutils.Version,
		address,
	)

	router := chi.NewRouter()

	db, err := InitDB("test.db")
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
