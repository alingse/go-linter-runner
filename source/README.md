## Update

```bash
cat keywords.txt|xargs -I {} sh search-github.sh {} >> all.github.json

cat all.github.json|jq -r -c '.[]' | jq -r -c 'select(.isFork|not)|select(.isArchived|not)'

## sort by star+issue+watcher
```


## Awesome

```
curl https://github.com/avelino/awesome-go/blob/main/README.md

cat awesome-go-readme | grep -o 'https://github.com/[^)]*' |sort|uniq > popular.txti
```

