#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

./dist/cleoc_linux_amd64_v1/cleoc \
    --url https://demo.cleodora.org \
    add forecast \
    --title "Will Cleodora have 5 or more contributors by the end of 2023?" \
    --resolves 2024-01-01T00:00:00+01:00 \
    --description "A contributor is any person who makes any change to the" \
        "Git repository (which includes cleosrv, cleoc, website etc.)." \
    --probability Yes=70 \
    --probability No=30 \
    --reason "Just my gut feeling"
