#!/bin/bash

# generate-container.sh
# Mux video, audio, and subtitle streams into various container formats
# Supports multiple audio tracks and subtitle tracks

show_usage() {
    cat << EOF
Usage: $0 --video FILE --output FILE [options]

Required:
    --video FILE        Input video file
    --output FILE       Output container file

Optional:
    --audio FILE        Add audio track (can be specified multiple times)
    --audio-lang CODE   Language code for last audio track (eng, heb, etc.)
    --audio-title TEXT  Title for last audio track
    --subtitle FILE     Add subtitle track (can be specified multiple times)
    --sub-lang CODE     Language code for last subtitle track
    --sub-title TEXT    Title for last subtitle track
    --sub-forced        Mark last subtitle as forced
    --container FORMAT  Force container format (mkv, mp4, avi, mov, webm)
    --copy-codecs       Copy all codecs without re-encoding (faster)
    --help, -h          Show this help

Examples:
    # Simple mux
    $0 --video video.mp4 --audio audio.mp3 --output movie.mkv

    # Multiple audio tracks with metadata
    $0 --video video.mp4 \\
       --audio audio_eng.aac --audio-lang eng --audio-title "English" \\
       --audio audio_heb.ac3 --audio-lang heb --audio-title "Hebrew" \\
       --output movie.mkv

    # With subtitles
    $0 --video video.mp4 --audio audio.aac \\
       --subtitle subs_eng.srt --sub-lang eng \\
       --subtitle subs_heb.srt --sub-lang heb --sub-forced \\
       --output movie.mkv

    # Quick copy without re-encoding
    $0 --video video.mp4 --audio audio.aac --subtitle subs.srt \\
       --output movie.mkv --copy-codecs
EOF
}

# Arrays for tracks
declare -a VIDEO_FILES
declare -a AUDIO_FILES
declare -a AUDIO_LANGS
declare -a AUDIO_TITLES
declare -a SUBTITLE_FILES
declare -a SUB_LANGS
declare -a SUB_TITLES
declare -a SUB_FORCED

# Variables
OUTPUT_FILE=""
CONTAINER=""
COPY_CODECS=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --video)
            VIDEO_FILES+=("$2")
            shift 2
            ;;
        --audio)
            AUDIO_FILES+=("$2")
            AUDIO_LANGS+=("")
            AUDIO_TITLES+=("")
            shift 2
            ;;
        --audio-lang)
            # Set language for last audio track
            idx=$((${#AUDIO_LANGS[@]} - 1))
            AUDIO_LANGS[$idx]="$2"
            shift 2
            ;;
        --audio-title)
            idx=$((${#AUDIO_TITLES[@]} - 1))
            AUDIO_TITLES[$idx]="$2"
            shift 2
            ;;
        --subtitle)
            SUBTITLE_FILES+=("$2")
            SUB_LANGS+=("")
            SUB_TITLES+=("")
            SUB_FORCED+=("0")
            shift 2
            ;;
        --sub-lang)
            idx=$((${#SUB_LANGS[@]} - 1))
            SUB_LANGS[$idx]="$2"
            shift 2
            ;;
        --sub-title)
            idx=$((${#SUB_TITLES[@]} - 1))
            SUB_TITLES[$idx]="$2"
            shift 2
            ;;
        --sub-forced)
            idx=$((${#SUB_FORCED[@]} - 1))
            SUB_FORCED[$idx]="1"
            shift
            ;;
        --output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        --container)
            CONTAINER="$2"
            shift 2
            ;;
        --copy-codecs)
            COPY_CODECS=true
            shift
            ;;
        --help|-h)
            show_usage
            exit 0
            ;;
        *)
            echo "Error: Unknown option $1"
            show_usage
            exit 1
            ;;
    esac
done

# Validate inputs
if [ ${#VIDEO_FILES[@]} -eq 0 ]; then
    echo "Error: No video file specified"
    show_usage
    exit 1
fi

if [ -z "$OUTPUT_FILE" ]; then
    echo "Error: No output file specified"
    show_usage
    exit 1
fi

# Check video file exists
if [ ! -f "${VIDEO_FILES[0]}" ]; then
    echo "Error: Video file not found: ${VIDEO_FILES[0]}"
    exit 1
fi

# Check audio files exist
for audio in "${AUDIO_FILES[@]}"; do
    if [ ! -f "$audio" ]; then
        echo "Error: Audio file not found: $audio"
        exit 1
    fi
done

# Check subtitle files exist
for sub in "${SUBTITLE_FILES[@]}"; do
    if [ ! -f "$sub" ]; then
        echo "Error: Subtitle file not found: $sub"
        exit 1
    fi
done

# Auto-detect container from output extension if not specified
if [ -z "$CONTAINER" ]; then
    EXT="${OUTPUT_FILE##*.}"
    CONTAINER="${EXT,,}"
fi

echo "Muxing into container: $CONTAINER"
echo "Video: ${VIDEO_FILES[0]}"
echo "Audio tracks: ${#AUDIO_FILES[@]}"
echo "Subtitle tracks: ${#SUBTITLE_FILES[@]}"
echo

# Build ffmpeg command
FFMPEG_CMD="ffmpeg -y"

# Add video input
FFMPEG_CMD="$FFMPEG_CMD -i \"${VIDEO_FILES[0]}\""

# Add audio inputs
for audio in "${AUDIO_FILES[@]}"; do
    FFMPEG_CMD="$FFMPEG_CMD -i \"$audio\""
done

# Add subtitle inputs
for sub in "${SUBTITLE_FILES[@]}"; do
    FFMPEG_CMD="$FFMPEG_CMD -i \"$sub\""
done

# Build mapping
FFMPEG_CMD="$FFMPEG_CMD -map 0:v"

# Map audio tracks
for i in "${!AUDIO_FILES[@]}"; do
    input_idx=$((i + 1))
    FFMPEG_CMD="$FFMPEG_CMD -map $input_idx:a"
done

# Map subtitle tracks
for i in "${!SUBTITLE_FILES[@]}"; do
    input_idx=$((${#AUDIO_FILES[@]} + i + 1))
    FFMPEG_CMD="$FFMPEG_CMD -map $input_idx:s"
done

# Codec options
if [ "$COPY_CODECS" = true ]; then
    FFMPEG_CMD="$FFMPEG_CMD -c copy"
else
    FFMPEG_CMD="$FFMPEG_CMD -c:v copy"
    
    # Copy audio by default, but some containers might need transcoding
    case "$CONTAINER" in
        mp4|m4v)
            # MP4 doesn't support all codecs well
            FFMPEG_CMD="$FFMPEG_CMD -c:a aac"
            ;;
        *)
            FFMPEG_CMD="$FFMPEG_CMD -c:a copy"
            ;;
    esac
    
    FFMPEG_CMD="$FFMPEG_CMD -c:s copy"
fi

# Add metadata for audio tracks
for i in "${!AUDIO_FILES[@]}"; do
    if [ -n "${AUDIO_LANGS[$i]}" ]; then
        FFMPEG_CMD="$FFMPEG_CMD -metadata:s:a:$i language=${AUDIO_LANGS[$i]}"
    fi
    if [ -n "${AUDIO_TITLES[$i]}" ]; then
        FFMPEG_CMD="$FFMPEG_CMD -metadata:s:a:$i title=\"${AUDIO_TITLES[$i]}\""
    fi
done

# Add metadata for subtitle tracks
for i in "${!SUBTITLE_FILES[@]}"; do
    if [ -n "${SUB_LANGS[$i]}" ]; then
        FFMPEG_CMD="$FFMPEG_CMD -metadata:s:s:$i language=${SUB_LANGS[$i]}"
    fi
    if [ -n "${SUB_TITLES[$i]}" ]; then
        FFMPEG_CMD="$FFMPEG_CMD -metadata:s:s:$i title=\"${SUB_TITLES[$i]}\""
    fi
    if [ "${SUB_FORCED[$i]}" = "1" ]; then
        FFMPEG_CMD="$FFMPEG_CMD -disposition:s:$i forced"
    fi
done

# Container-specific options
case "$CONTAINER" in
    mkv|webm)
        # MKV/WebM support everything
        ;;
    mp4|m4v)
        # MP4 has limitations
        FFMPEG_CMD="$FFMPEG_CMD -movflags +faststart"
        ;;
    avi)
        # AVI is limited
        ;;
    mov)
        # MOV (QuickTime)
        FFMPEG_CMD="$FFMPEG_CMD -movflags +faststart"
        ;;
esac

FFMPEG_CMD="$FFMPEG_CMD \"$OUTPUT_FILE\""

# Execute
echo "Command: $FFMPEG_CMD"
echo

if eval "$FFMPEG_CMD" 2>&1; then
    echo
    echo "✓ Container created: $OUTPUT_FILE"
    
    if command -v ffprobe|| exit 1; then
        echo
        echo "Container info:"
        ffprobe -v quiet -show_entries stream=index,codec_type,codec_name:stream_tags=language,title \
                -of compact "$OUTPUT_FILE" 2>/dev/null | while IFS='|' read -r line; do
            echo "  $line"
        done
    fi
else
    echo "✗ Error creating container"
    exit 1
fi
