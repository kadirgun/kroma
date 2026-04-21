# Git Shallow Clone - Prevent .git Bloat

```bash
cd chromium/src
git config fetch.depth 1
git config fetch.recursesubmodules off
git config fetch.prune true
git config push.recursesubmodules off
git config gc.auto 0
git config gc.autopacklimit 0
git config gc.autodetach false
```

## If .git Bloats

```bash
cd chromium/src
git gc --prune=now
```
