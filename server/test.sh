#!/bin/bash

go test ./pkg/... -cover || exit 1

./test/scripts/run_integration_tests.sh || exit 1
