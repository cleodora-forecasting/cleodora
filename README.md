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

## Manual setup

Note that instead of setting all this up on your computer, you can use
**Gitpod** or **GitHub Codespaces** for pre-configured dev environments
directly in your browser. See below.

* [Install Go](https://go.dev/doc/install) 1.18 or higher
* [Install npm](https://nodejs.org/en/download/) v16.18.0 or higher
* (optional) [Install hugo](https://gohugo.io/installation/) (extended flavor)
  v0.97.3 or higher to make changes to the cleodora.org website

```bash
git clone https://github.com/cleodora-forecasting/cleodora
cd cleodora
./installDependencies.sh
```

## Gitpod

A simple dev environment, directly in the browser.

<a href="https://gitpod.io/#https://github.com/cleodora-forecasting/cleodora">
  <img
    src="https://img.shields.io/badge/Contribute%20with-Gitpod-908a85?logo=gitpod"
    alt="Contribute with Gitpod"
  />
</a>


## GitHub Codespaces

**IMPORTANT:** Codespaces ports are always _private_ by default. After starting
the API server you need to set the API port (8080) to public if you want to
access it from the frontend. For example with this command `gh codespace ports
visibility 8080:public -c $CODESPACE_NAME` .

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://github.com/codespaces/new?hide_repo_select=true&ref=main&repo=548549126&machine=basicLinux32gb)


# Run

## Backend (GraphQL API)

```bash
# Top level directory
go run .
```

* GraphQL playground: http://localhost:8080/playground/
* GraphQL API: http://localhost:8080/query

## Frontend (React app)

```bash
cd frontend
npm start
## Optionally, you can override the backend URL
# REACT_APP_API_URL = http://localhost:5555 npm start
```

Open http://localhost:3000

The backend should also be running, otherwise you will get an error.

## Backend with statically embedded frontend

```bash
cd frontend
npm run build
cd ..
go run -tags production .
```

Open http://localhost:8080

# Build

```bash
make build
```

You can find the binary containing everything under `build/cleosrv` . Run it
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

Start the app

```bash
make build
./build/cleosrv
```

Run the tests

```bash
cd e2e_tests
npm install
node_modules/.bin/mocha --timeout 15000 frontPageTest.spec.js
```
