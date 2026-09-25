# Upstream synchronization

This fork tracks NVIDIA/Personal-AI-Router on the develop branch.

## Automatic sync

.github/workflows/sync-upstream.yml runs every six hours and can also be started manually from GitHub Actions.

When NVIDIA has new commits, the workflow performs a normal non-fast-forward merge into this fork's develop branch. Fork-specific changes remain in the history.

If the merge is clean, the workflow pushes the synchronized branch. If upstream and fork changes conflict, the merge is aborted and nothing is pushed. Resolve that conflict locally, push the result, and scheduled synchronization can continue normally.

## Manual sync

    git remote add upstream https://github.com/NVIDIA/Personal-AI-Router.git
    git fetch upstream develop
    git checkout develop
    git merge --no-edit --no-ff upstream/develop
    git push origin develop

The fork's legacy GPU, Tailscale/MeshLLM, and node-explorer changes are intentionally kept as fork-local commits. When resolving conflicts, preserve upstream changes that are unrelated to those features and reapply the fork-specific behavior where lines overlap.
