---
title: "User Documentation"
weight: 1
# bookFlatSection: false
bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---

# User Documentation

**IMPORTANT** There is no release 0.1.0 yet. This documentation is in
preparation for the release which is coming very soon.


## cleosrv

`cleosrv` is the Cleodora server that includes all you need in a single binary.


### Installing cleosrv

Download the latest release for [Windows
(64bit)](https://github.com/cleodora-forecasting/cleodora/releases/download/v0.1.0/TODO),
[Linux (64bit)](), [Linux ARM (64bit)]() or [Mac OS (64bit)]().

You may also get the `cleodora` Docker container from Docker Hub.


### Running cleosrv

You can start `cleosrv` in a terminal/console window. On Windows you can also
double-click on it and such a console will open automatically. Now `cleosrv` is
running, congratulations!

Open a web browser such as Firefox, Google Chrome, Microsoft Edge or Safari and
type http://localhost:8080 into the address bar. That's it!

If you want to stop `cleosrv` you can close the console window or type Ctrl+C
inside that console.


### Configuring cleosrv

You do not need to do any of this to run `cleosrv`, just if you want to change
some default option.

You can pass parameters on the command line to specify the path to a different
database or a different config file. Execute `cleosrv --help`.

You can also create (or edit) the default config file named `cleosrv.yml`:

```yaml
address: localhost:9999
```

The default location for this file is:

* On Linux: `~/.config/`
* On Mac OS: `~/Library/Application Support/`
* On Windows: `%LOCALAPPDATA%` (just type that into an Explorer address bar)


### Back up the database

`cleosrv` stores all forecasting data in a SQLite database. This is a single
file on your computer. If you make a copy of this file, you have backed up
everything you need for running `cleosrv`.

You can see the location of this file by looking in the console window where
`cleosrv` is running. There should be one line starting with `Database: ` which
tells you the location. Stop `cleosrv` and make a copy of that file to
whichever location you want.


## cleoc

`cleoc` is a Cleodora CLI (command-line interface) client to interact with a
Cleodora server. You do not need to use it. In particular, if you have never
used a CLI before (or don't know what that is), you may not get much benefit
out of it.


### Installing cleoc

Download the latest release for [Windows
(64bit)](https://github.com/cleodora-forecasting/cleodora/releases/download/v0.1.0/TODO),
[Linux (64bit)](), [Linux ARM (64bit)]() or [Mac OS (64bit)]().


### Running cleoc

Execute it in a console window. Try `cleoc --help` to get started.


### Configuring cleoc

You can pass multiple CLI parameters to configure `cleoc`. Check the `--help`
output. You may also create a config file. The default location is displayed as
part of `--help`.

```yaml
url: http://localhost:9999
```


## Backwards Compatibility

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
