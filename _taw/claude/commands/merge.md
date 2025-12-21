---
description: Merge worktree branch into project's current branch
allowed-tools: Bash(git:*), Bash(echo:*), Bash(cat:*)
---

# Merge Branch

## Step 1: Check Current State

Run these commands to understand the current state:

```bash
echo "TASK_NAME=$TASK_NAME"
echo "PROJECT_DIR=$PROJECT_DIR"
echo "WORKTREE_DIR=$WORKTREE_DIR"
```

```bash
echo "Worktree branch: $(git branch --show-current)"
echo "Project branch: $(git -C $PROJECT_DIR branch --show-current)"
```

```bash
git status -s
```

## Step 2: Pre-merge Checks

### Check for uncommitted changes

If there are uncommitted changes in worktree, ask user whether to:
- Commit them first
- Stash them
- Abort merge

### Confirm merge target

The merge will be:
```
{worktree branch} -> {project branch}
```

Ask user to confirm before proceeding.

## Step 3: Perform Merge

```bash
git -C $PROJECT_DIR merge $TASK_NAME --no-ff -m "Merge branch '$TASK_NAME'"
```

Use `--no-ff` to preserve branch history.

## Step 4: Handle Conflicts (if any)

If merge conflicts occur:

```bash
git -C $PROJECT_DIR diff --name-only --diff-filter=U
```

- List conflicted files
- Inform user about conflicts
- Do NOT auto-resolve - let user decide

## Step 5: Show Result

```bash
git -C $PROJECT_DIR log --oneline -5
```

## After Merge

Inform user they can now run `/done` to clean up the task (worktree, branch, window).

Proceed with merge.
