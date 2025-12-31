#!/bin/bash

# generate-subtitle-variations.sh
# Generates test cases for subtitle handling: external, embedded, multiple languages, encodings

SCRIPT_DIR=$(dirname "$(realpath "${BASH_SOURCE[0]}")")

show_usage() {
    cat << EOF
Usage: $0 <output_dir>

Generates subtitle test cases including:
- External SRT files
- Embedded subtitles in MKV
- Multiple subtitle tracks
- Different languages (English, Hebrew)
- Forced subtitles
- Different encodings (UTF-8, Windows-1255)
- ASS/SSA styled subtitles

Output structure:
    output_dir/
    ├── external_srt/
    ├── embedded_subtitles/
    ├── multi_language_subs/
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
echo "Generating Subtitle Variation Test Cases"
echo "=========================================="
echo

# Helper function to generate base video
generate_base_video() {
    local output="$1"
    local text="$2"
    
    if [ -f "$SCRIPT_DIR/generate-video.sh" ]; then
        $SCRIPT_DIR/generate-video.sh "$text" 15 "$output" \
            --video-codec libx264 --audio-codec aac >/dev/null 2>&1
    else
        ffmpeg -y -f lavfi -i color=c=purple:s=1280x720:d=15 \
               -f lavfi -i sine=frequency=440:duration=15 \
               -c:v libx264 -preset ultrafast -c:a aac -t 15 "$output" >/dev/null 2>&1
    fi
}

# Helper function to generate subtitle
generate_subtitle() {
    local output="$1"
    local text="$2"
    local duration="${3:-15}"
    
    if [ -f "$SCRIPT_DIR/generate-subtitle.sh" ]; then
        $SCRIPT_DIR/generate-subtitle.sh "$duration" "$text" "$output" >/dev/null 2>&1
    else
        # Fallback: simple SRT
        cat > "$output" << EOF
1
00:00:00,000 --> 00:00:05,000
$text - 1

2
00:00:05,000 --> 00:00:10,000
$text - 2

3
00:00:10,000 --> 00:00:15,000
$text - 3

EOF
    fi
}

#------------------------------------------------------------------------------
# Test Case 1: External SRT subtitles
#------------------------------------------------------------------------------
echo "Creating: External SRT subtitles"
test_dir="$OUTPUT_DIR/external_srt"
mkdir -p "$test_dir"

generate_base_video "$test_dir/movie.mp4" "External SRT Test"
generate_subtitle "$test_dir/movie.srt" "External English Subtitle"
generate_subtitle "$test_dir/movie.heb.srt" "כתוביות בעברית"

cat > "$test_dir/README.txt" << EOF
Test Case: External SRT Subtitles

Contains:
- movie.mp4 (video with audio)
- movie.srt (English subtitles)
- movie.heb.srt (Hebrew subtitles)

Purpose: Test external subtitle detection and loading
EOF

echo "  ✓ Created: $test_dir/"
echo

#------------------------------------------------------------------------------
# Test Case 2: Embedded subtitles in MKV
#------------------------------------------------------------------------------
echo "Creating: Embedded subtitles in MKV"
test_dir="$OUTPUT_DIR/embedded_subtitles"
mkdir -p "$test_dir"

video_temp="$test_dir/temp_video.mp4"
sub1="$test_dir/temp_eng.srt"
sub2="$test_dir/temp_heb.srt"

generate_base_video "$video_temp" "Embedded Subs Test"
generate_subtitle "$sub1" "Embedded English"
generate_subtitle "$sub2" "כתוביות משובצות"

if [ -f "$SCRIPT_DIR/mux-container.sh" ]; then
    $SCRIPT_DIR/mux-container.sh --video "$video_temp" \
        --subtitle "$sub1" --sub-lang eng --sub-title "English" \
        --subtitle "$sub2" --sub-lang heb --sub-title "Hebrew" \
        --output "$test_dir/movie.mkv" --copy-codecs >/dev/null 2>&1
else
    ffmpeg -y -i "$video_temp" -i "$sub1" -i "$sub2" \
           -map 0 -map 1 -map 2 \
           -c copy -c:s srt \
           -metadata:s:s:0 language=eng -metadata:s:s:0 title="English" \
           -metadata:s:s:1 language=heb -metadata:s:s:1 title="Hebrew" \
           "$test_dir/movie.mkv" >/dev/null 2>&1
fi

rm -f "$video_temp" "$sub1" "$sub2"
echo "  ✓ Created: $test_dir/movie.mkv"
echo

#------------------------------------------------------------------------------
# Test Case 3: Multiple subtitle tracks (3+ languages)
#------------------------------------------------------------------------------
echo "Creating: Multiple subtitle tracks (English, Hebrew, Forced)"
test_dir="$OUTPUT_DIR/multi_language_subs"
mkdir -p "$test_dir"

video_temp="$test_dir/temp_video.mp4"
sub_eng="$test_dir/temp_eng.srt"
sub_heb="$test_dir/temp_heb.srt"
sub_forced="$test_dir/temp_forced.srt"

generate_base_video "$video_temp" "Multi-Sub Test"
generate_subtitle "$sub_eng" "English Full"
generate_subtitle "$sub_heb" "עברית מלאה"

# Create forced subtitle (shorter, only important parts)
cat > "$sub_forced" << EOF
1
00:00:03,000 --> 00:00:06,000
[Important Text]

2
00:00:12,000 --> 00:00:14,000
[Critical Info]

EOF

ffmpeg -y -i "$video_temp" -i "$sub_eng" -i "$sub_heb" -i "$sub_forced" \
       -map 0 -map 1 -map 2 -map 3 \
       -c copy -c:s srt \
       -metadata:s:s:0 language=eng -metadata:s:s:0 title="English" \
       -metadata:s:s:1 language=heb -metadata:s:s:1 title="Hebrew" \
       -metadata:s:s:2 language=eng -metadata:s:s:2 title="English Forced" \
       -disposition:s:2 forced \
       "$test_dir/movie.mkv" >/dev/null 2>&1

rm -f "$video_temp" "$sub_eng" "$sub_heb" "$sub_forced"
echo "  ✓ Created: $test_dir/movie.mkv"
echo

#------------------------------------------------------------------------------
# Test Case 4: Both external and embedded subtitles
#------------------------------------------------------------------------------
echo "Creating: Both external SRT and embedded subtitles"
test_dir="$OUTPUT_DIR/external_and_embedded"
mkdir -p "$test_dir"

video_temp="$test_dir/temp_video.mp4"
sub_embedded="$test_dir/temp_embedded.srt"

generate_base_video "$video_temp" "Mixed Subs Test"
generate_subtitle "$sub_embedded" "Embedded Subtitle"

# Create MKV with embedded subtitle
ffmpeg -y -i "$video_temp" -i "$sub_embedded" \
       -map 0 -map 1 \
       -c copy -c:s srt \
       -metadata:s:s:0 language=eng \
       "$test_dir/movie.mkv" >/dev/null 2>&1

# Also create external subtitle
generate_subtitle "$test_dir/movie.srt" "External Subtitle"

rm -f "$video_temp" "$sub_embedded"
echo "  ✓ Created: $test_dir/movie.mkv + movie.srt"
echo

#------------------------------------------------------------------------------
# Test Case 5: ASS/SSA styled subtitles
#------------------------------------------------------------------------------
echo "Creating: ASS styled subtitles"
test_dir="$OUTPUT_DIR/ass_styled_subs"
mkdir -p "$test_dir"

generate_base_video "$test_dir/movie.mp4" "ASS Subtitle Test"

# Create ASS subtitle with styling
cat > "$test_dir/movie.ass" << 'EOF'
[Script Info]
Title: Test ASS Subtitle
ScriptType: v4.00+

[V4+ Styles]
Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding
Style: Default,Arial,20,&H00FFFFFF,&H000000FF,&H00000000,&H00000000,0,0,0,0,100,100,0,0,1,2,0,2,10,10,10,1

[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
Dialogue: 0,0:00:00.00,0:00:05.00,Default,,0,0,0,,Styled subtitle line 1
Dialogue: 0,0:00:05.00,0:00:10.00,Default,,0,0,0,,{\i1}Italic text{\i0} in subtitle
Dialogue: 0,0:00:10.00,0:00:15.00,Default,,0,0,0,,{\c&HFF0000&}Blue colored text{\c}
EOF

echo "  ✓ Created: $test_dir/movie.mp4 + movie.ass"
echo

#------------------------------------------------------------------------------
# Test Case 6: Hebrew encoding test (UTF-8 vs Windows-1255)
#------------------------------------------------------------------------------
echo "Creating: Hebrew encoding test (UTF-8 and Windows-1255)"
test_dir="$OUTPUT_DIR/hebrew_encoding"
mkdir -p "$test_dir"

generate_base_video "$test_dir/movie.mp4" "Hebrew Encoding Test"

# UTF-8 (correct)
cat > "$test_dir/movie.utf8.srt" << EOF
1
00:00:00,000 --> 00:00:05,000
כתוביות בעברית - UTF-8

2
00:00:05,000 --> 00:00:10,000
שורה שנייה בעברית

3
00:00:10,000 --> 00:00:15,000
שלום עולם

EOF

# Note: We can't easily generate Windows-1255 in bash, so just document it
cat > "$test_dir/README.txt" << EOF
Test Case: Hebrew Encoding

Contains:
- movie.utf8.srt (UTF-8 encoded - correct)

Purpose: Test Hebrew subtitle encoding detection
Note: In real scenarios, you might encounter Windows-1255 encoded files
that display incorrectly if not detected properly.
EOF

echo "  ✓ Created: $test_dir/"
echo

#------------------------------------------------------------------------------
# Test Case 7: No subtitles at all
#------------------------------------------------------------------------------
echo "Creating: No subtitles"
test_dir="$OUTPUT_DIR/no_subtitles"
mkdir -p "$test_dir"

generate_base_video "$test_dir/movie.mp4" "No Subtitles Test"

cat > "$test_dir/README.txt" << EOF
Test Case: No Subtitles

Contains only video with audio, no subtitle files.

Purpose: Test handling of videos without subtitles
EOF

echo "  ✓ Created: $test_dir/"
echo

echo "=========================================="
echo "Subtitle variation test cases generated in: $OUTPUT_DIR"
echo "Total test cases: 7"
echo "=========================================="
