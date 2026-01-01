#!/bin/bash

# generate-broken-files.sh
# Generates test cases for corrupted, truncated, and problematic files

SCRIPT_DIR=$(dirname "$(realpath "${BASH_SOURCE[0]}")")

show_usage() {
    cat << EOF
Usage: $0 <output_dir>

Generates broken/problematic file test cases including:
- Truncated/incomplete files
- Zero-byte files
- Wrong file extensions
- Corrupted headers
- Very short clips (<1 second)
- Extremely long files (simulated)
- Missing metadata

Output structure:
    output_dir/
    ├── truncated_video/
    ├── zero_byte_file/
    ├── wrong_extension/
    └── ...
EOF
}

if [ $# -ne 1 ]; then
    show_usage
    exit 1
fi

OUTPUT_DIR="$1"
mkdir -p "$OUTPUT_DIR"

echo "=========================================="
echo "Generating Broken/Problematic File Test Cases"
echo "=========================================="
echo

#------------------------------------------------------------------------------
# Test Case 1: Zero-byte file
#------------------------------------------------------------------------------
echo "Creating: Zero-byte file"
test_dir="$OUTPUT_DIR/zero_byte_file"
mkdir -p "$test_dir"

touch "$test_dir/movie.mp4"

cat > "$test_dir/README.txt" << EOF
Test Case: Zero-Byte File

File exists but has no content (0 bytes).

Purpose: Test handling of empty/placeholder files.
Should be detected and skipped or flagged as error.
EOF

echo "  ✓ Created: $test_dir/movie.mp4 (0 bytes)"
echo

#------------------------------------------------------------------------------
# Test Case 2: Truncated video (incomplete)
#------------------------------------------------------------------------------
echo "Creating: Truncated/incomplete video"
test_dir="$OUTPUT_DIR/truncated_video"
mkdir -p "$test_dir"

# Generate full video then truncate it
temp_video="$test_dir/temp_full.mp4"
ffmpeg -y -f lavfi -i "color=c=gray:s=1280x720:d=15" \
       -f lavfi -i "sine=frequency=440:duration=15" \
       -c:v libx264 -preset ultrafast -c:a aac \
       -t 15 "$temp_video"|| exit 1

# Truncate to 30% of original size (simulates incomplete download)
size=$(stat -f%z "$temp_video" 2>/dev/null || stat -c%s "$temp_video" 2>/dev/null)
truncate_size=$((size * 30 / 100))
dd if="$temp_video" of="$test_dir/movie.mp4" bs=1 count=$truncate_size 2>/dev/null

rm -f "$temp_video"

cat > "$test_dir/README.txt" << EOF
Test Case: Truncated Video

Video file is incomplete - truncated to ~30% of original size.
Simulates interrupted downloads or corrupted files.

Purpose: Test handling of incomplete files.
Should detect corruption and either skip or report error.
EOF

echo "  ✓ Created: $test_dir/movie.mp4 (truncated)"
echo

#------------------------------------------------------------------------------
# Test Case 3: Wrong file extension
#------------------------------------------------------------------------------
echo "Creating: Wrong file extension"
test_dir="$OUTPUT_DIR/wrong_extension"
mkdir -p "$test_dir"

# Create MP4 but name it .avi
ffmpeg -y -f lavfi -i "color=c=yellow:s=1280x720:d=10" \
       -f lavfi -i "sine=frequency=440:duration=10" \
       -c:v libx264 -preset ultrafast -c:a aac \
       -t 10 "$test_dir/movie.avi"|| exit 1

# Create AVI but name it .mp4
ffmpeg -y -f lavfi -i "color=c=cyan:s=1280x720:d=10" \
       -f lavfi -i "sine=frequency=440:duration=10" \
       -c:v mpeg4 -c:a libmp3lame \
       -t 10 "$test_dir/video.mp4"|| exit 1

cat > "$test_dir/README.txt" << EOF
Test Case: Wrong File Extension

Contains:
- movie.avi - Actually MP4 format (H.264/AAC)
- video.mp4 - Actually AVI format (MPEG-4/MP3)

Purpose: Test container detection based on actual format, not extension.
Should detect real format and handle appropriately.
EOF

echo "  ✓ Created: $test_dir/ (mismatched extensions)"
echo

#------------------------------------------------------------------------------
# Test Case 4: Very short clip (<1 second)
#------------------------------------------------------------------------------
echo "Creating: Very short clip (<1 second)"
test_dir="$OUTPUT_DIR/very_short_clip"
mkdir -p "$test_dir"

ffmpeg -y -f lavfi -i "color=c=pink:s=1280x720:d=0.5" \
       -f lavfi -i "sine=frequency=440:duration=0.5" \
       -c:v libx264 -preset ultrafast -c:a aac \
       -t 0.5 "$test_dir/clip.mp4"|| exit 1

cat > "$test_dir/README.txt" << EOF
Test Case: Very Short Clip

Video is only 0.5 seconds long.

Purpose: Test handling of extremely short videos.
Some media tools have issues with <1s videos.
EOF

echo "  ✓ Created: $test_dir/clip.mp4 (0.5 seconds)"
echo

#------------------------------------------------------------------------------
# Test Case 5: Missing duration metadata
#------------------------------------------------------------------------------
echo "Creating: Missing duration metadata"
test_dir="$OUTPUT_DIR/missing_duration"
mkdir -p "$test_dir"

# Create video without proper duration metadata
ffmpeg -y -f lavfi -i "color=c=brown:s=1280x720:d=15" \
       -c:v libx264 -preset ultrafast -copyts -fflags +bitexact \
       "$test_dir/video.mp4"|| exit 1

cat > "$test_dir/README.txt" << EOF
Test Case: Missing Duration Metadata

Video may not have accurate duration in metadata.

Purpose: Test handling when duration needs to be calculated
by scanning the entire file.
EOF

echo "  ✓ Created: $test_dir/video.mp4"
echo

#------------------------------------------------------------------------------
# Test Case 6: Corrupted file header
#------------------------------------------------------------------------------
echo "Creating: Corrupted file header"
test_dir="$OUTPUT_DIR/corrupted_header"
mkdir -p "$test_dir"

# Create valid video then corrupt first 100 bytes
temp_video="$test_dir/temp.mp4"
ffmpeg -y -f lavfi -i "color=c=black:s=1280x720:d=10" \
       -c:v libx264 -preset ultrafast \
       -t 10 "$temp_video"|| exit 1

# Corrupt header by overwriting first 100 bytes with random data
dd if=/dev/urandom of="$test_dir/movie.mp4" bs=100 count=1 2>/dev/null
dd if="$temp_video" of="$test_dir/movie.mp4" bs=1 skip=100 seek=100 2>/dev/null

rm -f "$temp_video"

cat > "$test_dir/README.txt" << EOF
Test Case: Corrupted File Header

First 100 bytes are corrupted/invalid.
File may not be recognized as valid video.

Purpose: Test error handling for corrupted files.
Should detect and skip or report error.
EOF

echo "  ✓ Created: $test_dir/movie.mp4 (corrupted header)"
echo

#------------------------------------------------------------------------------
# Test Case 7: Text file with video extension
#------------------------------------------------------------------------------
echo "Creating: Text file with video extension"
test_dir="$OUTPUT_DIR/text_file_as_video"
mkdir -p "$test_dir"

cat > "$test_dir/fake_video.mp4" << EOF
This is not actually a video file.
It's just a text file with a .mp4 extension.
This should be detected and rejected.
EOF

cat > "$test_dir/README.txt" << EOF
Test Case: Text File Masquerading as Video

A plain text file with .mp4 extension.

Purpose: Test that actual file content is validated,
not just the extension.
EOF

echo "  ✓ Created: $test_dir/fake_video.mp4 (text file)"
echo

#------------------------------------------------------------------------------
# Test Case 8: Symlink to non-existent file
#------------------------------------------------------------------------------
echo "Creating: Broken symlink"
test_dir="$OUTPUT_DIR/broken_symlink"
mkdir -p "$test_dir"

# Create symlink to non-existent file
ln -s "/nonexistent/path/to/video.mp4" "$test_dir/movie.mp4" 2>/dev/null

cat > "$test_dir/README.txt" << EOF
Test Case: Broken Symlink

Symlink points to non-existent file.

Purpose: Test handling of broken symlinks.
Should detect and skip or report error.
EOF

echo "  ✓ Created: $test_dir/movie.mp4 (broken symlink)"
echo

#------------------------------------------------------------------------------
# Test Case 9: File with special characters in name
#------------------------------------------------------------------------------
echo "Creating: Special characters in filename"
test_dir="$OUTPUT_DIR/special_characters"
mkdir -p "$test_dir"

ffmpeg -y -f lavfi -i "color=c=white:s=1280x720:d=10" \
       -c:v libx264 -preset ultrafast \
       -t 10 "$test_dir/movie [2024] (1080p) {H.264}.mp4"|| exit 1

ffmpeg -y -f lavfi -i "color=c=black:s=1280x720:d=10" \
       -c:v libx264 -preset ultrafast \
       -t 10 "$test_dir/file with spaces & symbols!.mp4"|| exit 1

cat > "$test_dir/README.txt" << EOF
Test Case: Special Characters in Filename

Files with brackets, parentheses, spaces, and symbols.

Purpose: Test filename parsing and shell escaping.
Common in real-world media libraries.
EOF

echo "  ✓ Created: $test_dir/ (special characters)"
echo

#------------------------------------------------------------------------------
# Test Case 10: Hidden files (dotfiles)
#------------------------------------------------------------------------------
echo "Creating: Hidden system files"
test_dir="$OUTPUT_DIR/hidden_files"
mkdir -p "$test_dir"

# Create actual video
ffmpeg -y -f lavfi -i "color=c=blue:s=1280x720:d=10" \
       -c:v libx264 -preset ultrafast \
       -t 10 "$test_dir/movie.mp4"|| exit 1

# Create hidden system files that should be ignored
touch "$test_dir/.DS_Store"
touch "$test_dir/Thumbs.db"
touch "$test_dir/.hidden_file"

cat > "$test_dir/README.txt" << EOF
Test Case: Hidden System Files

Directory contains:
- movie.mp4 (actual video)
- .DS_Store (macOS metadata)
- Thumbs.db (Windows thumbnail cache)
- .hidden_file (hidden file)

Purpose: Test filtering of system/hidden files.
These should be ignored during media scanning.
EOF

echo "  ✓ Created: $test_dir/ (with hidden files)"
echo

echo "=========================================="
echo "Broken/problematic file test cases generated in: $OUTPUT_DIR"
echo "Total test cases: 10"
echo "=========================================="
