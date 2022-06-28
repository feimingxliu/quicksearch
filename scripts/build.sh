#!/bin/bash
export VERSION=`git describe --tags --always`
export BUILD_DATE=`date -u '+%Y-%m-%d %I:%M:%S'`
export COMMIT_HASH=`git rev-parse HEAD`
export BRANCH=`git branch --show-current`
export LDFLAGS="-w -s -X github.com/feimingxliu/quicksearch/pkg/about.Branch=${BRANCH} -X github.com/feimingxliu/quicksearch/pkg/about.Version=${VERSION} -X 'github.com/feimingxliu/quicksearch/pkg/about.BuildDate=${BUILD_DATE}' -X github.com/feimingxliu/quicksearch/pkg/about.CommitHash=${COMMIT_HASH}"
go build -ldflags="$LDFLAGS" -o bin/quicksearch github.com/feimingxliu/quicksearch/cmd/quicksearch