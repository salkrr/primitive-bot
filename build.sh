#!/bin/bash

DIR="builds"
mkdir -p $DIR
cd "$DIR" || return 1

ARCH=("386" "amd64" "arm" "arm64")
OS=("linux" "darwin" "windows")
for arch in "${ARCH[@]}"; do
  for os in "${OS[@]}"; do
    BIN_PATH="primitive-bot"
    ARCHIVE_PATH="primitive_bot_${os}_${arch}.tar.gz"
    if [[ "$os" == "windows" ]]; then
      BIN_PATH="$BIN_PATH.exe"
    fi

    env GOARCH="$arch" GOOS="$os" go build -o "$BIN_PATH" ../cmd/primitive-bot/*.go
    if [ $? -eq 0 ]; then
      tar cfz "$ARCHIVE_PATH" "$BIN_PATH"
      rm "$BIN_PATH"
    fi
  done
done
