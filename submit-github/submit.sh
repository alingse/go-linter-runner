#!bin/bash
set -e
echo "will run linter yaml "$1" for repo:"$2
ls $1
gh workflow run .github/workflows/run-on-github.yml -F yaml_file=$1 -F repo_url=$2
