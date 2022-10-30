<img align="right" src="./logo_full.png">

# Cleodora Forecasting

Software to track personal forecasts/predictions and systematically improve at
making them.

Examples of such forecasts:

* Will "The Fabelmans" win "Best Picture" at the Oscars 2023?
* Will I get an A in my upcoming exam?
* Will there be nice weather on my birthday?
* Will the number of contributors for "Cleodora" be more than 3 at the end of 2022?

All information, including development guidelines, requirements and roadmap can
be found on the website https://cleodora.org .

# Run

Without embedded frontend:

```bash
go run .
```

With statically embedded frontend:

```bash
cd frontend
npm run build
cd ..
go run -tags frontend .
```

# Build

```bash
make build
```

# Playing around

One console:
```bash
go run .
```

Another console:
```bash
cd frontend
npm start
```

Open http://localhost:8080/playground/ in a browser and create some forecasts:

```graphql
mutation createForecast {
  createForecast(
    input: {
        summary: "Will Fabelmans win?",
        description: "Oscars 2023",
        closes: "2023-03-01T12:00:00+01:00",
    }
  ) {
    id
    summary
    description
    created
    closes
  }
}
```

Open http://localhost:3000 in a browser and see the list of forecasts.
