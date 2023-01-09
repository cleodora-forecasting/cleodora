# Ideas

A collection of thoughts and ideas that haven't yet made it into actual tasks
(and may never do so).

## cleoc (CLI)

### Design Command-Line Tools People Love

Title: GopherCon 2019: Carolyn Van Slyck - Design Command-Line Tools People Love

Video: https://www.youtube.com/watch?v=eMz0vni6PAw

Code: https://github.com/carolynvs/emote

Slides: https://carolynvanslyck.com/talk/go/cli/

* Use verb-noun
* Avoid positional arguments unless they are all the same type (e.g. delete a b c -> fine)
* Add a `--json` flag to every command
* Default to human output
* Dates -> go humanize
* Use aliases to compromise between brevity and discoverability (have 'aliases' section in help)
* combine multiple commands to achieve tasks. avoid piping except for scripts.

https://github.com/getporter/porter is a professional tool that applies many of
these ideas and can serve as inspiration.

cmd/porter/main.go

pkg/porter/porter.go


### Command Line Interface Guidelines

https://clig.dev/


## Mobile Client

### Rationale

* Most people carry a smartphone around with them
* It's ideal for making quick things right when you think of them e.g. create a
  new forecast that just popped into your head
* Having a desktop application you need to start when booting your PC is much
  more cumbersom


### Plan for 1.0

Right now the idea is to have a web application that people can run on their
PC, their home server or on some private server in the cloud.  This web
application can then be accessed with any browser, including one on your
smartphone.  The Web UI must be mobile friendly.


### Problems

* Starting the (single binary) web application on the PC and opening
  localhost:N is already somewhat challenging for many people. Running the
  application as a service, with a fixed IP etc. to be able to access it on the
  go will not be an option for many users.
* If the mobile device has no internet connection or is not in the same network
  as the web application server then it can't be used


### Possible Solutions

#### Batch Import

Implement an easy batch import function in the web UI e.g. to import a simple
CSV, for example:

```csv
Will I get a promotion in 3 months?;40%;2023-03-01
Will the price of gas drop to 1.50 € tomorrow?;5%;2022-11-14
```

Then it would be possible to create forecasts by writing into a text file on
the phone or any other device. When having access to the Web UI just copy paste
into it.


#### Simple Mobile 'Creator' App

A simple app that just allows recording new forecasts and creates them in the
Cleodora app server when a "Sync" button is pressed. Then the forecasts
disappear. It's only for creating forecasts on the go and pushing them to the
server when you get home and start the PC.


#### Full Blown Mobile App

Completely separate code base that has at least the most important features of
the Cleodora web application.  Probably some sync functionality needed to
synchronize phone with desktop to leverage the power of the web application
version.
