#!/bin/bash

TIMESTAMP=$(date +%Y-%m-%d-%H%M%S)
OUTPUT_DIR="test-$TIMESTAMP"

show_usage() {
    cat << EOF
Usage: $0 [output_dir] [options]

Builds complete test library

Options:
    --help, -h      Show this help

Examples:
    $0                          # Build everything in ./test-library
    $0 /tmp/test                # Build in custom directory
EOF
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --help|-h)
            show_usage
            exit 0
            ;;
        -*)
            echo "Error: Unknown option $1"
            show_usage
            exit 1
            ;;
        *)
            OUTPUT_DIR="$1"
            shift
            ;;
    esac
done

echo "=========================================="
echo "Building Complete Test Library"
echo "=========================================="
echo "Output directory: $OUTPUT_DIR"
echo

# Create root directory
mkdir -p "$OUTPUT_DIR"

./build-issue-library.sh "$OUTPUT_DIR/by-issue"
./build-codec-library.sh "$OUTPUT_DIR/by-codec"
./build-media-library.sh "$OUTPUT_DIR/library"
