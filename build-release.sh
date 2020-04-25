#!/usr/bin/env bash

set -e

platforms=("windows/amd64" "windows/386" "darwin/386" "darwin/amd64" "linux/386" "linux/amd64" "linux/arm" "linux/arm64" "freebsd/386" "freebsd/amd64")
PACKAGE="azssh"

rm -rf release
mkdir -p release

for platform in "${platforms[@]}"
do
    echo "Building for platform:" $platform
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$PACKAGE'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+=".exe"
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o release/$output_name
    if [ $? -ne 0 ]; then
        echo "An error has occurred! Aborting the script execution..."
        exit 1
    fi
done
