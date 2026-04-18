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

## When to activate

Activate this skill when the user:

- Runs **`patcher push`** and encounters **"Rejected hunk"** errors
- Needs to **resolve patch conflicts** in `.rej` files
- Wants to **extract changes** from modified files into `.diff` patches
- Needs to **systematically apply** all patches and handle failures

## Workflow

### Step 1: Push patches to apply them

```bash
cd /path/to/kroma
patcher push
```

If all patches apply successfully, you're done.

### Step 2: Detect and inspect rejected hunks

If you see **"Rejected hunk"** error with a filename (e.g., `chrome/browser/chrome_content_browser_client.cc`):

1. Look for the corresponding `.rej` file in the same directory:

   ```bash
   # Example: if push failed on chrome/browser/chrome_content_browser_client.cc
   # the .rej file will be at the same path with .rej extension
   chromium/src/chrome/browser/chrome_content_browser_client.cc.rej
   ```

2. Read the `.rej` file to understand what changes failed:
   - The `.rej` file shows the exact hunks that patch couldn't apply
   - Each hunk shows the expected context and the attempted change
   - Context mismatch is the most common cause

### Step 3: Resolve the conflict

1. Open the source file (`chromium/src/chrome/browser/chrome_content_browser_client.cc`)
2. Find the location where the hunk should go (context in `.rej` file is your guide)
3. Manually apply the changes from the `.rej` file
4. Save the file

### Step 4: Pull the patch again to mark it as applied

After resolving conflicts:

```bash
patcher pull chrome/browser/chrome_content_browser_client.cc
```

This generates or updates `.diff` files in the `patches/` directory.

Ensure:

- You're running commands from the `kroma/` root directory
- The `chromium/src/` submodule exists and is checked out
- File paths in patches are relative to `chromium/src/`

## Examples

**Example**

```bash
# Push fails with "Rejected hunk" for chrome_content_browser_client.cc
patcher push

# Find the .rej file
cat chromium/src/chrome/browser/chrome_content_browser_client.cc.rej

# Edit the target file to match context, then delete the .rej file
rm chromium/src/chrome/browser/chrome_content_browser_client.cc.rej

# Pull the patch to update .diff files
patcher pull chrome/browser/chrome_content_browser_client.cc
```
