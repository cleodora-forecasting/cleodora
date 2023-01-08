<img align="right" src="./design/logo_full.png">

# Cleodora Forecasting

![Website Build](https://github.com/cleodora-forecasting/cleodora/actions/workflows/website.yml/badge.svg)
![golangci-lint](https://github.com/cleodora-forecasting/cleodora/actions/workflows/golangci-lint.yml/badge.svg)
[![Go tests](https://github.com/cleodora-forecasting/cleodora/actions/workflows/go-tests.yml/badge.svg)](https://github.com/cleodora-forecasting/cleodora/actions/workflows/go-tests.yml)

[Demo](https://demo.cleodora.org/) (courtesy of fly.io) - will be reset at
regular intervals.

Software to track personal forecasts/predictions and systematically improve at
making them.

Examples of such forecasts:

* Will "The Fabelmans" win "Best Picture" at the Oscars 2023?
* Will I get an A in my upcoming exam?
* Will there be nice weather on my birthday?
* Will the number of contributors for "Cleodora" be more than 3 at the end of 2022?

Below an example mockup for how a user kept updating their "Will aliens be found on
the moon in 2022?" forecast due to new information. Analyzing such a forecasting
history is what will help them improve. You can find more such mockups on the
website https://cleodora.org/docs/development/roadmap/mockups/ .

![Mockup: History of a forecast](website/content/docs/development/roadmap/mockups/mockups_cleodora_history.jpg)


# Features

Concrete next steps can be found as GitHub issues.

Less detailed possible future features can be found on the
[website](https://cleodora.org/docs/development/roadmap/) (or also in here under
[./website/content/docs/development/roadmap](./website/content/docs/development/roadmap/)).


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
./runE2ETests.sh
```


# Deploy

## Docker

```bash
make clean
make build
DOCKER_TAG=0.1.0.dev.`git rev-parse --short HEAD`
echo "DOCKER_TAG: ${DOCKER_TAG}"
docker build --tag cleodora:${DOCKER_TAG} .
docker run -p 8080:8080 -v cleodora_data:/data cleodora:${DOCKER_TAG}
docker push cleodora/cleodora:${DOCKER_TAG}
```

Always start the container with a named volume (and keep using the same name,
`-v cleodora_data:/data` in the example below, even when updating the image):

```bash
docker run -p 8080:8080 -v cleodora_data:/data cleodora:VERSION
```

The reason is that this image will use an anonymous volume `/data` by default
to store the data. This means if you just stop a container and start a new one,
you will lose your data (e.g. list of forecasts). There are some other things
you can do to avoid this (but the best and easiest is using a named volume as
mentioned above):

* Use a bind mount.
* Before deleting the old container, start the new one with `--volumes-from`
  option to use the same (anonymous) volume. Then you can delete the old
  container.
* Disaster recovery: Find the anonymous volume and copy the data into a new
  volume. This will only work if the volume hasn't been deleted (luckily
  `docker rm` does not delete such volumes by default). See for example [this
  link](https://github.com/moby/moby/issues/31154#issuecomment-360531460).


## fly.io (demo.cleodora.org)

```bash
make clean
make build
flyctl deploy --local-only # use local Docker to build
```
