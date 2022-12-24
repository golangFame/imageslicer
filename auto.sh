#!/bin/bash

echo "pulling"

git pull

git describe --tags

echo "what's the tag"

read tag

tag=$(echo "$tag" | sed 's/ //g')

export TAG=$tag


chmod +x scripts/release.sh
chmod +x scripts/tag.sh

./scripts/release.sh
./scripts/tag.sh

git push origin --tags