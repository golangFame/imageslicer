#! /bin/bash

#git checkout main
#
#echo "pulling"
#
#git pull

git describe --tags

echo "what's the tag"

read tag

tag=$(echo "$tag" | sed 's/ //g')

export TAG=$tag

chmod +x scripts/release.sh
chmod +x scripts/tag.sh

./scripts/release.sh
./scripts/tag.sh

go work sync

git add .
git commit -m "final deps-$tag"

git push
git push origin --tags

echo "check the latest version, sleeping for 1 min for the go registry to sync up"

sleep 1

go list -m -json github.com/golangFame/imageslicer/examples/basic@latest
go list -m -json github.com/golangFame/imageslicer/examples/websocket-server@latest
go list -m -json github.com/golangFame/imageslicer@latest
