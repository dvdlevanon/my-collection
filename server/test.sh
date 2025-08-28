#!/bin/bash

go test ./pkg/... -cover 

./test/scripts/run_integration_tests.sh
