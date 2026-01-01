#!/bin/bash

# generate-audio.sh
# Generates audio files with various codecs, channels, and sample rates
# Can either generate synthetic audio or transcode from source

SCRIPT_DIR=$(dirname "$(realpath "${BASH_SOURCE[0]}")")
DEFAULT_SOURCE_AUDIO="$SCRIPT_DIR/audio.mp3"

show_usage() {
    cat << EOF
Usage: $0 <duration> <output_file> [options]

Arguments:
    duration        Duration in seconds
    output_file     Output audio file path

Options:
    --source FILE   Source audio file to transcode (default: $DEFAULT_SOURCE_AUDIO)
    --codec CODEC   Audio codec (default: auto-detect from extension)
                    Options: aac, libmp3lame, ac3, libopus, flac, pcm_s16le
    --channels N    Number of channels: 1 (mono), 2 (stereo), 6 (5.1), 8 (7.1)
    --sample-rate N Sample rate in Hz (default: 48000)
                    Common: 44100, 48000, 96000
    --bitrate N     Audio bitrate (e.g., 128k, 320k)
    --generate      Generate synthetic audio instead of transcoding
    --tone FREQ     Tone frequency for synthetic audio (default: 440)
    --help, -h      Show this help

Examples:
    # Transcode to AC3 5.1
    $0 20 output.ac3 --channels 6

    # Generate synthetic stereo MP3
    $0 30 output.mp3 --generate --channels 2

    # Generate mono tone
    $0 10 beep.wav --generate --channels 1 --tone 1000
EOF
}

# Defaults
DURATION=""
OUTPUT_FILE=""
SOURCE_AUDIO="$DEFAULT_SOURCE_AUDIO"
CODEC=""
CHANNELS=""
SAMPLE_RATE=""
BITRATE=""
GENERATE=false
TONE_FREQ=440

# Parse arguments
if [ $# -lt 2 ]; then
    show_usage
    exit 1
fi

DURATION="$1"
OUTPUT_FILE="$2"
shift 2

while [[ $# -gt 0 ]]; do
    case $1 in
        --source)
            SOURCE_AUDIO="$2"
            shift 2
            ;;
        --codec)
            CODEC="$2"
            shift 2
            ;;
        --channels)
            CHANNELS="$2"
            shift 2
            ;;
        --sample-rate)
            SAMPLE_RATE="$2"
            shift 2
            ;;
        --bitrate)
            BITRATE="$2"
            shift 2
            ;;
        --generate)
            GENERATE=true
            shift
            ;;
        --tone)
            TONE_FREQ="$2"
            shift 2
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

# Validate duration
if ! [[ "$DURATION" =~ ^[0-9]+$ ]] || [ "$DURATION" -le 0 ]; then
    echo "Error: Duration must be a positive integer"
    exit 1
fi

# Auto-detect codec from extension if not specified
if [ -z "$CODEC" ]; then
    EXT="${OUTPUT_FILE##*.}"
    case "${EXT,,}" in
        mp3) CODEC="libmp3lame" ;;
        aac|m4a) CODEC="aac" ;;
        ac3) CODEC="ac3" ;;
        opus) CODEC="libopus" ;;
        flac) CODEC="flac" ;;
        wav) CODEC="pcm_s16le" ;;
        ogg) CODEC="libvorbis" ;;
        *) CODEC="aac" ;;  # Default
    esac
fi

# Build ffmpeg command
if [ "$GENERATE" = true ]; then
    # Generate synthetic audio
    echo "Generating synthetic audio..."
    
    INPUT_OPTS="-f lavfi"
    
    if [ -n "$CHANNELS" ] && [ "$CHANNELS" -gt 1 ]; then
        # Multi-channel sine wave
        INPUT_OPTS="$INPUT_OPTS -i sine=frequency=$TONE_FREQ:duration=$DURATION"
    else
        # Simple mono sine
        INPUT_OPTS="$INPUT_OPTS -i sine=frequency=$TONE_FREQ:duration=$DURATION"
    fi
else
    # Transcode from source
    if [ ! -f "$SOURCE_AUDIO" ]; then
        echo "Error: Source audio file not found: $SOURCE_AUDIO"
        exit 1
    fi
    
    echo "Transcoding from: $SOURCE_AUDIO"
    INPUT_OPTS="-i \"$SOURCE_AUDIO\""
fi

# Build output options
OUTPUT_OPTS="-c:a $CODEC"

[ -n "$CHANNELS" ] && OUTPUT_OPTS="$OUTPUT_OPTS -ac $CHANNELS"
[ -n "$SAMPLE_RATE" ] && OUTPUT_OPTS="$OUTPUT_OPTS -ar $SAMPLE_RATE"
[ -n "$BITRATE" ] && OUTPUT_OPTS="$OUTPUT_OPTS -b:a $BITRATE"

OUTPUT_OPTS="$OUTPUT_OPTS -t $DURATION"

# Execute
FFMPEG_CMD="ffmpeg -y $INPUT_OPTS $OUTPUT_OPTS \"$OUTPUT_FILE\""

echo "Command: $FFMPEG_CMD"
echo

if eval "$FFMPEG_CMD" 2>&1; then
    echo
    echo "✓ Audio file generated: $OUTPUT_FILE"
    
    if command -v ffprobe >/dev/null || exit 1; then
        PROBE_INFO=$(ffprobe -v quiet -select_streams a:0 -show_entries stream=codec_name,channels,sample_rate -of csv=p=0 "$OUTPUT_FILE" 2>/dev/null)
        IFS=',' read -r codec channels srate <<< "$PROBE_INFO"
        echo "  Codec: $codec"
        echo "  Channels: $channels"
        echo "  Sample rate: ${srate}Hz"
    fi
else
    echo "✗ Error generating audio file"
    exit 1
fi
