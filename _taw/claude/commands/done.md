---
description: Clean up completed task (worktree, agent dir, window)
allowed-tools: Bash(git:*), Bash(rm:*), Bash(tmux:*), Bash(ls:*), Bash(cat:*), Bash(gh:*), Bash(echo:*), Bash(printenv:*), Bash(test:*), Bash(true:*)
---

# Task Cleanup

## Step 1: Get Environment

First, get the environment variables:
```bash
echo "TASK_NAME=$TASK_NAME"
echo "TAW_DIR=$TAW_DIR"
echo "PROJECT_DIR=$PROJECT_DIR"
echo "WORKTREE_DIR=$WORKTREE_DIR"
```

## Step 2: Check Merge Status

### Method A: Check saved PR number

```bash
cat $TAW_DIR/agents/$TASK_NAME/.pr
```

If PR number exists, check its merge status:
```bash
gh pr view <PR_NUMBER> --json merged,state -q '"\(.state) merged:\(.merged)"'
```

### Method B: Check if commits are in main (fallback)

If no PR file exists:
```bash
git -C $PROJECT_DIR fetch origin main
```

```bash
git -C $WORKTREE_DIR rev-parse HEAD
```

Then check if that commit is in main:
```bash
git -C $PROJECT_DIR branch -r --contains <COMMIT_HASH> | grep origin/main
```

## Decision Logic

- **PR merged: true** OR **commit in main** → Proceed without confirmation
- Otherwise → Ask user to confirm

## Step 3: Cleanup

Execute these commands. Ignore errors - some resources may already be cleaned up.

### 3.1 Clean stale worktree references

```bash
git -C $PROJECT_DIR worktree prune
```

### 3.2 Remove worktree

```bash
test -d $WORKTREE_DIR && git -C $PROJECT_DIR worktree remove $WORKTREE_DIR --force || true
```

### 3.3 Remove agent directory

```bash
rm -rf $TAW_DIR/agents/$TASK_NAME || true
```

### 3.4 Delete local branch (if exists)

```bash
git -C $PROJECT_DIR branch --list $TASK_NAME | grep -q $TASK_NAME && git -C $PROJECT_DIR branch -D $TASK_NAME || true
```

### 3.5 Close tmux window

```bash
tmux kill-window
```

## Important

- Use `|| true` to prevent errors from stopping cleanup
- If any command fails, continue with the next step
- Report what was cleaned up at the end

Proceed with cleanup.
