#!/bin/bash

# generate-audio-tracks.sh
# Generates test cases for multiple audio tracks, different channels, languages

SCRIPT_DIR=$(dirname "$(realpath "${BASH_SOURCE[0]}")")

show_usage() {
    cat << EOF
Usage: $0 <output_dir>

Generates multi-audio track test cases including:
- Multiple languages (English, Hebrew)
- Different channel configurations (stereo, 5.1, 7.1)
- Mixed audio codecs in same file
- Audio sync issues
- Missing/corrupt audio tracks

Output structure:
    output_dir/
    ├── dual_audio_eng_heb/
    ├── stereo_vs_51_surround/
    ├── three_audio_tracks/
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
echo "Generating Multi-Audio Track Test Cases"
echo "=========================================="
echo

# Helper function to generate base video
generate_base_video() {
    local output="$1"
    local text="$2"
    
    if [ -f "$SCRIPT_DIR/generate-video.sh" ]; then
        $SCRIPT_DIR/generate-video.sh "$text" 15 "$output" \
            --video-codec libx264 --no-audio >/dev/null 2>&1
    else
        ffmpeg -y -f lavfi -i color=c=green:s=1280x720:d=15 \
               -c:v libx264 -preset ultrafast -t 15 "$output" >/dev/null 2>&1
    fi
}

# Helper function to generate audio
generate_audio_track() {
    local output="$1"
    local codec="$2"
    local channels="$3"
    
    if [ -f "$SCRIPT_DIR/generate-audio.sh" ]; then
        $SCRIPT_DIR/generate-audio.sh 15 "$output" \
            --codec "$codec" --channels "$channels" >/dev/null 2>&1
    else
        if [ -f "$SCRIPT_DIR/audio.mp3" ]; then
            ffmpeg -y -i "$SCRIPT_DIR/audio.mp3" -c:a "$codec" -ac "$channels" -t 15 "$output" >/dev/null 2>&1
        else
            ffmpeg -y -f lavfi -i "sine=frequency=440:duration=15" \
                   -c:a "$codec" -ac "$channels" "$output" >/dev/null 2>&1
        fi
    fi
}

#------------------------------------------------------------------------------
# Test Case 1: Dual audio (English + Hebrew)
#------------------------------------------------------------------------------
echo "Creating: Dual audio tracks (English AAC + Hebrew AC3)"
test_dir="$OUTPUT_DIR/dual_audio_eng_heb"
mkdir -p "$test_dir"

video_temp="$test_dir/temp_video.mp4"
audio1="$test_dir/temp_audio_eng.aac"
audio2="$test_dir/temp_audio_heb.ac3"

generate_base_video "$video_temp" "Dual Audio Test"
generate_audio_track "$audio1" "aac" "2"
generate_audio_track "$audio2" "ac3" "2"

if [ -f "$SCRIPT_DIR/mux-container.sh" ]; then
    $SCRIPT_DIR/mux-container.sh --video "$video_temp" \
        --audio "$audio1" --audio-lang eng --audio-title "English" \
        --audio "$audio2" --audio-lang heb --audio-title "Hebrew" \
        --output "$test_dir/movie.mkv" --copy-codecs >/dev/null 2>&1
else
    ffmpeg -y -i "$video_temp" -i "$audio1" -i "$audio2" \
           -map 0:v -map 1:a -map 2:a \
           -c copy \
           -metadata:s:a:0 language=eng -metadata:s:a:0 title="English" \
           -metadata:s:a:1 language=heb -metadata:s:a:1 title="Hebrew" \
           "$test_dir/movie.mkv" >/dev/null 2>&1
fi

rm -f "$video_temp" "$audio1" "$audio2"
echo "  ✓ Created: $test_dir/movie.mkv"
echo

#------------------------------------------------------------------------------
# Test Case 2: Stereo vs 5.1 surround
#------------------------------------------------------------------------------
echo "Creating: Stereo vs 5.1 surround audio"
test_dir="$OUTPUT_DIR/stereo_vs_51_surround"
mkdir -p "$test_dir"

video_temp="$test_dir/temp_video.mp4"
audio_stereo="$test_dir/temp_stereo.aac"
audio_51="$test_dir/temp_51.ac3"

generate_base_video "$video_temp" "Surround Sound Test"
generate_audio_track "$audio_stereo" "aac" "2"
generate_audio_track "$audio_51" "ac3" "6"

ffmpeg -y -i "$video_temp" -i "$audio_stereo" -i "$audio_51" \
       -map 0:v -map 1:a -map 2:a \
       -c copy \
       -metadata:s:a:0 title="Stereo" \
       -metadata:s:a:1 title="5.1 Surround" \
       "$test_dir/movie.mkv" >/dev/null 2>&1

rm -f "$video_temp" "$audio_stereo" "$audio_51"
echo "  ✓ Created: $test_dir/movie.mkv"
echo

#------------------------------------------------------------------------------
# Test Case 3: Three audio tracks (multilingual)
#------------------------------------------------------------------------------
echo "Creating: Three audio tracks (English, Hebrew, Commentary)"
test_dir="$OUTPUT_DIR/three_audio_tracks"
mkdir -p "$test_dir"

video_temp="$test_dir/temp_video.mp4"
audio1="$test_dir/temp_eng.aac"
audio2="$test_dir/temp_heb.ac3"
audio3="$test_dir/temp_commentary.mp3"

generate_base_video "$video_temp" "Multi-Language Test"
generate_audio_track "$audio1" "aac" "2"
generate_audio_track "$audio2" "ac3" "6"
generate_audio_track "$audio3" "libmp3lame" "2"

ffmpeg -y -i "$video_temp" -i "$audio1" -i "$audio2" -i "$audio3" \
       -map 0:v -map 1:a -map 2:a -map 3:a \
       -c copy \
       -metadata:s:a:0 language=eng -metadata:s:a:0 title="English Stereo" \
       -metadata:s:a:1 language=heb -metadata:s:a:1 title="Hebrew 5.1" \
       -metadata:s:a:2 language=eng -metadata:s:a:2 title="Director Commentary" \
       "$test_dir/movie.mkv" >/dev/null 2>&1

rm -f "$video_temp" "$audio1" "$audio2" "$audio3"
echo "  ✓ Created: $test_dir/movie.mkv"
echo

#------------------------------------------------------------------------------
# Test Case 4: Mixed codecs (AAC, AC3, DTS in one file)
#------------------------------------------------------------------------------
echo "Creating: Mixed audio codecs (AAC + AC3 + DTS)"
test_dir="$OUTPUT_DIR/mixed_audio_codecs"
mkdir -p "$test_dir"

video_temp="$test_dir/temp_video.mp4"
audio_aac="$test_dir/temp.aac"
audio_ac3="$test_dir/temp.ac3"
audio_dts="$test_dir/temp.dts"

generate_base_video "$video_temp" "Mixed Codec Test"
generate_audio_track "$audio_aac" "aac" "2"
generate_audio_track "$audio_ac3" "ac3" "6"
generate_audio_track "$audio_dts" "dts" "6"

ffmpeg -y -i "$video_temp" -i "$audio_aac" -i "$audio_ac3" -i "$audio_dts" \
       -map 0:v -map 1:a -map 2:a -map 3:a \
       -c copy \
       -metadata:s:a:0 title="AAC Stereo" \
       -metadata:s:a:1 title="AC3 5.1" \
       -metadata:s:a:2 title="DTS 5.1" \
       "$test_dir/movie.mkv" >/dev/null 2>&1

rm -f "$video_temp" "$audio_aac" "$audio_ac3" "$audio_dts"
echo "  ✓ Created: $test_dir/movie.mkv"
echo

#------------------------------------------------------------------------------
# Test Case 5: 7.1 surround
#------------------------------------------------------------------------------
echo "Creating: 7.1 surround sound"
test_dir="$OUTPUT_DIR/71_surround"
mkdir -p "$test_dir"

video_temp="$test_dir/temp_video.mp4"
audio_71="$test_dir/temp_71.ac3"

generate_base_video "$video_temp" "7.1 Surround Test"
generate_audio_track "$audio_71" "ac3" "8"

ffmpeg -y -i "$video_temp" -i "$audio_71" \
       -map 0:v -map 1:a \
       -c copy \
       -metadata:s:a:0 title="7.1 Surround" \
       "$test_dir/movie.mkv" >/dev/null 2>&1

rm -f "$video_temp" "$audio_71"
echo "  ✓ Created: $test_dir/movie.mkv"
echo

#------------------------------------------------------------------------------
# Test Case 6: Video only - no audio at all
#------------------------------------------------------------------------------
echo "Creating: Video only (no audio track)"
test_dir="$OUTPUT_DIR/no_audio_track"
mkdir -p "$test_dir"

ffmpeg -y -f lavfi -i color=c=red:s=1280x720:d=15 \
       -c:v libx264 -preset ultrafast -t 15 \
       "$test_dir/video.mp4" >/dev/null 2>&1

echo "  ✓ Created: $test_dir/video.mp4"
echo

echo "=========================================="
echo "Multi-audio track test cases generated in: $OUTPUT_DIR"
echo "Total test cases: 6"
echo "=========================================="
