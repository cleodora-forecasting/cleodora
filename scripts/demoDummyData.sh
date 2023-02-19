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

./dist/cleoc_linux_amd64_v1/cleoc \
    --url https://demo.cleodora.org \
    add forecast \
    --title "Will I go at least 10 times to the gym this month?" \
    --resolves `date -Iseconds -d "+30 days"` \
    --probability Yes=65 \
    --probability No=35 \
    --reason "Past experience"

./dist/cleoc_linux_amd64_v1/cleoc \
    --url https://demo.cleodora.org \
    add forecast \
    --title "Will I finish the financial report this week?" \
    --resolves `date -Iseconds -d "+7 days"` \
    --probability Yes=15 \
    --probability No=85 \
    --reason "I haven't even started"

./dist/cleoc_linux_amd64_v1/cleoc \
    --url https://demo.cleodora.org \
    add forecast \
    --title "How many cookies will I eat in the evening?" \
    --resolves `date -Iseconds -d "+1 day"` \
    --probability None=5 \
    --probability "Less than 5=20" \
    --probability "5 to 15=60" \
    --probability "15 or more=15" \
    --reason "I know myself"
