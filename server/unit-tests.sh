#!/bin/bash

go test ./pkg/... -count=1 -cover  $@ || exit 1
