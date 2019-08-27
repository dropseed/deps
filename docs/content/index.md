---
description: "Documentation and examples for automatically updating your dependencies with deps"
---

# Overview

Deps is a command line tool that runs in CI to automate dependency updates. It will automatically create branches, commit changes, and send you pull requests with links to release notes.

## 1. Try it locally

Installing `deps` on your machine is the easiest way to get started and see how things work.
You don't need an API token to run it locally.

[Install it on your computer →](/local/)

## 2. Automate in CI

When you're ready to automate your dependency updates,
you'll install `deps` to run in your CI provider.
Each CI system is different,
but generally all you need to do is set some environment variables and add a scheduled job for dependency updates.

[Set it up in your CI →](/ci/)
