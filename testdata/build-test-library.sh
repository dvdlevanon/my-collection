#!/bin/bash

# build-test-library.sh
# Master script to build complete test library with both hierarchies

SCRIPT_DIR=$(dirname "$(realpath "${BASH_SOURCE[0]}")")

show_usage() {
    cat << EOF
Usage: $0 [output_dir] [options]

Builds complete test library with two hierarchies:
1. by_issue/     - Test cases organized by problem type
2. lib1/lib2/lib3/ - Category-based media library structure

Options:
    --issue-only    Build only by_issue/ hierarchy
    --category-only Build only lib1/lib2/lib3/ hierarchy
    --clean         Remove existing test library before building
    --help, -h      Show this help

Default output: ./test-library

Examples:
    $0                          # Build everything in ./test-library
    $0 /tmp/test                # Build in custom directory
    $0 --clean                  # Clean and rebuild
    $0 --issue-only             # Only edge case tests
EOF
}

# Defaults
OUTPUT_DIR="./test-library"
BUILD_ISSUE=true
BUILD_CATEGORY=true
CLEAN=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --issue-only)
            BUILD_CATEGORY=false
            shift
            ;;
        --category-only)
            BUILD_ISSUE=false
            shift
            ;;
        --clean)
            CLEAN=true
            shift
            ;;
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
echo "Build by_issue: $BUILD_ISSUE"
echo "Build categories: $BUILD_CATEGORY"
echo

# Clean if requested
if [ "$CLEAN" = true ]; then
    if [ -d "$OUTPUT_DIR" ]; then
        echo "Cleaning existing test library..."
        rm -rf "$OUTPUT_DIR"
        echo "✓ Cleaned: $OUTPUT_DIR"
        echo
    fi
fi

# Create root directory
mkdir -p "$OUTPUT_DIR"
START_TIME=$(date +%s)

#------------------------------------------------------------------------------
# Build by_issue/ hierarchy
#------------------------------------------------------------------------------
if [ "$BUILD_ISSUE" = true ]; then
    echo ""
    echo "=========================================="
    echo "STEP 1: Building by_issue/ Hierarchy"
    echo "=========================================="
    
    if [ -f "$SCRIPT_DIR/build-by-issue-hierarchy.sh" ]; then
        bash "$SCRIPT_DIR/build-by-issue-hierarchy.sh" "$OUTPUT_DIR"
    else
        echo "✗ build-by-issue-hierarchy.sh not found"
        echo "  Creating placeholder..."
        mkdir -p "$OUTPUT_DIR/by_issue"
        echo "Script not found" > "$OUTPUT_DIR/by_issue/MISSING.txt"
    fi
fi

#------------------------------------------------------------------------------
# Build lib1/lib2/lib3/ category hierarchy
#------------------------------------------------------------------------------
if [ "$BUILD_CATEGORY" = true ]; then
    echo ""
    echo "=========================================="
    echo "STEP 2: Building Category Hierarchy"
    echo "=========================================="
    
    if [ -f "$SCRIPT_DIR/create-media-library.sh" ]; then
        bash "$SCRIPT_DIR/create-media-library.sh" "$OUTPUT_DIR"
    else
        echo "✗ create-media-library.sh not found"
        echo "  Creating placeholder structure..."
        mkdir -p "$OUTPUT_DIR/lib1" "$OUTPUT_DIR/lib2" "$OUTPUT_DIR/lib3"
        echo "Script not found" > "$OUTPUT_DIR/lib1/MISSING.txt"
    fi
fi

#------------------------------------------------------------------------------
# Generate master README
#------------------------------------------------------------------------------
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

cat > "$OUTPUT_DIR/README.md" << EOF
# Complete Test Library for my-collection

Generated on: $(date)
Build time: ${DURATION} seconds

## Two-Hierarchy Structure

This test library contains two separate hierarchies designed to test different aspects of media library management:

### 1. Issue-Based Hierarchy (\`by_issue/\`)

Organized by problem type to test edge case handling:

\`\`\`
by_issue/
├── codec_variations/      # Different codecs and containers
├── audio_tracks/          # Multi-audio, channels, languages
├── subtitle_variations/   # External, embedded, multi-language subs
├── resolution_variations/ # Various resolutions and aspect ratios
├── broken_files/          # Corrupted, truncated, malformed files
└── metadata_issues/       # Frame rates, color spaces, interlacing
\`\`\`

**Purpose:** Test specific edge cases and error handling

### 2. Category-Based Hierarchy (\`lib1/\`, \`lib2/\`, \`lib3/\`)

Organized by content type to test library organization:

\`\`\`
lib1/  # Movies
lib2/  # TV Series
lib3/  # Short Clips
by_codec/  # Same content in different codecs
\`\`\`

**Purpose:** Test category detection, filesystem sync, and codec conversion features

## Statistics

$(if [ "$BUILD_ISSUE" = true ] && [ -d "$OUTPUT_DIR/by_issue" ]; then
    echo "### by_issue/ Hierarchy"
    echo ""
    find "$OUTPUT_DIR/by_issue" -type d -maxdepth 1 -mindepth 1 2>/dev/null | sort | while read dir; do
        name=$(basename "$dir")
        count=$(find "$dir" -type d -mindepth 1 2>/dev/null | wc -l | tr -d ' ')
        files=$(find "$dir" -type f \( -name "*.mp4" -o -name "*.mkv" -o -name "*.avi" -o -name "*.webm" \) 2>/dev/null | wc -l | tr -d ' ')
        echo "- **$name**: $count test cases, $files video files"
    done
    echo ""
fi)

$(if [ "$BUILD_CATEGORY" = true ] && [ -d "$OUTPUT_DIR/lib1" ]; then
    echo "### Category Hierarchy"
    echo ""
    echo "- **lib1/**: $(find "$OUTPUT_DIR/lib1" -type f \( -name "*.mp4" -o -name "*.mkv" -o -name "*.webm" \) 2>/dev/null | wc -l | tr -d ' ') movies"
    echo "- **lib2/**: $(find "$OUTPUT_DIR/lib2" -type f \( -name "*.mp4" -o -name "*.mkv" -o -name "*.webm" \) 2>/dev/null | wc -l | tr -d ' ') TV episodes"
    echo "- **lib3/**: $(find "$OUTPUT_DIR/lib3" -type f \( -name "*.mp4" -o -name "*.mkv" -o -name "*.webm" \) 2>/dev/null | wc -l | tr -d ' ') clips"
    if [ -d "$OUTPUT_DIR/by_codec" ]; then
        echo "- **by_codec/**: $(find "$OUTPUT_DIR/by_codec" -type f \( -name "*.mp4" -o -name "*.mkv" -o -name "*.webm" \) 2>/dev/null | wc -l | tr -d ' ') codec test files"
    fi
    echo ""
fi)

**Total video files:** $(find "$OUTPUT_DIR" -type f \( -name "*.mp4" -o -name "*.mkv" -o -name "*.avi" -o -name "*.webm" -o -name "*.mpg" \) 2>/dev/null | wc -l | tr -d ' ')

## Testing Workflow

### For Edge Case Testing (by_issue/)

1. Run media scanner on each subdirectory
2. Verify proper detection and handling
3. Check error recovery for broken files
4. Validate metadata extraction

### For Category/Codec Testing (lib1/lib2/lib3/)

1. Test category detection and organization
2. Test codec conversion workflows
3. Verify filesystem sync operations
4. Compare quality across codec variations

## Cleanup

To remove the test library:

\`\`\`bash
rm -rf $OUTPUT_DIR
\`\`\`

To rebuild:

\`\`\`bash
$0 --clean
\`\`\`

## Notes

- All videos are short (10-20 seconds) for fast testing
- Files are optimized for speed, not quality
- Some test cases intentionally contain errors
- Not all files are fully generated (some are placeholders)

For detailed information about specific test cases, see the README files in each subdirectory.
EOF

echo ""
echo "=========================================="
echo "Build Complete!"
echo "=========================================="
echo "Location: $OUTPUT_DIR"
echo "Build time: ${DURATION} seconds"
echo ""
echo "Statistics:"

if [ "$BUILD_ISSUE" = true ] && [ -d "$OUTPUT_DIR/by_issue" ]; then
    echo "  by_issue test cases: $(find "$OUTPUT_DIR/by_issue" -type d -mindepth 2 2>/dev/null | wc -l | tr -d ' ')"
fi

if [ "$BUILD_CATEGORY" = true ]; then
    if [ -d "$OUTPUT_DIR/lib1" ]; then
        echo "  Category files: $(find "$OUTPUT_DIR/lib1" "$OUTPUT_DIR/lib2" "$OUTPUT_DIR/lib3" -type f \( -name "*.mp4" -o -name "*.mkv" -o -name "*.webm" \) 2>/dev/null | wc -l | tr -d ' ')"
    fi
    if [ -d "$OUTPUT_DIR/by_codec" ]; then
        echo "  Codec test files: $(find "$OUTPUT_DIR/by_codec" -type f \( -name "*.mp4" -o -name "*.mkv" -o -name "*.webm" \) 2>/dev/null | wc -l | tr -d ' ')"
    fi
fi

echo "  Total video files: $(find "$OUTPUT_DIR" -type f \( -name "*.mp4" -o -name "*.mkv" -o -name "*.avi" -o -name "*.webm" -o -name "*.mpg" \) 2>/dev/null | wc -l | tr -d ' ')"
echo ""
echo "See $OUTPUT_DIR/README.md for complete documentation"
echo ""
echo "Quick test commands:"
echo "  tree -L 2 $OUTPUT_DIR              # View structure"
echo "  du -sh $OUTPUT_DIR                 # Check size"
echo "  find $OUTPUT_DIR -name '*.mp4'     # List all MP4s"
