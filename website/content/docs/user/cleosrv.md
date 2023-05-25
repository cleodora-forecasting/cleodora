---
title: "cleosrv"
weight: 1
# bookFlatSection: false
bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---

# cleosrv

**Latest release: 0.3.0**

`cleosrv` (Cleodora server) is the main application. It includes a beautiful
web UI (user interface) that you can access from your browser as well as a
GraphQL API (which you can ignore if you don't know what that is).

All `cleosrv` data is stored in a single file, a SQLite database, that will be
automatically created it if does not exist.

![cleosrv web frontend](/cleosrv_frontend.png "cleosrv web frontend")


## Installing cleosrv

Download the latest release for
[Windows (64bit)](https://github.com/cleodora-forecasting/cleodora/releases/download/v0.3.0/cleosrv_0.3.0_Windows_64bit.zip),
[Linux (64bit)](https://github.com/cleodora-forecasting/cleodora/releases/download/v0.3.0/cleosrv_0.3.0_Linux_64bit.tar.gz),
[Linux ARM (64bit)](https://github.com/cleodora-forecasting/cleodora/releases/download/v0.3.0/cleosrv_0.3.0_Linux_ARM64.tar.gz) or
[MacOS (64bit)](https://github.com/cleodora-forecasting/cleodora/releases/download/v0.3.0/cleosrv_0.3.0_macOS_64bit.tar.gz).

Unpack the archive wherever you want it. There will only be a single executable
inside named `cleosrv` or `cleosrv.exe`.


## Updating cleosrv

Download the latest release (see above) and replace the existing `cleosrv` or
`cleosrv.exe` executable on your PC. All data in your database will be
transformed automatically to be compatible with the new version when you start
`cleosrv`.

Before updating it's a good idea to back up the database. See below for
instructions.

After updating to a newer version do not open the same database with the old
`cleosrv` version.


## Running cleosrv

You can start `cleosrv` in a terminal/console window. On Windows you can also
double-click on the `cleosrv.exe` file you downloaded and such a console will
open automatically. Now `cleosrv` is running, congratulations! You may be asked by
your firewall whether you want to allow this application to run. Answer "Yes".

Open a web browser such as Firefox, Google Chrome, Microsoft Edge or Safari and
type `http://localhost:8080` into the address bar. That's it!

If you want to stop `cleosrv` you can close the console window or hit Ctrl+C
(possibly twice) inside that console.


## Configuring cleosrv

You do not need to do any of this to run `cleosrv`, just if you want to change
some default option.

You can pass parameters on the command line to specify the path to a different
database or a different config file. Execute `cleosrv --help`.

You can also create (or edit) the default config file:

**cleosrv.yml**

```yaml
# Uncomment (remove the leading # from) any parameter you want to use

# Default address the server should listen on. Use 0.0.0.0:8080 if you want to
# expose it in your local network. Don't include 'http://'.
#address: localhost:8080

# Path to the database if you don't want to use the default
#database: /path/to/some/cleosrv.db

#frontend:
#  # A text that is displayed in the footer of the web frontend
#  footer_text: ""
```

The default location for this file is:

* On Linux: `~/.config/`
* On Mac OS: `~/Library/Application Support/`
* On Windows: `%LOCALAPPDATA%` (just type that into an Explorer address bar)


## Backing up the database

`cleosrv` stores all forecasting data in a SQLite database. This is a single
file on your computer. If you make a copy of this file, you have backed up
everything you need for running `cleosrv`.

You can see the location of this file by looking in the console window where
`cleosrv` is running. There should be one line starting with `Database: ` which
tells you the location. Stop `cleosrv` and make a copy of that file to
whichever location you want.

![cleosrv console](/cleosrv_console.png "cleosrv console")


## Advanced

**Note:** The below topics may be advanced and you can safely ignore them.


### Docker

You may also get the `cleodora` Docker image from Docker Hub:

`docker pull cleodora/cleodora`

The image only includes `cleosrv` (not `cleoc`).

When running on Docker you should specify a named volume (`cleodora_data` here),
otherwise your data will be stored on an anonymous volume and you may lose it:

```bash
docker run -p 8080:8080 -v cleodora_data:/data cleodora/cleodora:latest
```


### GraphQL API

The endpoint is `BASE_URL/query` e.g. http://localhost:8080/query .

You can access a GraphiQL (graphical interactive in-browser GraphQL IDE)
instance under `BASE_URL/playground/` e.g. http://localhost:8080/playground/ .


### Backwards Compatibility

Any **data** you store in Cleodora will always be **migrated automatically**
when you update to newer official releases. For example if you start with
0.1.0, then update to 0.2.0, then update to 0.7.0 (skipping a few versions) you
will always be able to see and update existing forecasts as well as add new
ones.

Of course, just like any other piece of software, Cleodora might contain bugs
which can lead to something not working as expected, also during an update.
Backing up your database before an update is highly recommended.

If you want to experiment with development versions between official releases
you should use a separate database for that (`--database` parameter). You can
also make a copy of your real database file and experiment with that.

The command line parameters, GraphQL API format and config file format will not
necessarily be backwards compatible before release 1.0.0 . Changes will be
documented in the changelog. The reason for this is that we don't want to drag
potential bad early design decisions along forever.
