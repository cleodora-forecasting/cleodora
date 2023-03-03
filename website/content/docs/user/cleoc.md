---
title: "cleoc"
weight: 2
# bookFlatSection: false
bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---

# cleoc

**Latest release: 0.1.1**

`cleoc` (Cleodora client) is a Cleodora CLI (command-line interface) client to
interact with a `cleosrv` Cleodora server.

You do not need to use `cleoc`. In particular, if you have never used a CLI
before (or don't know what that is), you may not get much benefit out of it.

To use `cleoc` you always need `cleosrv` to run as well, either on your local
machine or somewhere else.

![cleoc console](/cleoc_console.png "cleoc console")


## Installing cleoc

Download the latest release for
[Windows (64bit)](https://github.com/cleodora-forecasting/cleodora/releases/download/v0.1.1/cleoc_0.1.1_windows_amd64.zip),
[Linux (64bit)](https://github.com/cleodora-forecasting/cleodora/releases/download/v0.1.1/cleoc_0.1.1_linux_amd64.tar.gz),
[Linux ARM (64bit)](https://github.com/cleodora-forecasting/cleodora/releases/download/v0.1.1/cleoc_0.1.1_linux_arm64.tar.gz) or
[Mac OS (64bit)](https://github.com/cleodora-forecasting/cleodora/releases/download/v0.1.1/cleoc_0.1.1_darwin_amd64.tar.gz).

Unpack the archive wherever you want it. There will only be a single executable
inside named `cleoc` or `cleoc.exe`.


## Updating cleoc

Download the latest release (see above) and replace the existing `cleoc` or
`cleoc.exe` executable on your PC.


## Running cleoc

Execute it in a console window. Try `cleoc --help` to get started.


## Configuring cleoc

You can pass multiple CLI parameters to configure `cleoc`. Check the `--help`
output. You may also create a config file. The default location is displayed as
part of `--help`.

**cleoc.yml**

```yaml
# Uncomment (remove the leading # from) any parameter you want to use

# URL of the cleosrv
#url: http://localhost:8080
```
