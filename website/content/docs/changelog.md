---
title: "Changelog"
weight: 9
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

The numbers in parentheses are GitHub issues that you can find
[here](https://github.com/cleodora-forecasting/cleodora/issues). They are
useful if you want to learn some more details about a specific changelog entry.


## Unreleased

### Added

* API, Frontend: Resolve forecasts (i.e. decide which outcome was actually the
  correct one). (#154, #155)
* Frontend: Forecast details dialog that can be opened by clicking on the
  forecast title. (#296)


### Fixed

* Frontend: Display more detailed error messages. (#173, #212)
* API: Ensure the forecast date values are logical i.e. _created_ must be
  before _closes_ which in turn must be before _resolves_. Existing data is
  automatically migrated to enforce this constraint. (#234, #264)
* Frontend: Simplify the _Add Forecast_ form on the main page. (#233)
* Frontend: When adding a new outcome in the _Add Forecast_ form the focus
  moves to that new outcome so you can start typing directly. (#161)


### Changed

* API: Date values passed to the GraphQL API may contain any time zone, but
  they will always be converted to UTC in the backend and returned as such. In
  the web frontend everything is converted to the time zone of the browser
  before displaying it. Existing data is automatically transformed to use UTC.
  (#265, #264)
* Frontend: Prettify and simplify the forecast list.


[Git history](https://github.com/cleodora-forecasting/cleodora/compare/v0.1.1...HEAD)


## v0.1.1

### Added

* Doc: Improved user and developer documentation


### Fixed

* Internal: Build issues when preparing a release


[Git history](https://github.com/cleodora-forecasting/cleodora/compare/v0.1.0...v0.1.1)


## v0.1.0

### Added

- `cleosrv` application that runs a local webserver on Windows, Linux or MacOS
- Docker image `cleodora` to run `cleosrv`
- `cleoc` CLI (command line interface) application to interact with a `cleosrv`
  server
- GraphQL API that allows listing forecasts and associated outcomes and
  probabilities
- List forecasts, outcomes and probabilities in the web frontend included in
  `cleosrv`
- Create new forecasts with associated outcomes and probabilities via `cleoc`
- Create new forecasts with associated outcomes and probabilities via the web
  frontend included in `cleosrv`


[Git history](https://github.com/cleodora-forecasting/cleodora/commits/v0.1.0)
