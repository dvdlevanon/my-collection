#!/bin/bash

# Check if correct number of arguments provided
if [ $# -ne 3 ]; then
    echo "Usage: $0 <length_in_seconds> \"custom_text\" output.srt"
    echo "Example: $0 30 \"Hello World\" output.srt"
    exit 1
fi

# Parse arguments
LENGTH=$1
CUSTOM_TEXT=$2
OUTPUT_FILE=$3

# Validate that length is a positive number
if ! [[ "$LENGTH" =~ ^[0-9]+$ ]] || [ "$LENGTH" -le 0 ]; then
    echo "Error: Length must be a positive integer"
    exit 1
fi

# Create/clear the output file
> "$OUTPUT_FILE"

# Calculate number of subtitle segments (every 5 seconds)
SEGMENTS=$((LENGTH / 5))
REMAINING=$((LENGTH % 5))

# Function to format time for SRT (HH:MM:SS,mmm)
format_time() {
    local total_seconds=$1
    local hours=$((total_seconds / 3600))
    local minutes=$(((total_seconds % 3600) / 60))
    local seconds=$((total_seconds % 60))
    printf "%02d:%02d:%02d,000" $hours $minutes $seconds
}

# Generate subtitle segments
counter=1
for ((i=0; i<SEGMENTS; i++)); do
    start_time=$((i * 5))
    end_time=$(((i + 1) * 5))
    
    # Write subtitle entry
    echo "$counter" >> "$OUTPUT_FILE"
    echo "$(format_time $start_time) --> $(format_time $end_time)" >> "$OUTPUT_FILE"
    echo "$CUSTOM_TEXT - $counter" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    
    counter=$((counter + 1))
done

# Handle remaining time if length is not divisible by 5
if [ $REMAINING -gt 0 ]; then
    start_time=$((SEGMENTS * 5))
    end_time=$LENGTH
    
    echo "$counter" >> "$OUTPUT_FILE"
    echo "$(format_time $start_time) --> $(format_time $end_time)" >> "$OUTPUT_FILE"
    echo "$CUSTOM_TEXT - $counter" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
fi

echo "Subtitle file '$OUTPUT_FILE' generated successfully!"
echo "Duration: ${LENGTH} seconds"
echo "Segments: $counter"
