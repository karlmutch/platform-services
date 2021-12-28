#!/bin/bash -x

if ( find /project -maxdepth 0 -empty | read v );
then
  echo "source code must be mounted into the /project directory"
  exit -1
fi

set -e
set -o pipefail

[ -e internal/gen/experimentsrv ] || mkdir -p internal/gen/experimentsrv
#[ -e vendor/github.com/karlmutch/platform-services/internal ] || mkdir -p vendor/github.com/karlmutch/platform-services/internal
#[ -e vendor/github.com/karlmutch/platform-services/internal/gen ] || ln -s `pwd`/internal/gen vendor/github.com/karlmutch/platform-services/internal/gen
protoc -Icmd/experimentsrv -I/usr/include/google --plugin=`which protoc-gen-go` --go_out=./internal/gen/experimentsrv --go_opt=paths=source_relative --plugin=`which protoc-gen-go-grpc` --go-grpc_out=./internal/gen/experimentsrv cmd/experimentsrv/experimentsrv.proto --go-grpc_opt=paths=source_relative

if [ "$1" == "gen" ]; then
    exit 0
fi

export HASH=`git rev-parse HEAD`

go mod vendor
go mod tidy

export SEMVER=`semver`
TAG_PARTS=$(echo $SEMVER | sed "s/-/\n-/g" | sed "s/\./\n\./g" | sed "s/+/\n+/g")
PATCH=""
for part in $TAG_PARTS
do
    start=`echo "$part" | cut -c1-1`
    if [ "$start" == "+" ]; then
        break
    fi
    if [ "$start" == "-" ]; then
        PATCH+=$part
    fi
done

flags='-X github.com/karlmutch/platform-services/internal/version.GitHash="$HASH" -X github.com/karlmutch/platform-services/internal/version.SemVer="$SEMVER"'
flags="$(eval echo $flags)"

mkdir -p cmd/experimentsrv/bin
CGO_ENABLED=0 go build -asmflags -trimpath -ldflags "$flags" -o cmd/experimentsrv/bin/experimentsrv cmd/experimentsrv/*.go
go build -asmflags -trimpath -ldflags "$flags" -race -o cmd/experimentsrv/bin/experimentsrv-race cmd/experimentsrv/*.go
CGO_ENABLED=0 go test -asmflags -trimpath -ldflags "$flags" -coverpkg="." -c -o cmd/experimentsrv/bin/experimentsrv-run-coverage cmd/experimentsrv/*.go
CGO_ENABLED=0 go test -asmflags -trimpath -ldflags "$flags" -coverpkg="." -c -o bin/experimentsrv-test-coverage cmd/experimentsrv/*.go
go test -asmflags -trimpath -ldflags "$flags" -race -c -o cmd/experimentsrv/bin/experimentsrv-test cmd/experimentsrv/*.go
if [ -z "$PATCH" ]; then
    if ! [ -z "${SEMVER}" ]; then
        if ! [ -z "${GITHUB_TOKEN}" ]; then
            github-release release --user karlmutch --repo platform-services --tag ${SEMVER} --pre-release || true
            github-release upload --user karlmutch --repo platform-services  --tag ${SEMVER} --name experimentsrv --file cmd/experimentsrv/bin/experimentsrv
        fi
    fi
fi
