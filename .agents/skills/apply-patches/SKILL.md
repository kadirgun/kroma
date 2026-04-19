---
name: apply-patches
description: >-
  Automate the patch push/pull workflow for Chromium repository modifications.
  Applies patches to the codebase, detects and resolves rejected hunks,
  and manages patch diff extraction. Use when pushing patches fails with
  rejection errors or when systematically extracting changes into patches.
user-invocable: true
---

# Apply Patches Workflow

## When to use

Use this skill when a Chromium patch cannot be applied cleanly with **`patcher push`** and a `.rej` file is generated.

## Workflow

### 1. Run `patcher push`

```bash
patcher push
```

If the command succeeds, no further action is needed.

### 2. Inspect the `.rej` file

If patcher push returns "Rejected hunk" error:

1. Read the `.rej` file to understand what changes failed
2. Create a todo list with each change that needs to be made to the target file
3. Apply each change to the target file using the edit tool

### 3. Edit the target file manually

- Apply the same changes from the rejected patch to the target file
- Use the editor and patch context to preserve nearby code
- Do not use `git checkout` to revert the file

### 4. Remove the `.rej` file

After the manual edits are complete:

```bash
Remove-Item "path/to/target.rej" -Force
```

### 5. Pull the patched file

Pull the repaired patch file back into the patches directory, which will create a new patch file with the applied changes:

```bash
patcher pull path/to/target
```

## Rules to follow

- Never use `git checkout` to reset files during patch resolution
- Edit the target file manually based on `.rej` contents
- Delete the `.rej` file only after edits are finished
- Keep a todo list for each change if multiple hunks failed
- Note that patch file paths are relative to `chromium/src/`

## Example

```bash
# Attempt to push the patch
patcher push

# If a reject occurs, inspect the .rej file and fix the target file
# Remove the reject file once edits are done
Remove-Item "chromium/src/chrome/browser/BUILD.gn.rej" -Force

# Then pull the repaired patch file back in
patcher pull chrome/browser/BUILD.gn
```
