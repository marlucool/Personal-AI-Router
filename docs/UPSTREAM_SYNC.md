# Upstream synchronization

This fork tracks NVIDIA/Personal-AI-Router on the develop branch.

## Automatic sync

The .github/workflows/sync-upstream.yml workflow runs every six hours and can also be started manually from GitHub Actions.

When NVIDIA has new commits, the workflow performs a normal non-fast-forward merge into the fork's develop branch. Fork-specific changes remain in the history.

If the merge is clean, the workflow pushes the synchronized branch. If upstream and fork changes conflict, the merge is aborted and nothing is pushed. Resolve the conflict locally, test it, and push the resolved merge before the next scheduled sync.

## Manual sync

    git remote add upstream https://github.com/NVIDIA/Personal-AI-Router.git
    git fetch upstream develop
    git checkout develop
    git merge --no-edit --no-ff upstream/develop
    git push origin develop

The sync workflow does not force-push develop and does not rewrite active feature branches.
