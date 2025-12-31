#!/bin/bash

SCRIPT_DIR=$(dirname "$(realpath "${BASH_SOURCE[0]}")")

# Media Library Directory Structure Generator
# Creates two separate hierarchies:
# 1. Category-based (lib1, lib2, lib3) - for testing category organization
# 2. Codec-based (by_codec/) - for testing codec conversion features

# Get current timestamp for unique directory name
TIMESTAMP=$(date +%Y-%m-%d-%H%M%S)
BASE_DIR=$1
[ -z "$BASE_DIR" ] && BASE_DIR="media-library-$TIMESTAMP"

echo "Creating media library structure in: $BASE_DIR"

# Create base directory
mkdir -p "$BASE_DIR"
cd "$BASE_DIR"

# Define codec combinations for testing
# Format: "video_codec:audio_codec:extension"
declare -A CODEC_PROFILES=(
    ["h264"]="libx264:aac:mp4"
    ["h265"]="libx265:aac:mp4"
    ["vp9"]="libvpx-vp9:libopus:webm"
    ["av1"]="libsvtav1:libopus:mp4"
)

# Function to generate video with specific codec
generate_video_with_codec() {
    local title="$1"
    local duration="$2"
    local output_path="$3"
    local video_codec="$4"
    local audio_codec="$5"
    
    if [ -f "$SCRIPT_DIR/generate-video.sh" ]; then
        $SCRIPT_DIR/generate-video.sh "$title" "$duration" "$output_path" \
            --video-codec "$video_codec" --audio-codec "$audio_codec"
    else
        echo "Warning: generate-video.sh not found, creating placeholder for $output_path"
        touch "$output_path"
    fi
}

# Content list for reuse
movies=(
    "action_movie_1"
    "action_movie_2" 
    "comedy_movie_1"
    "comedy_movie_2"
    "drama_movie_1"
    "horror_movie_1"
    "sci_fi_movie_1"
    "romance_movie_1"
    "thriller_movie_1"
    "documentary_1"
)

series_data=(
    "drama_series:2:3,3,3"
    "comedy_series:2:10,10"
    "action_series:1:1"
)

clip_categories=(
    "music_videos"
    "trailers" 
    "comedy_skits"
    "tutorials"
    "vlogs"
    "gaming_highlights"
    "nature_clips"
    "sports_highlights"
    "news_clips"
    "animations"
)

################################################################################
# HIERARCHY 1: CATEGORY-BASED (lib1, lib2, lib3)
################################################################################
echo ""
echo "=========================================="
echo "Creating CATEGORY-BASED hierarchy..."
echo "=========================================="

# LIB1 - Movies (standard format: H.264 + AAC)
echo "Creating lib1 (Movies)..."
mkdir -p lib1

IFS=':' read -r vcodec acodec ext <<< "${CODEC_PROFILES[h264]}"

for movie in "${movies[@]}"; do
    mkdir -p "lib1/$movie"
    echo "  Generating $movie..."
    generate_video_with_codec "$movie" 20 "lib1/$movie/movie.${ext}" "$vcodec" "$acodec"
    
    if [ -f "$SCRIPT_DIR/generate-subtitle.sh" ]; then
        $SCRIPT_DIR/generate-subtitle.sh 20 "$movie" "lib1/$movie/movie.srt"
    else
        touch "lib1/$movie/movie.srt"
    fi
    
    # Add extras for some movies
    if [[ $movie == *"drama"* ]] || [[ $movie == *"documentary"* ]]; then
        mkdir -p "lib1/$movie/extras"
        generate_video_with_codec "$movie - behind the scenes" 15 \
            "lib1/$movie/extras/behind_the_scenes.${ext}" "$vcodec" "$acodec"
        touch "lib1/$movie/extras/behind_the_scenes.srt"
    fi
done

# LIB2 - TV Series
echo "Creating lib2 (TV Series)..."
mkdir -p lib2

for series_info in "${series_data[@]}"; do
    IFS=':' read -r series_name season_count episodes <<< "$series_info"
    mkdir -p "lib2/$series_name"
    
    IFS=',' read -ra EPISODES <<< "$episodes"
    for i in $(seq 1 $season_count); do
        season_dir="lib2/$series_name/season$i"
        mkdir -p "$season_dir"
        
        echo "  Generating $series_name season $i..."
        
        episode_count=${EPISODES[$((i-1))]}
        for j in $(seq 1 $episode_count); do
            episode_num=$(printf "%02d" $j)
            season_num=$(printf "%02d" $i)
            
            # Generate only first 2 episodes per season
            if [ $j -le 2 ]; then
                generate_video_with_codec "$series_name S${season_num}E${episode_num}" 15 \
                    "$season_dir/S${season_num}E${episode_num}.${ext}" "$vcodec" "$acodec"
            else
                touch "$season_dir/S${season_num}E${episode_num}.${ext}"
            fi
            
            touch "$season_dir/S${season_num}E${episode_num}.srt"
        done
    done
done

# LIB3 - Short Clips
echo "Creating lib3 (Short Clips)..."
mkdir -p lib3

for category in "${clip_categories[@]}"; do
    mkdir -p "lib3/$category"
    
    echo "  Generating clips for $category..."
    
    case $category in
        "music_videos")
            clips=("pop_hit_2024" "rock_anthem" "indie_track" "jazz_session")
            ;;
        "trailers")
            clips=("action_movie_trailer" "horror_teaser" "comedy_preview")
            ;;
        "comedy_skits")
            clips=("office_parody" "stand_up_clip" "sketch_01")
            ;;
        "tutorials")
            clips=("cooking_basics" "guitar_lesson_1" "photography_tips")
            ;;
        *)
            clips=("clip_01" "clip_02" "clip_03")
            ;;
    esac
    
    for clip in "${clips[@]}"; do
        # Generate only first clip per category
        if [ "$clip" == "${clips[0]}" ]; then
            generate_video_with_codec "$clip" 10 \
                "lib3/$category/${clip}.${ext}" "$vcodec" "$acodec"
        else
            touch "lib3/$category/${clip}.${ext}"
        fi
    done
done

################################################################################
# HIERARCHY 2: CODEC-BASED (by_codec/)
################################################################################
echo ""
echo "=========================================="
echo "Creating CODEC-BASED hierarchy..."
echo "=========================================="

# Create codec-organized structure
for codec_name in "${!CODEC_PROFILES[@]}"; do
    IFS=':' read -r vcodec acodec ext <<< "${CODEC_PROFILES[$codec_name]}"
    
    echo "Creating codec hierarchy for: $codec_name"
    
    # Movies organized by codec
    echo "  Creating by_codec/$codec_name/movies..."
    for movie in "${movies[@]}"; do
        movie_dir="by_codec/$codec_name/movies/$movie"
        mkdir -p "$movie_dir"
        
        # Generate only first 3 movies per codec to save time
        if [[ "$movie" == "action_movie_1" ]] || [[ "$movie" == "comedy_movie_1" ]] || [[ "$movie" == "drama_movie_1" ]]; then
            echo "    Generating $movie in $codec_name..."
            generate_video_with_codec "$movie" 20 "$movie_dir/movie.${ext}" "$vcodec" "$acodec"
        else
            touch "$movie_dir/movie.${ext}"
        fi
        
        touch "$movie_dir/movie.srt"
    done
    
    # TV Series organized by codec
    echo "  Creating by_codec/$codec_name/series..."
    for series_info in "${series_data[@]}"; do
        IFS=':' read -r series_name season_count episodes <<< "$series_info"
        
        IFS=',' read -ra EPISODES <<< "$episodes"
        for i in $(seq 1 $season_count); do
            season_dir="by_codec/$codec_name/series/$series_name/season$i"
            mkdir -p "$season_dir"
            
            episode_count=${EPISODES[$((i-1))]}
            for j in $(seq 1 $episode_count); do
                episode_num=$(printf "%02d" $j)
                season_num=$(printf "%02d" $i)
                
                # Generate only first episode of first season for first series
                if [[ "$series_name" == "drama_series" ]] && [ $i -eq 1 ] && [ $j -eq 1 ]; then
                    echo "    Generating $series_name S${season_num}E${episode_num} in $codec_name..."
                    generate_video_with_codec "$series_name S${season_num}E${episode_num}" 15 \
                        "$season_dir/S${season_num}E${episode_num}.${ext}" "$vcodec" "$acodec"
                else
                    touch "$season_dir/S${season_num}E${episode_num}.${ext}"
                fi
                
                touch "$season_dir/S${season_num}E${episode_num}.srt"
            done
        done
    done
    
    # Short clips organized by codec
    echo "  Creating by_codec/$codec_name/clips..."
    for category in "${clip_categories[@]}"; do
        clip_dir="by_codec/$codec_name/clips/$category"
        mkdir -p "$clip_dir"
        
        # Generate only first category
        if [[ "$category" == "music_videos" ]]; then
            echo "    Generating clips for $category in $codec_name..."
            generate_video_with_codec "sample_clip" 10 "$clip_dir/clip_01.${ext}" "$vcodec" "$acodec"
            touch "$clip_dir/clip_02.${ext}"
            touch "$clip_dir/clip_03.${ext}"
        else
            touch "$clip_dir/clip_01.${ext}"
            touch "$clip_dir/clip_02.${ext}"
        fi
    done
done

################################################################################
# DOCUMENTATION
################################################################################

# Create codec reference file
cat > CODEC_REFERENCE.md << EOF
# Codec Reference

This media library contains two separate hierarchies for testing different features:

## 1. CATEGORY-BASED Hierarchy (lib1, lib2, lib3)

For testing category organization and media library scanning:

- **lib1/**: Movies organized by title (all in H.264+AAC)
- **lib2/**: TV series organized by show and season (all in H.264+AAC)
- **lib3/**: Short clips organized by category (all in H.264+AAC)

## 2. CODEC-BASED Hierarchy (by_codec/)

For testing codec detection and conversion features:

Structure: \`by_codec/{video_codec}/{content_type}/{title}/\`

### Available Codecs:

#### H.264 (by_codec/h264/)
- Video: libx264 (H.264/AVC)
- Audio: AAC
- Container: MP4
- Use case: Universal compatibility, baseline for conversions

#### H.265 (by_codec/h265/)
- Video: libx265 (H.265/HEVC)
- Audio: AAC
- Container: MP4
- Use case: Better compression, test transcoding performance

#### VP9 (by_codec/vp9/)
- Video: VP9
- Audio: Opus
- Container: WebM
- Use case: Web streaming, open source codec testing

#### AV1 (by_codec/av1/)
- Video: SVT-AV1 (libsvtav1)
- Audio: Opus
- Container: MP4
- Use case: Next-gen codec, much faster than libaom-av1

### Codec Hierarchy Contents:

Each codec directory contains:
- \`movies/\`: Same movie titles in different codecs
- \`series/\`: Same series/episodes in different codecs
- \`clips/\`: Same clip categories in different codecs

## Testing Scenarios:

### Category Features:
1. Library scanning and organization
2. Metadata extraction
3. Playlist generation by category
4. Search and filtering by content type

### Codec Features:
1. Codec detection and identification
2. Transcoding H.264 → H.265
3. Transcoding H.264 → VP9 (for web)
4. Transcoding to AV1 (next-gen)
5. Batch conversion by codec
6. Quality/size comparison across codecs
7. Performance benchmarking

EOF

# Create main README
cat > README.md << EOF
# Media Library Test Collection

Generated on: $(date)

## Two-Hierarchy Structure

### 1. Category-Based Hierarchy
- **lib1/**: $(find lib1 -type f -name "*.mp4" -o -name "*.webm" 2>/dev/null | wc -l) movies
- **lib2/**: $(find lib2 -type f -name "*.mp4" -o -name "*.webm" 2>/dev/null | wc -l) TV episodes
- **lib3/**: $(find lib3 -type f -name "*.mp4" -o -name "*.webm" 2>/dev/null | wc -l) short clips

Purpose: Test category organization, library scanning, and metadata management

### 2. Codec-Based Hierarchy
- **by_codec/h264/**: Movies, series, clips in H.264+AAC
- **by_codec/h265/**: Movies, series, clips in H.265+AAC
- **by_codec/vp9/**: Movies, series, clips in VP9+Opus
- **by_codec/av1/**: Movies, series, clips in AV1+Opus

Purpose: Test codec detection, conversion, and transcoding features

## Quick Stats

Total video files: $(find . -type f \( -name "*.mp4" -o -name "*.webm" \) 2>/dev/null | wc -l)

### By Hierarchy:
- Category-based: $(find lib1 lib2 lib3 -type f \( -name "*.mp4" -o -name "*.webm" \) 2>/dev/null | wc -l)
- Codec-based: $(find by_codec -type f \( -name "*.mp4" -o -name "*.webm" \) 2>/dev/null | wc -l)

### By Codec (in codec hierarchy):
- H.264: $(find by_codec/h264 -type f \( -name "*.mp4" -o -name "*.webm" \) 2>/dev/null | wc -l)
- H.265: $(find by_codec/h265 -type f \( -name "*.mp4" -o -name "*.webm" \) 2>/dev/null | wc -l)
- VP9: $(find by_codec/vp9 -type f \( -name "*.mp4" -o -name "*.webm" \) 2>/dev/null | wc -l)
- AV1: $(find by_codec/av1 -type f \( -name "*.mp4" -o -name "*.webm" \) 2>/dev/null | wc -l)

See CODEC_REFERENCE.md for detailed information and testing scenarios.
EOF

echo ""
echo "=========================================="
echo "Media library structure created successfully!"
echo "=========================================="
echo "Directory: $BASE_DIR"
echo ""
echo "HIERARCHY 1 - Category-based (lib1, lib2, lib3):"
echo "  Movies: $(find lib1 -type f \( -name "*.mp4" -o -name "*.webm" \) 2>/dev/null | wc -l)"
echo "  TV Episodes: $(find lib2 -type f \( -name "*.mp4" -o -name "*.webm" \) 2>/dev/null | wc -l)"
echo "  Clips: $(find lib3 -type f \( -name "*.mp4" -o -name "*.webm" \) 2>/dev/null | wc -l)"
echo ""
echo "HIERARCHY 2 - Codec-based (by_codec/):"
echo "  H.264: $(find by_codec/h264 -type f \( -name "*.mp4" -o -name "*.webm" \) 2>/dev/null | wc -l) files"
echo "  H.265: $(find by_codec/h265 -type f \( -name "*.mp4" -o -name "*.webm" \) 2>/dev/null | wc -l) files"
echo "  VP9: $(find by_codec/vp9 -type f \( -name "*.mp4" -o -name "*.webm" \) 2>/dev/null | wc -l) files"
echo "  AV1: $(find by_codec/av1 -type f \( -name "*.mp4" -o -name "*.webm" \) 2>/dev/null | wc -l) files"
echo ""
echo "To view the structure:"
echo "  tree -L 3 $BASE_DIR"
