package cleosrv

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"modernc.org/sqlite"

	"github.com/cleodora-forecasting/cleodora/cleosrv/ent"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/forecast"
	"github.com/cleodora-forecasting/cleodora/cleosrv/graph"
	"github.com/cleodora-forecasting/cleodora/cleoutils"
)

type App struct {
	Out      io.Writer
	Err      io.Writer
	Config   *Config
	dbClient *ent.Client
}

func NewApp(c *Config) (*App, error) {
	app := &App{
		Out:    os.Stdout,
		Err:    os.Stderr,
		Config: c,
	}

	err := os.MkdirAll(filepath.Dir(app.Config.Database), 0770)
	if err != nil {
		return nil, fmt.Errorf("error making directories for database %v: %w", app.Config.Database,
			err)
	}

	// &_fk=1 is missing from the connection string because it's specific to github.
	// com/mattn/go-sqlite3 . Instead the custom driver always enables foreign keys.
	client, err := ent.Open("sqlite3", app.Config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %v", err)
	}

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %v", err)
	}

	app.dbClient = client

	return app, nil
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
	srv := a.GetServer()
	router := chi.NewRouter()
	configureCORS(router, srv)
	router.Handle("/playground/",
		playground.Handler("GraphQL playground", "/query"),
	)
	router.Handle("/query", srv)
	serveFrontend(router, a.Config.Frontend.FooterText)
	return http.ListenAndServe(a.Config.Address, router)
}

func (a *App) Close() error {
	if a.dbClient == nil {
		return nil
	}
	err := a.dbClient.Close()
	if err != nil {
		return fmt.Errorf("error closing DB: %w", err)
	}
	a.dbClient = nil
	return nil
}

// GetServer is already called in Start() and should not normally be used.
// It is required for the tests, but should maybe be replaced with a more
// encapsulated solution.
func (a *App) GetServer() *handler.Server {
	srv := handler.NewDefaultServer(graph.NewSchema(a.dbClient))
	srv.Use(entgql.Transactioner{TxOpener: a.dbClient})
	return srv
}

// AddDummyData creates data in the DB
func (a *App) AddDummyData() error {
	if a.dbClient == nil {
		return errors.New("the database has not been initialized")
	}

	if err := a.addDummyData_Fabelmans(); err != nil {
		return err
	}
	if err := a.addDummyData_CPE(); err != nil {
		return err
	}
	if err := a.addDummyData_Contributors(); err != nil {
		return err
	}
	return nil
}

func (a *App) addDummyData_Fabelmans() error {
	ctx := context.Background()

	outcomeYes, err := a.dbClient.Outcome.Create().SetText("Yes").Save(ctx)
	if err != nil {
		return err
	}
	probabilityYes, err := a.dbClient.Probability.Create().SetValue(30).SetOutcome(outcomeYes).
		Save(ctx)
	if err != nil {
		return err
	}
	outcomeNo, err := a.dbClient.Outcome.Create().SetText("No").Save(ctx)
	if err != nil {
		return err
	}
	probabilityNo, err := a.dbClient.Probability.Create().SetValue(70).SetOutcome(outcomeNo).Save(
		ctx)
	if err != nil {
		return err
	}

	estimate, err := a.dbClient.Estimate.Create().
		SetReason("It's a great film and it's of the type that the Academy loves!").
		SetCreated(timeParseOrPanic(
			time.RFC3339,
			"2022-10-30T17:05:00+01:00",
		)).
		AddProbabilities(probabilityYes, probabilityNo).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = a.dbClient.Forecast.Create().
		SetTitle("Will \"The Fabelmans\" win \"Best Picture\" at the Oscars 2023?").
		SetCreated(timeParseOrPanic(
			time.RFC3339,
			"2022-10-30T17:05:00+01:00",
		)).
		SetResolves(timeParseOrPanic(
			time.RFC3339,
			"2023-03-11T23:59:00+01:00",
		)).
		AddEstimates(estimate).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) addDummyData_CPE() error {
	ctx := context.Background()

	outcomeGradeA, err := a.dbClient.Outcome.Create().
		SetText("C2 - Grade A").
		SetCorrect(true).
		Save(ctx)
	if err != nil {
		return err
	}
	probabilityGradeA, err := a.dbClient.Probability.Create().
		SetValue(40).
		SetOutcome(outcomeGradeA).
		Save(ctx)
	if err != nil {
		return err
	}
	outcomeGradeB, err := a.dbClient.Outcome.Create().SetText("C2 - Grade B").Save(ctx)
	if err != nil {
		return err
	}
	probabilityGradeB, err := a.dbClient.Probability.Create().SetValue(30).SetOutcome(
		outcomeGradeB).
		Save(ctx)
	if err != nil {
		return err
	}
	outcomeGradeC, err := a.dbClient.Outcome.Create().SetText("C2 - Grade C").Save(ctx)
	if err != nil {
		return err
	}
	probabilityGradeC, err := a.dbClient.Probability.Create().SetValue(20).SetOutcome(
		outcomeGradeC).
		Save(ctx)
	if err != nil {
		return err
	}
	outcomeC1, err := a.dbClient.Outcome.Create().SetText("C1").Save(ctx)
	if err != nil {
		return err
	}
	probabilityC1, err := a.dbClient.Probability.Create().SetValue(8).SetOutcome(
		outcomeC1).
		Save(ctx)
	if err != nil {
		return err
	}
	outcomeFail, err := a.dbClient.Outcome.Create().SetText("Fail").Save(ctx)
	if err != nil {
		return err
	}
	probabilityFail, err := a.dbClient.Probability.Create().SetValue(2).SetOutcome(
		outcomeFail).
		Save(ctx)
	if err != nil {
		return err
	}

	estimate, err := a.dbClient.Estimate.Create().
		SetReason("I'm well prepared and performed well on test"+
			" exams.").
		SetCreated(timeParseOrPanic(
			time.RFC3339,
			"2022-10-15T13:10:00+02:00",
		)).
		AddProbabilities(probabilityGradeA, probabilityGradeB, probabilityGradeC, probabilityC1,
			probabilityFail).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = a.dbClient.Forecast.Create().
		SetTitle("What grade will I get in my upcoming CPE exam?").
		SetDescription("CPE C2 exam. Grade C1 is the worst passing grade. " +
			"It's a language exam using the Common European Framework of" +
			" Reference for Languages.").
		SetCreated(timeParseOrPanic(
			time.RFC3339,
			"2022-10-30T17:05:00+01:00",
		)).
		SetResolves(timeParseOrPanic(
			time.RFC3339,
			"2023-03-11T23:59:00+01:00",
		)).
		SetCloses(timeParseOrPanic(
			time.RFC3339,
			"2022-11-11T23:59:00+01:00",
		)).
		SetResolution(forecast.ResolutionRESOLVED).
		AddEstimates(estimate).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) addDummyData_Contributors() error {
	ctx := context.Background()

	outcomeYes, err := a.dbClient.Outcome.Create().SetText("Yes").Save(ctx)
	if err != nil {
		return err
	}

	outcomeNo, err := a.dbClient.Outcome.Create().SetText("No").Save(ctx)
	if err != nil {
		return err
	}

	probability1Yes, err := a.dbClient.Probability.Create().SetValue(15).SetOutcome(outcomeYes).
		Save(ctx)
	if err != nil {
		return err
	}
	probability1No, err := a.dbClient.Probability.Create().SetValue(85).SetOutcome(outcomeNo).Save(
		ctx)
	if err != nil {
		return err
	}
	estimate1, err := a.dbClient.Estimate.Create().
		SetReason("It's a new project and people are usually busy.").
		SetCreated(timeParseOrPanic(
			time.RFC3339,
			"2022-10-01T11:00:00+02:00",
		)).
		AddProbabilities(probability1Yes, probability1No).
		Save(ctx)
	if err != nil {
		return err
	}

	probability2Yes, err := a.dbClient.Probability.Create().SetValue(1).SetOutcome(outcomeYes).
		Save(ctx)
	if err != nil {
		return err
	}
	probability2No, err := a.dbClient.Probability.Create().SetValue(99).SetOutcome(outcomeNo).Save(
		ctx)
	if err != nil {
		return err
	}
	estimate2, err := a.dbClient.Estimate.Create().
		SetReason("Despite multiple people expressing interest nobody "+
			"has contributed so far. The year is almost over.").
		SetCreated(timeParseOrPanic(
			time.RFC3339,
			"2022-10-01T11:00:00+02:00",
		)).
		AddProbabilities(probability2Yes, probability2No).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = a.dbClient.Forecast.Create().
		SetTitle("Will the number of contributors to \"Cleodora\" be more "+
			"than 3 at the end of 2022?").
		SetDescription("A contributor is any person who has made a commit"+
			" in any Git repository of the cleodora-forecasting GitHub"+
			" organization.").
		SetCreated(timeParseOrPanic(
			time.RFC3339,
			"2022-10-01T11:00:00+02:00",
		)).
		SetResolves(timeParseOrPanic(
			time.RFC3339,
			"2022-12-31T23:59:00+01:00",
		)).
		AddEstimates(estimate1, estimate2).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

type sqlite3Driver struct {
	*sqlite.Driver
}

type sqlite3DriverConn interface {
	Exec(string, []driver.Value) (driver.Result, error)
}

func (d sqlite3Driver) Open(name string) (conn driver.Conn, err error) {
	conn, err = d.Driver.Open(name)
	if err != nil {
		return
	}
	_, err = conn.(sqlite3DriverConn).Exec("PRAGMA foreign_keys = ON;", nil)
	if err != nil {
		_ = conn.Close()
	}
	return
}

func init() {
	// Register a custom driver to use Go-only SQLite
	sql.Register("sqlite3", sqlite3Driver{Driver: &sqlite.Driver{}})
}

// timeParseOrPanic parses the time and converts it to UTC
func timeParseOrPanic(layout string, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}

	return t.UTC()
}
