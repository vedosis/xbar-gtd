#!/bin/bash

WORKSPACE=$BUILD_WORKSPACE_DIRECTORY
if [[ -z $WORKSPACE ]]; then
  echo "Missing workspace directory, aborting installation of precommit hooks"
  exit 0
fi

cd $WORKSPACE

if [ ! -e .git/hooks/pre-commit ]; then
  echo "📦 Installing git hooks... ($WORKSPACE)"
  mkdir -p .git/hooks
  ln -sf $WORKSPACE/tools/pre-commit.sh .git/hooks/pre-commit
fi

