#!/bin/bash

echo "pulling"

git pull

git describe --tags

echo "what's the tag"

read tag

tag=$(echo "$tag" | sed 's/ //g')

export TAG=$tag

export go=../linux/go/bin/go # is the 19 version

chmod +x scripts/release.sh
chmod +x scripts/tag.sh


./scripts/release.sh
./scripts/tag.sh

git add .
git commit -m "final deps-$tag"

git push origin --tags