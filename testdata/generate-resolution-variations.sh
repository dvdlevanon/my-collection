#!/bin/bash

# generate-resolution-variations.sh
# Generates test cases for various resolutions, aspect ratios, and orientations

SCRIPT_DIR=$(dirname "$(realpath "${BASH_SOURCE[0]}")")

show_usage() {
    cat << EOF
Usage: $0 <output_dir>

Generates resolution/aspect ratio test cases including:
- Standard resolutions (480p, 720p, 1080p, 4K)
- Non-standard resolutions (DVD PAL, old rips)
- Portrait orientation (phone videos)
- Ultra-wide aspect ratios
- Mixed resolutions in same directory
- Wrong aspect ratio metadata

Output structure:
    output_dir/
    ├── 480p_sd/
    ├── 720p_hd/
    ├── 1080p_fullhd/
    ├── 4k_uhd/
    ├── portrait_phone/
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
echo "Generating Resolution Variation Test Cases"
echo "=========================================="
echo

# Test case definitions: name:resolution:description
test_cases=(
    "480p_sd:640x480:Standard Definition 480p"
    "720p_hd:1280x720:HD 720p"
    "1080p_fullhd:1920x1080:Full HD 1080p"
    "4k_uhd:3840x2160:4K Ultra HD"
    "dvd_pal:720x576:DVD PAL Resolution"
    "dvd_ntsc:720x480:DVD NTSC Resolution"
    "old_vcd:352x240:Old VCD Quality"
    "portrait_phone:1080x1920:Vertical Phone Video"
    "portrait_9_16:720x1280:Portrait 9:16"
    "ultrawide_219:2560x1080:Ultra-wide 21:9"
    "square_instagram:1080x1080:Square 1:1 (Instagram)"
    "weird_resolution:1440x1080:Non-standard Anamorphic"
    "tiny_thumbnail:160x120:Tiny Thumbnail Size"
)

for test_case in "${test_cases[@]}"; do
    IFS=':' read -r name resolution description <<< "$test_case"
    
    test_dir="$OUTPUT_DIR/$name"
    mkdir -p "$test_dir"
    
    echo "Creating: $description ($resolution)"
    
    # Generate video at specific resolution
    if [ -f "$SCRIPT_DIR/generate-video.sh" ]; then
        $SCRIPT_DIR/generate-video.sh "$description" 15 "$test_dir/video.mp4" \
            --resolution "$resolution" --video-codec libx264 --audio-codec aac >/dev/null 2>&1
    else
        # Fallback: generate with ffmpeg
        ffmpeg -y -f lavfi -i "color=c=orange:s=$resolution:d=15" \
               -f lavfi -i "sine=frequency=440:duration=15" \
               -c:v libx264 -preset ultrafast -c:a aac \
               -t 15 "$test_dir/video.mp4" >/dev/null 2>&1
    fi
    
    # Create README
    cat > "$test_dir/README.txt" << EOF
Test Case: $description

Resolution: $resolution
Aspect Ratio: $(echo $resolution | awk -F'x' '{printf "%.2f:1", $1/$2}')

Purpose: Test resolution detection and scaling
EOF
    
    echo "  ✓ Created: $test_dir/video.mp4"
    echo
done

#------------------------------------------------------------------------------
# Special Case: Mixed resolutions in same directory
#------------------------------------------------------------------------------
echo "Creating: Mixed resolutions in same directory"
test_dir="$OUTPUT_DIR/mixed_resolutions"
mkdir -p "$test_dir"

# Create multiple videos with different resolutions
resolutions=("640x480" "1280x720" "1920x1080" "3840x2160")
for i in "${!resolutions[@]}"; do
    res="${resolutions[$i]}"
    ffmpeg -y -f lavfi -i "color=c=blue:s=$res:d=10" \
           -f lavfi -i "sine=frequency=440:duration=10" \
           -c:v libx264 -preset ultrafast -c:a aac \
           -t 10 "$test_dir/video_${res/x/_}.mp4" >/dev/null 2>&1
    echo "  ✓ Created: video_${res/x/_}.mp4 ($res)"
done

cat > "$test_dir/README.txt" << EOF
Test Case: Mixed Resolutions

Contains multiple videos with different resolutions in the same directory:
- 640x480 (SD)
- 1280x720 (HD)
- 1920x1080 (Full HD)  
- 3840x2160 (4K)

Purpose: Test handling of mixed resolution libraries
EOF

echo

#------------------------------------------------------------------------------
# Special Case: Wrong aspect ratio metadata
#------------------------------------------------------------------------------
echo "Creating: Wrong aspect ratio metadata"
test_dir="$OUTPUT_DIR/wrong_aspect_ratio"
mkdir -p "$test_dir"

# Create 16:9 video but tag it as 4:3
ffmpeg -y -f lavfi -i "color=c=red:s=1920x1080:d=15" \
       -f lavfi -i "sine=frequency=440:duration=15" \
       -c:v libx264 -preset ultrafast -c:a aac \
       -aspect 4:3 \
       -t 15 "$test_dir/video.mp4" >/dev/null 2>&1

cat > "$test_dir/README.txt" << EOF
Test Case: Wrong Aspect Ratio Metadata

Video is actually 16:9 (1920x1080) but metadata claims 4:3.

Purpose: Test handling of incorrect aspect ratio metadata.
This causes stretching/squashing in some players.
EOF

echo "  ✓ Created: $test_dir/video.mp4"
echo

#------------------------------------------------------------------------------
# Special Case: Rotation metadata (phone videos)
#------------------------------------------------------------------------------
echo "Creating: Rotation metadata test"
test_dir="$OUTPUT_DIR/rotation_metadata"
mkdir -p "$test_dir"

# Create video with rotation metadata
ffmpeg -y -f lavfi -i "color=c=green:s=1080x1920:d=15" \
       -f lavfi -i "sine=frequency=440:duration=15" \
       -c:v libx264 -preset ultrafast -c:a aac \
       -metadata:s:v:0 rotate=90 \
       -t 15 "$test_dir/video.mp4" >/dev/null 2>&1

cat > "$test_dir/README.txt" << EOF
Test Case: Rotation Metadata

Portrait video (1080x1920) with 90-degree rotation metadata.
Common in phone recordings.

Purpose: Test proper handling of rotation metadata.
Some players ignore this and display sideways video.
EOF

echo "  ✓ Created: $test_dir/video.mp4"
echo

echo "=========================================="
echo "Resolution variation test cases generated in: $OUTPUT_DIR"
echo "Total test cases: $((${#test_cases[@]} + 3))"
echo "=========================================="
