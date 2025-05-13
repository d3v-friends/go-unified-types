#!/bin/zsh
set -e

GO_DEF_PATH=$1

# delete oly definitions
rm "$GO_DEF_PATH/*.pb.go" || true

gc-cli
# delete mac dummy files
rm "$GO_DEF_PATH/._*" || true
git add "$GO_DEF_PATH"
echo "success";