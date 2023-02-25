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

	err = migrateDB(db)
	if err != nil {
		return nil, fmt.Errorf("migrating data: %w", err)
	}

	return db, nil
}

// migrateDB creates the DB tables and migrates the schema and data to the
// current version if necessary.
func migrateDB(db *gorm.DB) error {
	type Migrations struct {
		ID      string `gorm:"uniqueIndex"`
		Applied time.Time
	}
	dbIsEmpty := !db.Migrator().HasTable("forecasts")
	err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Migrator().AutoMigrate(&Migrations{})
	})
	if err != nil {
		return fmt.Errorf("auto migrating 'migrations' table: %w", err)
	}
	if dbIsEmpty {
		err := db.Transaction(func(tx *gorm.DB) error {
			err := createDb(tx)
			if err != nil {
				return fmt.Errorf("creating tables: %w", err)
			}
			// save all migrations as done in the DB, without executing the
			// functions.
			var mEntries []Migrations
			for _, m := range dbMigrations {
				mEntries = append(
					mEntries,
					Migrations{
						ID:      m.ID,
						Applied: time.Now().UTC(),
					},
				)
			}
			ret := tx.Create(mEntries)
			if ret.Error != nil {
				return fmt.Errorf("storing migrations as done in DB: %w", ret.Error)
			}
			return nil
		})
		return err
	}
	for _, m := range dbMigrations {
		var count int64
		ret := db.Model(&Migrations{}).Where("id = ?", m.ID).Count(&count)
		if ret.Error != nil {
			return fmt.Errorf("selecting migration %v: %w", m.ID, ret.Error)
		}
		if count == 1 {
			continue // migration already ran in the past
		}
		err := db.Transaction(func(tx *gorm.DB) error {
			fmt.Printf("Running DB migration '%v'\n", m.ID)
			if m.Up != nil {
				err = m.Up(tx)
				if err != nil {
					return fmt.Errorf("running %v: %w", m.ID, err)
				}
			}
			ret = tx.Create(Migrations{
				ID:      m.ID,
				Applied: time.Now().UTC(),
			})
			if ret.Error != nil {
				return fmt.Errorf("saving migration %v: %w", m.ID, ret.Error)
			}
			fmt.Printf("Finished DB migration '%v'\n", m.ID)
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func createDb(db *gorm.DB) error {
	return db.AutoMigrate(
		&dbmodel.Forecast{},
		&dbmodel.Outcome{},
		&dbmodel.Probability{},
		&dbmodel.Estimate{},
	)
}

type dbMigration struct {
	ID string
	Up func(db *gorm.DB) error
}

// dbMigrations is a list of all the migrations that need to be applied to an
// old DB to update to the current version. They are run in order. Append new
// ones to the end.
var dbMigrations = []dbMigration{
	{
		ID: "0.1.0 init",
		Up: nil,
	},
	{
		ID: "0.2.0 convert forecast dates",
		Up: func(db *gorm.DB) error {
			type Forecast struct {
				gorm.Model
				Title       string
				Description string
				Created     time.Time
				Resolves    time.Time
				Closes      *time.Time
			}
			var forecasts []Forecast
			ret := db.Find(&forecasts)
			if ret.Error != nil {
				return fmt.Errorf("getting forecasts: %w", ret.Error)
			}
			for _, f := range forecasts {
				f.Created = f.Created.UTC()
				f.Resolves = f.Resolves.UTC()
				f.CreatedAt = f.CreatedAt.UTC()
				if f.Closes != nil && f.Closes.IsZero() {
					f.Closes = nil
				} else if f.Closes != nil {
					temp := f.Closes.UTC()
					f.Closes = &temp
				}
				if f.Resolves.Before(f.Created) {
					f.Resolves = f.Created
				}
				if f.Closes != nil {
					if f.Closes.Before(f.Created) {
						newCloses := f.Created
						f.Closes = &newCloses
					}
					if f.Resolves.Before(*f.Closes) {
						newCloses := f.Resolves
						f.Closes = &newCloses
					}
				}
				ret = db.Save(f)
				if ret.Error != nil {
					return fmt.Errorf("saving %v: %w", f.ID, ret.Error)
				}
			}
			return nil
		},
	},
	{
		ID: "0.2.0 convert estimate dates to UTC",
		Up: func(db *gorm.DB) error {
			type Estimate struct {
				gorm.Model
				Created time.Time
			}
			var estimates []Estimate
			ret := db.Find(&estimates)
			if ret.Error != nil {
				return fmt.Errorf("getting estimates: %w", ret.Error)
			}
			for _, e := range estimates {
				e.Created = e.Created.UTC()
				e.CreatedAt = e.CreatedAt.UTC()
				ret = db.Save(e)
				if ret.Error != nil {
					return fmt.Errorf("saving %v: %w", e.ID, ret.Error)
				}
			}
			return nil
		},
	},
	{
		ID: "0.2.0 convert outcome dates to UTC",
		Up: func(db *gorm.DB) error {
			type Outcome struct {
				gorm.Model
			}
			var outcomes []Outcome
			ret := db.Find(&outcomes)
			if ret.Error != nil {
				return fmt.Errorf("getting outcomes: %w", ret.Error)
			}
			for _, o := range outcomes {
				o.CreatedAt = o.CreatedAt.UTC()
				ret = db.Save(o)
				if ret.Error != nil {
					return fmt.Errorf("saving %v: %w", o.ID, ret.Error)
				}
			}
			return nil
		},
	},
	{
		ID: "0.2.0 convert probability dates to UTC",
		Up: func(db *gorm.DB) error {
			type Probability struct {
				gorm.Model
			}
			var probabilities []Probability
			ret := db.Find(&probabilities)
			if ret.Error != nil {
				return fmt.Errorf("getting probabilities: %w", ret.Error)
			}
			for _, p := range probabilities {
				p.CreatedAt = p.CreatedAt.UTC()
				ret = db.Save(p)
				if ret.Error != nil {
					return fmt.Errorf("saving %v: %w", p.ID, ret.Error)
				}
			}
			return nil
		},
	},
}
