#!/usr/bin/env bash

set -e

PACKAGE="azssh"
platforms=("aix/ppc64" "darwin/386" "darwin/amd64" "dragonfly/amd64" "freebsd/386" "freebsd/amd64" "freebsd/arm" "freebsd/arm64" "illumos/amd64" "linux/386" "linux/amd64" "linux/arm" "linux/arm64" "linux/mips" "linux/mips64" "linux/mips64le" "linux/mipsle" "linux/ppc64" "linux/ppc64le" "linux/riscv64" "linux/s390x" "netbsd/386" "netbsd/amd64" "netbsd/arm" "netbsd/arm64" "openbsd/386" "openbsd/amd64" "openbsd/arm" "openbsd/arm64" "solaris/amd64" "windows/386" "windows/amd64" "windows/arm")

rm -rf release
mkdir -p release

for platform in "${platforms[@]}"
do
    echo "Building for platform: $platform"
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$PACKAGE'-'$GOOS'-'$GOARCH
    if [ "$GOOS" = "windows" ]; then
        output_name+=".exe"
    fi

    env GOOS="$GOOS" GOARCH="$GOARCH" go build -o release/$output_name
    if [ $? -ne 0 ]; then
        echo "An error has occurred! Aborting the script execution..."
        exit 1
    fi
done
