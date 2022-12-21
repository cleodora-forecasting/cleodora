#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

# Script to merge all dependabot PRs (pull requests) into 'main'

if [[ "`git rev-parse --abbrev-ref HEAD`" != "main" ]]
then
    echo "Not on main! Exiting"
    exit 1
fi

git diff --quiet --exit-code || (echo "There are uncomitted changes! Exiting" && exit 1)

git fetch
git remote prune origin

ALL_PRS=`git for-each-ref --format='%(refname)' refs/remotes/origin/dependabot/`

for PR in ${ALL_PRS}
do
    echo "PR: ${PR}"
    git merge-base --is-ancestor "${PR}" HEAD && printf "Already merged\n\n" && continue || /bin/true
    git merge "${PR}" -m "Merge dependabot update"
    echo ""
done

echo "All PRs merged"

./installDependencies.sh
make lint
make generate

git diff --exit-code || (echo "Code was changed via lint/generate" && exit 1)

./runE2ETests.sh

echo "Successfully done. You must run 'git push' to publish the changes."
