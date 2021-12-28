#!/bin/bash -x

if ( find /project -maxdepth 0 -empty | read v );
then
  echo "source code must be mounted into the /project directory"
  exit 990
fi

set -e
set -o pipefail

[ -e internal/gen ] || mkdir internal/gen
swagger generate server -q -t internal/gen -f cmd/timesrv/swagger.yaml --exclude-main -A timesrv
# go install gen/...
#[ -e vendor/github.com/karlmutch/platform-services/internal ] || mkdir -p vendor/github.com/karlmutch/platform-services/internal
#[ -e vendor/github.com/karlmutch/platform-services/internal/gen ] || ln -s `pwd`/internal/gen vendor/github.com/karlmutch/platform-services/internal/gen
if [ "$1" == "gen" ]; then
    exit 0
fi

export HASH=`git rev-parse HEAD`

mkdir -p cmd/timesrv/bin
go build -asmflags -trimpath -ldflags "-X github.com/karlmutch/platform-services/internal/version.GitHash=$HASH" -o cmd/timesrv/bin/timesrv cmd/timesrv/*.go
go build -asmflags -trimpath -ldflags "-X github.com/karlmutch/platform-services/internal/version.GitHash=$HASH" -race -o cmd/timesrv/bin/timesrv-race cmd/timesrv/*.go
if ! [ -z "${TRAVIS_TAG}" ]; then
    if ! [ -z "${GITHUB_TOKEN}" ]; then
        github-release release --user karlmutch --repo platform-services --tag ${TRAVIS_TAG} --pre-release || true
        github-release upload --user karlmutch --repo platform-services  --tag ${TRAVIS_TAG} --name timesrv --file cmd/timesrv/bin/timesrv
    fi
fi
