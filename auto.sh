#!/bin/bash

echo "pulling"

git pull

git describe --tags

echo "what's the tag"

read tag

set TAG=$tag

sh scripts/release.sh
sh scripts/tag.sh

git push origin --tags