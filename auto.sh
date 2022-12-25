#!/bin/bash

echo "pulling"

git pull

go work sync

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

git push
git push origin --tags


go list -m -json github.com/goferHiro/image-slicer/examples/basic@latest
go list -m -json github.com/goferHiro/image-slicer/examples/websocket-server@latest
go list -m -json github.com/goferHiro/image-slicer@latest
