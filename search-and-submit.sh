#!bin/bash
echo "run linter "$1" for keyword "$2
gh search repos $2 --limit 100 --language=go --json url --jq '.[]|.url' | xargs -I {} bash submit.sh $1 {}
