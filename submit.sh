#!bin/bash
echo "will run linter "$1" for repo:"$2
gh workflow run .github/workflows/check-$1.yml -F repo_url=$2
