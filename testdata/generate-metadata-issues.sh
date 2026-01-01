#!/bin/bash

# generate-metadata-issues.sh
# Generates test cases for metadata issues, frame rates, color spaces, etc.

SCRIPT_DIR=$(dirname "$(realpath "${BASH_SOURCE[0]}")")

show_usage() {
    cat << EOF
Usage: $0 <output_dir>

Generates metadata/encoding issue test cases including:
- Variable frame rate (VFR) videos
- Non-standard frame rates (23.976, 29.97, 50, 60)
- Different color spaces (BT.709, BT.601)
- Interlaced vs progressive
- Different pixel formats
- Missing/wrong metadata

Output structure:
    output_dir/
    ├── vfr_screen_recording/
    ├── fps_23976/
    ├── fps_2997/
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
echo "Generating Metadata/Encoding Issue Test Cases"
echo "=========================================="
echo

#------------------------------------------------------------------------------
# Test Case 1: Variable Frame Rate (VFR)
#------------------------------------------------------------------------------
echo "Creating: Variable frame rate video"
test_dir="$OUTPUT_DIR/vfr_screen_recording"
mkdir -p "$test_dir"

# VFR is hard to generate properly, but we can create a video and document it
ffmpeg -y -f lavfi -i "color=c=navy:s=1920x1080:d=15" \
       -f lavfi -i "sine=frequency=440:duration=15" \
       -c:v libx264 -preset ultrafast -c:a aac \
       -vsync vfr \
       -t 15 "$test_dir/recording.mp4"|| exit 1

cat > "$test_dir/README.txt" << EOF
Test Case: Variable Frame Rate (VFR)

Screen recording with variable frame rate.
Common in OBS, screen capture software.

Purpose: Test VFR detection and conversion to CFR.
VFR can cause sync issues in some players/editors.
EOF

echo "  ✓ Created: $test_dir/recording.mp4"
echo

#------------------------------------------------------------------------------
# Test Case 2: 23.976 fps (Film standard)
#------------------------------------------------------------------------------
echo "Creating: 23.976 fps (film cadence)"
test_dir="$OUTPUT_DIR/fps_23976"
mkdir -p "$test_dir"

ffmpeg -y -f lavfi -i "color=c=sepia:s=1920x1080:d=15:rate=24000/1001" \
       -f lavfi -i "sine=frequency=440:duration=15" \
       -c:v libx264 -preset ultrafast -c:a aac \
       -r 24000/1001 \
       -t 15 "$test_dir/movie.mp4"|| exit 1

cat > "$test_dir/README.txt" << EOF
Test Case: 23.976 fps

NTSC film cadence (23.976 fps).
Common in cinema releases and streaming.

Purpose: Test proper handling of fractional frame rates.
EOF

echo "  ✓ Created: $test_dir/movie.mp4 (23.976 fps)"
echo

#------------------------------------------------------------------------------
# Test Case 3: 29.97 fps (NTSC)
#------------------------------------------------------------------------------
echo "Creating: 29.97 fps (NTSC)"
test_dir="$OUTPUT_DIR/fps_2997"
mkdir -p "$test_dir"

ffmpeg -y -f lavfi -i "color=c=teal:s=1280x720:d=15:rate=30000/1001" \
       -f lavfi -i "sine=frequency=440:duration=15" \
       -c:v libx264 -preset ultrafast -c:a aac \
       -r 30000/1001 \
       -t 15 "$test_dir/video.mp4"|| exit 1

cat > "$test_dir/README.txt" << EOF
Test Case: 29.97 fps

NTSC television frame rate (29.97 fps).
Common in North American broadcasts.

Purpose: Test NTSC frame rate handling.
EOF

echo "  ✓ Created: $test_dir/video.mp4 (29.97 fps)"
echo

#------------------------------------------------------------------------------
# Test Case 4: 50 fps (PAL)
#------------------------------------------------------------------------------
echo "Creating: 50 fps (PAL)"
test_dir="$OUTPUT_DIR/fps_50"
mkdir -p "$test_dir"

ffmpeg -y -f lavfi -i "color=c=lime:s=1280x720:d=10:rate=50" \
       -f lavfi -i "sine=frequency=440:duration=10" \
       -c:v libx264 -preset ultrafast -c:a aac \
       -r 50 \
       -t 10 "$test_dir/video.mp4"|| exit 1

cat > "$test_dir/README.txt" << EOF
Test Case: 50 fps

PAL high frame rate (50 fps).
Common in European broadcasts and sports.

Purpose: Test high frame rate handling.
EOF

echo "  ✓ Created: $test_dir/video.mp4 (50 fps)"
echo

#------------------------------------------------------------------------------
# Test Case 5: 60 fps (High frame rate)
#------------------------------------------------------------------------------
echo "Creating: 60 fps (gaming/action)"
test_dir="$OUTPUT_DIR/fps_60"
mkdir -p "$test_dir"

ffmpeg -y -f lavfi -i "color=c=red:s=1920x1080:d=10:rate=60" \
       -f lavfi -i "sine=frequency=440:duration=10" \
       -c:v libx264 -preset ultrafast -c:a aac \
       -r 60 \
       -t 10 "$test_dir/gameplay.mp4"|| exit 1

cat > "$test_dir/README.txt" << EOF
Test Case: 60 fps

High frame rate video (60 fps).
Common in gaming footage and action cameras.

Purpose: Test high frame rate handling and potential downscaling.
EOF

echo "  ✓ Created: $test_dir/gameplay.mp4 (60 fps)"
echo

#------------------------------------------------------------------------------
# Test Case 6: Interlaced video
#------------------------------------------------------------------------------
echo "Creating: Interlaced video (1080i)"
test_dir="$OUTPUT_DIR/interlaced_1080i"
mkdir -p "$test_dir"

ffmpeg -y -f lavfi -i "color=c=purple:s=1920x1080:d=15" \
       -f lavfi -i "sine=frequency=440:duration=15" \
       -c:v libx264 -preset ultrafast -c:a aac \
       -flags +ildct+ilme \
       -t 15 "$test_dir/broadcast.mp4"|| exit 1

cat > "$test_dir/README.txt" << EOF
Test Case: Interlaced Video

1080i interlaced video.
Common in older broadcasts and cameras.

Purpose: Test interlace detection and deinterlacing.
Should convert to progressive scan.
EOF

echo "  ✓ Created: $test_dir/broadcast.mp4 (interlaced)"
echo

#------------------------------------------------------------------------------
# Test Case 7: Different color spaces
#------------------------------------------------------------------------------
echo "Creating: Different color spaces (BT.709 vs BT.601)"
test_dir="$OUTPUT_DIR/color_spaces"
mkdir -p "$test_dir"

# BT.709 (HD standard)
ffmpeg -y -f lavfi -i "color=c=0xFF0000:s=1920x1080:d=10" \
       -f lavfi -i "sine=frequency=440:duration=10" \
       -c:v libx264 -preset ultrafast -c:a aac \
       -colorspace bt709 -color_primaries bt709 -color_trc bt709 \
       -t 10 "$test_dir/bt709_hd.mp4"|| exit 1

# BT.601 (SD standard)
ffmpeg -y -f lavfi -i "color=c=0x00FF00:s=720x480:d=10" \
       -f lavfi -i "sine=frequency=440:duration=10" \
       -c:v libx264 -preset ultrafast -c:a aac \
       -colorspace bt470bg -color_primaries bt470bg -color_trc bt470bg \
       -t 10 "$test_dir/bt601_sd.mp4"|| exit 1

cat > "$test_dir/README.txt" << EOF
Test Case: Color Space Variations

Contains:
- bt709_hd.mp4 - BT.709 color space (HD standard)
- bt601_sd.mp4 - BT.601 color space (SD standard)

Purpose: Test color space detection and conversion.
Incorrect handling causes color shifts.
EOF

echo "  ✓ Created: $test_dir/ (color spaces)"
echo

#------------------------------------------------------------------------------
# Test Case 8: Different pixel formats
#------------------------------------------------------------------------------
echo "Creating: Different pixel formats"
test_dir="$OUTPUT_DIR/pixel_formats"
mkdir -p "$test_dir"

# yuv420p (standard)
ffmpeg -y -f lavfi -i "color=c=blue:s=1280x720:d=10" \
       -c:v libx264 -preset ultrafast -pix_fmt yuv420p \
       -t 10 "$test_dir/yuv420p.mp4"|| exit 1

# yuv422p (higher chroma)
ffmpeg -y -f lavfi -i "color=c=green:s=1280x720:d=10" \
       -c:v libx264 -preset ultrafast -pix_fmt yuv422p \
       -t 10 "$test_dir/yuv422p.mp4"|| exit 1

cat > "$test_dir/README.txt" << EOF
Test Case: Pixel Format Variations

Contains:
- yuv420p.mp4 - Standard 4:2:0 subsampling
- yuv422p.mp4 - Higher quality 4:2:2 subsampling

Purpose: Test pixel format detection and conversion.
Most web players need yuv420p.
EOF

echo "  ✓ Created: $test_dir/ (pixel formats)"
echo

#------------------------------------------------------------------------------
# Test Case 9: Missing creation date metadata
#------------------------------------------------------------------------------
echo "Creating: Missing creation date"
test_dir="$OUTPUT_DIR/missing_creation_date"
mkdir -p "$test_dir"

ffmpeg -y -f lavfi -i "color=c=gray:s=1280x720:d=10" \
       -c:v libx264 -preset ultrafast \
       -fflags +bitexact -map_metadata -1 \
       -t 10 "$test_dir/video.mp4"|| exit 1

cat > "$test_dir/README.txt" << EOF
Test Case: Missing Creation Date

Video has no creation date metadata.

Purpose: Test handling when metadata is absent.
Should use file system timestamps or flag as unknown.
EOF

echo "  ✓ Created: $test_dir/video.mp4"
echo

#------------------------------------------------------------------------------
# Test Case 10: HDR metadata (simulated)
#------------------------------------------------------------------------------
echo "Creating: HDR metadata"
test_dir="$OUTPUT_DIR/hdr_metadata"
mkdir -p "$test_dir"

ffmpeg -y -f lavfi -i "color=c=white:s=1920x1080:d=10" \
       -c:v libx265 -preset ultrafast \
       -x265-params "hdr-opt=1:repeat-headers=1:colorprim=bt2020:transfer=smpte2084:colormatrix=bt2020nc" \
       -t 10 "$test_dir/hdr_video.mp4"|| exit 1

cat > "$test_dir/README.txt" << EOF
Test Case: HDR Metadata

Video with HDR (High Dynamic Range) metadata.
Uses BT.2020 color space and PQ transfer.

Purpose: Test HDR detection and tone mapping for SDR displays.
EOF

echo "  ✓ Created: $test_dir/hdr_video.mp4"
echo

echo "=========================================="
echo "Metadata/encoding issue test cases generated in: $OUTPUT_DIR"
echo "Total test cases: 10"
echo "=========================================="
