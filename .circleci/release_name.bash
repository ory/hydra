#!/bin/bash

release=$(curl -s "https://api.github.com/repos/$CIRCLE_PROJECT_USERNAME/$CIRCLE_PROJECT_REPONAME/releases")
tag=$(echo ${release} | jq -r ".[0].tag_name")
tag_name=$(echo ${release} | jq -r ".[0].name")

if [[ -n "$tag_name" ]]; then
    echo "export RELEASE_NAME=$tag_name" >> $BASH_ENV
elif [[ -n "$tag" ]]; then
    echo "export RELEASE_NAME=$tag" >> $BASH_ENV
else
    echo "export RELEASE_NAME=$CIRCLE_SHA1" >> $BASH_ENV
fi
