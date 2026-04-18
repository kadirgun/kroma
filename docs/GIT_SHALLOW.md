# Git Shallow Clone - Prevent .git Bloat

```bash
cd chromium/src
git config --local fetch.depth 1
git config --local fetch.recursesubmodules off
git config --local fetch.prune true
git config --local push.recursesubmodules off
```

## If .git Bloats

```bash
cd chromium/src
git gc --prune=now
```
