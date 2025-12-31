#!/bin/bash

# build-by-issue-hierarchy.sh
# Builds the by_issue/ hierarchy by running all issue generator scripts

SCRIPT_DIR=$(dirname "$(realpath "${BASH_SOURCE[0]}")")

show_usage() {
    cat << EOF
Usage: $0 [output_root_dir]

Builds the by_issue/ test hierarchy with all edge case test generators.

Creates:
    output_root_dir/by_issue/
    ├── codec_variations/
    ├── audio_tracks/
    ├── subtitle_variations/
    ├── resolution_variations/
    ├── broken_files/
    └── metadata_issues/

Default output: ./test-library/by_issue
EOF
}

# Parse arguments
OUTPUT_ROOT="${1:-./test-library}"
BY_ISSUE_DIR="$OUTPUT_ROOT/by_issue"

if [[ "$1" == "--help" ]] || [[ "$1" == "-h" ]]; then
    show_usage
    exit 0
fi

echo "=========================================="
echo "Building by_issue Test Hierarchy"
echo "=========================================="
echo "Output directory: $BY_ISSUE_DIR"
echo

# Create base directory
mkdir -p "$BY_ISSUE_DIR"

# Run each issue generator
generators=(
    "generate-codec-variations.sh:codec_variations:Codec Variations"
    "generate-audio-tracks.sh:audio_tracks:Multi-Audio Tracks"
    "generate-subtitle-variations.sh:subtitle_variations:Subtitle Variations"
    "generate-resolution-variations.sh:resolution_variations:Resolution/Aspect Ratios"
    "generate-broken-files.sh:broken_files:Broken/Corrupted Files"
    "generate-metadata-issues.sh:metadata_issues:Metadata/Encoding Issues"
)

total_generators=${#generators[@]}
current=0

for generator_info in "${generators[@]}"; do
    IFS=':' read -r script_name dir_name description <<< "$generator_info"
    current=$((current + 1))
    
    echo ""
    echo "[$current/$total_generators] Generating: $description"
    echo "----------------------------------------"
    
    script_path="$SCRIPT_DIR/$script_name"
    output_dir="$BY_ISSUE_DIR/$dir_name"
    
    if [ -f "$script_path" ]; then
        if bash "$script_path" "$output_dir"; then
            echo "✓ Completed: $description"
        else
            echo "✗ Failed: $description"
        fi
    else
        echo "✗ Generator not found: $script_path"
        echo "  Creating placeholder directory..."
        mkdir -p "$output_dir"
        echo "Generator script not found: $script_name" > "$output_dir/MISSING.txt"
    fi
done

echo ""
echo "=========================================="
echo "by_issue Hierarchy Complete"
echo "=========================================="
echo "Location: $BY_ISSUE_DIR"
echo ""

# Generate summary
cat > "$BY_ISSUE_DIR/README.md" << EOF
# Test Cases Organized by Issue Type

Generated on: $(date)

This hierarchy contains test cases organized by the type of problem they represent.
Use these to test specific edge case handling in your media library software.

## Structure:

$(find "$BY_ISSUE_DIR" -type d -maxdepth 1 -mindepth 1 | sort | while read dir; do
    name=$(basename "$dir")
    count=$(find "$dir" -type d -mindepth 1 | wc -l | tr -d ' ')
    echo "- **$name/** ($count test cases)"
done)

## Purpose:

Each subdirectory contains specific edge cases that commonly cause issues in media libraries:

- **codec_variations/**: Different video/audio codecs and containers
- **audio_tracks/**: Multiple audio tracks, channels, languages
- **subtitle_variations/**: External, embedded, multi-language subtitles
- **resolution_variations/**: Various resolutions, aspect ratios, orientations
- **broken_files/**: Corrupted, truncated, malformed files
- **metadata_issues/**: Frame rates, color spaces, interlacing, etc.

## Usage:

Run your media library scanner against these test cases to verify:
1. Proper codec detection
2. Audio track selection
3. Subtitle handling
4. Resolution scaling
5. Error recovery
6. Metadata extraction

Each test case directory contains a README.txt explaining the specific issue being tested.
EOF

echo "Summary:"
find "$BY_ISSUE_DIR" -type d -maxdepth 1 -mindepth 1 | sort | while read dir; do
    name=$(basename "$dir")
    if [ "$name" != "by_issue" ]; then
        count=$(find "$dir" -type d -mindepth 1 2>/dev/null | wc -l | tr -d ' ')
        echo "  - $name: $count test cases"
    fi
done

echo ""
echo "See $BY_ISSUE_DIR/README.md for details"
