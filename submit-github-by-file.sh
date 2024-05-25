#!bin/bash
echo "submit github actions "$1" for file: "$2
cat $2|xargs -I {} bash submit-github.sh $1 {}
