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


## Unreleased

[Git history](https://github.com/cleodora-forecasting/cleodora/compare/v0.3.0...HEAD)


## v0.3.0

Release date: 2023-05-25


### Added

* API, Frontend: Estimates from forecasts that have been resolved now include
  the Brier score.
  ([#298](https://github.com/cleodora-forecasting/cleodora/issues/298),
  [#499](https://github.com/cleodora-forecasting/cleodora/issues/499))
* The `cleosrv` binary now also supports a `--version/-v` flag in addition to
  the identical `version` subcommand because it's more convenient.

[Git history](https://github.com/cleodora-forecasting/cleodora/compare/v0.2.0...v0.3.0)


## v0.2.0

Release date: 2023-03-04


### Added

* API, Frontend: Resolve forecasts (i.e. decide which outcome was actually the
  correct one).
  ([#154](https://github.com/cleodora-forecasting/cleodora/issues/154),
  [#155](https://github.com/cleodora-forecasting/cleodora/issues/155))
* Frontend: Forecast details dialog that can be opened by clicking on the
  forecast title.
  ([#296](https://github.com/cleodora-forecasting/cleodora/issues/296))


### Fixed

* Frontend: Display more detailed error messages.
  ([#173](https://github.com/cleodora-forecasting/cleodora/issues/173),
  [#212](https://github.com/cleodora-forecasting/cleodora/issues/212))
* API: Ensure the forecast date values are logical i.e. _created_ must be
  before _closes_ which in turn must be before _resolves_. Existing data is
  automatically migrated to enforce this constraint.
  ([#234](https://github.com/cleodora-forecasting/cleodora/issues/234),
  [#264](https://github.com/cleodora-forecasting/cleodora/issues/264))
* Frontend: When adding a new outcome in the _Add Forecast_ form the focus
  moves to that new outcome so you can start typing directly.
  ([#161](https://github.com/cleodora-forecasting/cleodora/issues/161))


### Changed

* API: Date values passed to the GraphQL API may contain any time zone, but
  they will always be converted to UTC in the backend and returned as such. In
  the web frontend everything is converted to the time zone of the browser
  before displaying it. Existing data is automatically transformed to use UTC.
  ([#265](https://github.com/cleodora-forecasting/cleodora/issues/265),
  [#264](https://github.com/cleodora-forecasting/cleodora/issues/264))
* Frontend: Simplify the _Add Forecast_ form on the main page.
  ([#233](https://github.com/cleodora-forecasting/cleodora/issues/233))
* Frontend: Prettify and simplify the forecast list.
* Frontend: In _Add Forecast_ form show two outcome rows by default because
  making a forecast with a single outcome makes little or no sense.


[Git history](https://github.com/cleodora-forecasting/cleodora/compare/v0.1.1...v0.2.0)


## v0.1.1

Release date: 2023-01-11


### Added

* Doc: Improved user and developer documentation


### Fixed

* Internal: Build issues when preparing a release


[Git history](https://github.com/cleodora-forecasting/cleodora/compare/v0.1.0...v0.1.1)


## v0.1.0

Release date: 2023-01-10


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
