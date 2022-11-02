<img align="right" src="./design/logo_full.png">

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

# Dev Setup

* [Install Go](https://go.dev/doc/install) 1.18 or higher
* [Install npm](https://nodejs.org/en/download/) v16.18.0 or higher

```bash
git clone https://github.com/cleodora-forecasting/cleodora
cd cleodora
go get
cd frontend
npm install
```

## Gitpod

A simple dev environment, directly in the browser.

<a href="https://gitpod.io/#https://github.com/cleodora-forecasting/cleodora">
  <img
    src="https://img.shields.io/badge/Contribute%20with-Gitpod-908a85?logo=gitpod"
    alt="Contribute with Gitpod"
  />
</a> (Note: the frontend won't currently work with Gitpod because the backend
URL is hardcoded to localhost:8080).
<!-- https://github.com/gitpod-io/gitpod/issues/2466 -->

# Run

## Backend (GraphQL API)

```bash
go run .
```

* GraphQL playground: http://localhost:8080/playground/
* GraphQL API: http://localhost:8080/query

## Frontend (React app)

```bash
cd frontend
npm start
```

Open http://localhost:3000

The backend should also be running, otherwise you will get an error.

## Backend with statically embedded frontend

```bash
cd frontend
npm run build
cd ..
go run -tags frontend .
```

Open http://localhost:8080

# Build

```bash
make build
```

You can find the binary containing everything under `build/cleodora` . Run it
and access http://localhost:8080 in the browser.

# GraphQL playground (GraphiQL)

Start frontend and backend as described above.

Open http://localhost:8080/playground/ in a browser and create some forecasts:

```graphql
mutation createForecast {
    createForecast(
        input: {
            summary: "Will 'The Fabelmans' win 'Best Picture'?",
            description: "The new Steven Spielberg movie. Academy Award for Best Picture 2023.",
            resolves: "2023-03-01T12:00:00+01:00",
        }
    ) {
        id
        summary
        description
        created
        resolves
        closes
        resolution
    }
}
```

Open http://localhost:3000 in a browser and see the list of forecasts.

# Tests

## Setup

```bash
sudo apt install firefox-geckodriver
```

## Execution

```bash
cd e2e_tests
npm install
node_modules/.bin/mocha --timeout 15000 frontPageTest.spec.js
```
