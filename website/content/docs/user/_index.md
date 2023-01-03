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

## Installation

Release 0.1.0 is coming soon.


## Configuration

Release 0.1.0 is coming soon.


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
