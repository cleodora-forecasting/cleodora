---
title: "Ideas"
weight: 5
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---

# Ideas

A collection of thoughts and ideas that haven't yet made it into actual tasks
(and may never do so).

## cleoc (CLI)

GopherCon 2019: Carolyn Van Slyck - Design Command-Line Tools People Love

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
