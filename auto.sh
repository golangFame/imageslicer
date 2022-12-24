echo "what's the tag"

read tag

set TAG=$tag

git pull

sh scripts/release.sh
sh scripts/tag.sh

git push origin --tags