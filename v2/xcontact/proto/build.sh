#!/usr/bin/env sh

SCRIPT_DIR=$(dirname "$0")

protoc base.proto --go_out="$SCRIPT_DIR" && mv "$SCRIPT_DIR"/code.olapie.com/sugar/v2/xcontact/base.pb.go "$SCRIPT_DIR"/../
rm -r "$SCRIPT_DIR"/code.olapie.com
