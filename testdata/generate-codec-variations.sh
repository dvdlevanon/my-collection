#!/bin/bash

# generate-codec-variations.sh
# Generates test cases for codec detection and compatibility testing

SCRIPT_DIR=$(dirname "$(realpath "${BASH_SOURCE[0]}")")

show_usage() {
    cat << EOF
Usage: $0 <output_dir>

Generates various codec combination test cases including:
- Old codecs (DivX, XviD, MPEG2)
- Modern codecs (H.264, H.265, VP9, AV1)
- Various audio codecs (MP3, AAC, AC3, DTS, Opus)
- Different containers (MKV, AVI, MP4, MOV, WebM)
- Problematic combinations

Output structure:
    output_dir/
    ├── old_divx_avi/
    ├── xvid_mp3_avi/
    ├── h264_ac3_mkv/
    ├── h265_dts_mkv/
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
echo "Generating Codec Variation Test Cases"
echo "=========================================="
echo

# Test case definitions: name:container:vcodec:acodec:description
test_cases=(
    "old_divx_avi:avi:mpeg4:libmp3lame:Old DivX/XviD in AVI"
    "xvid_mp3_avi:avi:mpeg4:libmp3lame:XviD with MP3 audio"
    "h264_aac_mp4:mp4:libx264:aac:Standard H.264/AAC MP4"
    "h264_ac3_mkv:mkv:libx264:ac3:H.264 with AC3 audio"
    "h264_dts_mkv:mkv:libx264:dts:H.264 with DTS audio"
    "h265_aac_mp4:mp4:libx265:aac:H.265/HEVC with AAC"
    "h265_ac3_mkv:mkv:libx265:ac3:H.265 with AC3 surround"
    "vp9_opus_webm:webm:libvpx-vp9:libopus:VP9/Opus WebM"
    "av1_opus_mp4:mp4:libsvtav1:libopus:AV1 with Opus audio"
    "h264_no_audio:mp4:libx264::H.264 video only - no audio"
    "mpeg2_ac3_mpg:mpg:mpeg2video:ac3:MPEG2 (DVD-like)"
)

for test_case in "${test_cases[@]}"; do
    IFS=':' read -r name container vcodec acodec description <<< "$test_case"
    
    test_dir="$OUTPUT_DIR/$name"
    mkdir -p "$test_dir"
    
    echo "Creating: $description"
    echo "  Directory: $test_dir"
    
    # Generate video
    video_file="$test_dir/temp_video.mp4"
    if [ -f "$SCRIPT_DIR/generate-video.sh" ]; then
        $SCRIPT_DIR/generate-video.sh "$description" 15 "$video_file" \
            --video-codec libx264 --audio-codec aac --no-audio >/dev/null 2>&1
    else
        # Fallback: create with ffmpeg directly
        ffmpeg -y -f lavfi -i color=c=blue:s=1280x720:d=15 \
               -c:v libx264 -preset ultrafast -t 15 "$video_file" >/dev/null 2>&1
    fi
    
    # Generate audio if needed
    if [ -n "$acodec" ]; then
        audio_file="$test_dir/temp_audio.${acodec##*lib}"
        if [ -f "$SCRIPT_DIR/generate-audio.sh" ]; then
            $SCRIPT_DIR/generate-audio.sh 15 "$audio_file" \
                --codec "$acodec" --channels 2 >/dev/null 2>&1
        else
            # Fallback
            if [ -f "$SCRIPT_DIR/audio.mp3" ]; then
                ffmpeg -y -i "$SCRIPT_DIR/audio.mp3" -c:a "$acodec" -t 15 "$audio_file" >/dev/null 2>&1
            else
                ffmpeg -y -f lavfi -i "sine=frequency=440:duration=15" \
                       -c:a "$acodec" "$audio_file" >/dev/null 2>&1
            fi
        fi
    fi
    
    # Remux to target container with target codecs
    output_file="$test_dir/video.$container"
    
    if [ -n "$acodec" ] && [ -f "$audio_file" ]; then
        # Video + Audio
        ffmpeg -y -i "$video_file" -i "$audio_file" \
               -map 0:v -map 1:a \
               -c:v "$vcodec" -preset ultrafast \
               -c:a "$acodec" \
               -t 15 "$output_file" >/dev/null 2>&1
    else
        # Video only
        ffmpeg -y -i "$video_file" \
               -c:v "$vcodec" -preset ultrafast \
               -t 15 "$output_file" >/dev/null 2>&1
    fi
    
    # Cleanup temp files
    rm -f "$video_file" "$audio_file"
    
    # Create metadata file
    cat > "$test_dir/README.txt" << EOF
Test Case: $description

Video Codec: $vcodec
Audio Codec: ${acodec:-none}
Container: $container

Purpose: Test codec detection and compatibility
EOF
    
    echo "  ✓ Created: $output_file"
    echo
done

echo "=========================================="
echo "Codec variations generated in: $OUTPUT_DIR"
echo "Total test cases: ${#test_cases[@]}"
echo "=========================================="
