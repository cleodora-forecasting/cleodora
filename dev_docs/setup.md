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
./scripts/installDependencies.sh
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

A simple dev environment, directly in the browser.

**IMPORTANT:** Codespaces ports are always _private_ by default
([source](https://github.com/community/community/discussions/4068)). After
starting the codespace you need to set the API port (8080) to public if you
want to access it from the frontend
([technical background](https://github.com/community/community/discussions/4068)).
For example by clicking on the `PORTS`
tab, right click on `API (8080)` and then `Port Visibility` / `Public`.

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://github.com/codespaces/new?hide_repo_select=true&ref=main&repo=548549126&machine=basicLinux32gb)


# Run

## Backend (GraphQL API written in Go)

```bash
cd cleosrv
go run .
```

* GraphQL playground: http://localhost:8080/playground/
* GraphQL API: http://localhost:8080/query


## GUI / Frontend (React app)

Frequently the frontend is referred to as GUI (Graphical User Interface) in
this app because it's shorter.

```bash
cd frontend
npm start
## Optionally, you can override the backend URL
# REACT_APP_API_URL=http://localhost:5555 npm start
```

Open http://localhost:3000

The backend must also be running, otherwise you will get an error. This is
because the frontend queries the backend for data.


## Client (CLI written in Go)

The purpose is to interact with the API.

```bash
cd cleoc/cmd/
go run .
```


# Build

```bash
make build
```

You can find the binary containing frontend and backend under
`dist/cleosrv_*/cleosrv` . Run it and access http://localhost:8080 in the
browser. The frontend is embedded as static files inside the binary, so this
binary contains the entire Cleodora backend and frontend.

The client binary is `dist/cleoc_*/cleoc` .


# GraphQL playground (GraphiQL)

Start the backend as described above.

Open http://localhost:8080/playground/ in a browser.

You can list all the forecasts:

```graphql
query GetForecasts {
  forecasts {
    id
    title
    description
    created
    closes
    resolves
    resolution
  }
}
```

And create a forecast ...

```graphql
mutation createForecast($forecast: NewForecast!, $estimate: NewEstimate!) {
  createForecast(forecast: $forecast, estimate: $estimate) {
    id
    title
    estimates {
      id
      created
      reason
      probabilities {
        id
        value
        outcome {
          id
          text
          correct
        }
      }
    }
  }
}
```

... with variables:

```json
{
  "forecast": {
    "title": "Will it rain tomorrow?",
    "description": "It counts as rain if between 9am and 9pm there are 30 min or more of uninterrupted precipitation.",
    "resolves": "2022-01-31T10:00:00+01:00"
  },
  "estimate": {
    "reason": "My weather app says it will rain.",
    "probabilities": [
      {
        "value": 70,
        "outcome": {
          "text": "Yes"
        }
      },
      {
        "value": 30,
        "outcome": {
          "text": "No"
        }
      }
    ]
  }
}
```

If you start the frontend as described above you can see the result there as
well.

See [schema.graphql](./schema.graphql) to see the GraphQL schema, queries
and mutations.


# Tests

## Integration tests

```bash
go test ./...
```


## E2E (end to end) tests

Selenium end to end tests, also including the cleoc client.

```bash
sudo apt install firefox-geckodriver
```

```bash
./scripts/runE2ETests.sh
```
