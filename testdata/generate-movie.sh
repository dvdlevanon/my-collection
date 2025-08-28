#!/bin/bash

# Script to generate a video with text overlay and audio
# Usage: ./generate-video.sh "text" <length_in_seconds> output.mp4 [audio_file] [options]

SCRIPT_DIR=$(dirname "$(realpath "${BASH_SOURCE[0]}")")

# Default values
DEFAULT_AUDIO="$SCRIPT_DIR/audio.mp3"
DEFAULT_RESOLUTION="1920x1080"
DEFAULT_BG_COLOR="black"
DEFAULT_COUNTER_FONTSIZE=72
DEFAULT_TEXT_FONTSIZE=36
DEFAULT_COUNTER_COLOR="white"
DEFAULT_TEXT_COLOR="yellow"

# Function to show usage
show_usage() {
    cat << EOF
Usage: $0 "text" <length_in_seconds> output.mp4 [audio_file] [options]

Arguments:
    text                Text to display (required)
    length_in_seconds   Video duration in seconds (required)
    output.mp4          Output video file (required)
    audio_file          Audio file to use (optional, default: $DEFAULT_AUDIO)

Options:
    --resolution WxH    Video resolution (default: $DEFAULT_RESOLUTION)
    --bg-color COLOR    Background color (default: $DEFAULT_BG_COLOR)
    --counter-size SIZE Counter font size (default: $DEFAULT_COUNTER_FONTSIZE)
    --text-size SIZE    Text font size (default: $DEFAULT_TEXT_FONTSIZE)
    --counter-color COLOR Counter color (default: $DEFAULT_COUNTER_COLOR)
    --text-color COLOR  Text color (default: $DEFAULT_TEXT_COLOR)
    --no-audio         Generate video without audio
    --help, -h         Show this help message

Examples:
    $0 "Hello World" 30 output.mp4
    $0 "Test Video" 60 test.mp4 custom_audio.wav --resolution 1280x720
    $0 "Silent Video" 15 silent.mp4 --no-audio --bg-color blue
EOF
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to validate file exists
check_file_exists() {
    if [ ! -f "$1" ]; then
        echo "Error: File '$1' not found"
        return 1
    fi
    return 0
}

# Function to validate positive integer
validate_positive_int() {
    if ! [[ "$1" =~ ^[0-9]+$ ]] || [ "$1" -le 0 ]; then
        echo "Error: '$1' must be a positive integer"
        return 1
    fi
    return 0
}

# Function to validate resolution format
validate_resolution() {
    if ! [[ "$1" =~ ^[0-9]+x[0-9]+$ ]]; then
        echo "Error: Resolution must be in format WIDTHxHEIGHT (e.g., 1920x1080)"
        return 1
    fi
    return 0
}

# Parse arguments
if [ $# -lt 3 ]; then
    echo "Error: Missing required arguments"
    show_usage
    exit 1
fi

# Handle help flag
if [[ "$1" == "--help" || "$1" == "-h" ]]; then
    show_usage
    exit 0
fi

# Required arguments
TEXT="$1"
LENGTH="$2"
OUTPUT_FILE="$3"
shift 3

# Optional arguments with defaults
AUDIO_FILE="$DEFAULT_AUDIO"
RESOLUTION="$DEFAULT_RESOLUTION"
BG_COLOR="$DEFAULT_BG_COLOR"
COUNTER_FONTSIZE="$DEFAULT_COUNTER_FONTSIZE"
TEXT_FONTSIZE="$DEFAULT_TEXT_FONTSIZE"
COUNTER_COLOR="$DEFAULT_COUNTER_COLOR"
TEXT_COLOR="$DEFAULT_TEXT_COLOR"
USE_AUDIO=true

# Parse remaining arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --resolution)
            RESOLUTION="$2"
            shift 2
            ;;
        --bg-color)
            BG_COLOR="$2"
            shift 2
            ;;
        --counter-size)
            COUNTER_FONTSIZE="$2"
            shift 2
            ;;
        --text-size)
            TEXT_FONTSIZE="$2"
            shift 2
            ;;
        --counter-color)
            COUNTER_COLOR="$2"
            shift 2
            ;;
        --text-color)
            TEXT_COLOR="$2"
            shift 2
            ;;
        --no-audio)
            USE_AUDIO=false
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
            # If it's not a flag and we haven't set audio file from default, assume it's audio file
            if [ "$AUDIO_FILE" == "$DEFAULT_AUDIO" ] && [ -f "$1" ]; then
                AUDIO_FILE="$1"
            else
                echo "Error: Unknown argument '$1'"
                show_usage
                exit 1
            fi
            shift
            ;;
    esac
done

# Validation
echo "Validating inputs..."

# Check if ffmpeg is installed
if ! command_exists ffmpeg; then
    echo "Error: ffmpeg is not installed or not in PATH"
    exit 1
fi

# Validate length
if ! validate_positive_int "$LENGTH"; then
    exit 1
fi

# Validate resolution
if ! validate_resolution "$RESOLUTION"; then
    exit 1
fi

# Validate font sizes
if ! validate_positive_int "$COUNTER_FONTSIZE"; then
    echo "Error: Counter font size must be a positive integer"
    exit 1
fi

if ! validate_positive_int "$TEXT_FONTSIZE"; then
    echo "Error: Text font size must be a positive integer"
    exit 1
fi

# Check if audio file exists (if using audio)
if [ "$USE_AUDIO" == true ]; then
    if ! check_file_exists "$AUDIO_FILE"; then
        echo "Hint: Use --no-audio flag to generate video without audio"
        exit 1
    fi
fi

# Check if output directory exists and is writable
OUTPUT_DIR=$(dirname "$OUTPUT_FILE")
if [ ! -d "$OUTPUT_DIR" ]; then
    echo "Error: Output directory '$OUTPUT_DIR' does not exist"
    exit 1
fi

if [ ! -w "$OUTPUT_DIR" ]; then
    echo "Error: Output directory '$OUTPUT_DIR' is not writable"
    exit 1
fi

# Escape text for ffmpeg filter
ESCAPED_TEXT=$(echo "$TEXT" | sed "s/'/'\\\\''/g" | sed 's/:/\\:/g')

echo "Generating video with the following settings:"
echo "  Text: $TEXT"
echo "  Duration: ${LENGTH}s"
echo "  Resolution: $RESOLUTION"
echo "  Background: $BG_COLOR"
echo "  Audio: $([ "$USE_AUDIO" == true ] && echo "Yes ($AUDIO_FILE)" || echo "No")"
echo "  Output: $OUTPUT_FILE"
echo

# Build FFmpeg command
FFMPEG_CMD="ffmpeg -y -f lavfi -i color=c=$BG_COLOR:s=$RESOLUTION:d=$LENGTH"

if [ "$USE_AUDIO" == true ]; then
    FFMPEG_CMD="$FFMPEG_CMD -i \"$AUDIO_FILE\""
fi

# Build filter complex
FILTER_COMPLEX="[0:v]drawtext=text='%{eif\\:n\\:d}':fontsize=$COUNTER_FONTSIZE:fontcolor=$COUNTER_COLOR:x=(w-text_w)/2:y=(h-text_h)/2-50:enable='between(t,0,$LENGTH)',drawtext=text='$ESCAPED_TEXT':fontsize=$TEXT_FONTSIZE:fontcolor=$TEXT_COLOR:x=(w-text_w)/2:y=(h-text_h)/2+100:enable='between(t,0,$LENGTH)'[v]"

FFMPEG_CMD="$FFMPEG_CMD -filter_complex \"$FILTER_COMPLEX\" -map \"[v]\""

if [ "$USE_AUDIO" == true ]; then
    FFMPEG_CMD="$FFMPEG_CMD -map 1:a -c:a aac -shortest"
fi

FFMPEG_CMD="$FFMPEG_CMD -c:v libx264 -pix_fmt yuv420p -t $LENGTH \"$OUTPUT_FILE\""

# Execute FFmpeg command
echo "Running FFmpeg..."
echo "Command: $FFMPEG_CMD"
echo

if eval "$FFMPEG_CMD"; then
    echo
    echo "✓ Video generated successfully: $OUTPUT_FILE"
    
    # Show file info
    if command_exists ffprobe; then
        echo
        echo "Video information:"
        ffprobe -v quiet -show_entries format=duration,size -show_entries stream=width,height,codec_name -of csv=p=0 "$OUTPUT_FILE" 2>/dev/null | while IFS=',' read -r width height vcodec duration size acodec; do
            echo "  Resolution: ${width}x${height}"
            echo "  Duration: ${duration}s"
            echo "  Size: $(( size / 1024 / 1024 ))MB"
            echo "  Video codec: $vcodec"
            [ -n "$acodec" ] && echo "  Audio codec: $acodec"
        done
    fi
else
    echo "✗ Error: Failed to generate video"
    exit 1
fi