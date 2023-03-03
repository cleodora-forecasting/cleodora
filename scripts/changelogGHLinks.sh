#!/bin/bash
set -o errexit
set -o nonunset
set -o pipefail

# Look for everything of the form #123 that is not surrounded with [] (i.e. not
# [#123]) and replace it with a link to GitHub.
#
# So #123 would be converted into
# [#123](https://github.com/cleodora-forecasting/cleodora/issues/123)

sed -i "s;\([^[]\)#\([[:digit:]]\+\)\([^]]\);\1[#\2](https://github.com/cleodora-forecasting/cleodora/issues/\2)\3;g" website/content/docs/changelog.md
