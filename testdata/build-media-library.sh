#!/bin/bash

SCRIPT_DIR=$(dirname "$(realpath "${BASH_SOURCE[0]}")")
BASE_DIR=$1

echo "Creating media library structure in: $BASE_DIR"

# Create base directory
mkdir -p "$BASE_DIR"
cd "$BASE_DIR"

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

# Movies
vcodec=libsvtav1
acodec=libopus
ext=mp4
echo "Creating movies (Movies)..."
mkdir -p movies
for movie in "${movies[@]}"; do
    mkdir -p "movies/$movie"
    echo "  Generating $movie..."
    generate_video_with_codec "$movie" 20 "movies/$movie/movie.${ext}" "$vcodec" "$acodec"
    $SCRIPT_DIR/generate-subtitle.sh 20 "$movie" "movies/$movie/movie.srt"
    
    # Add extras for some movies
    if [[ $movie == *"drama"* ]] || [[ $movie == *"documentary"* ]]; then
        mkdir -p "movies/$movie/extras"
        generate_video_with_codec "$movie - behind the scenes" 15 \
            "movies/$movie/extras/behind_the_scenes.${ext}" "$vcodec" "$acodec"
    fi
done

# TV Series
vcodec=libx264
acodec=libmp3lame
ext=avi
echo "Creating TV Series..."
mkdir -p tv-series
for series_info in "${series_data[@]}"; do
    IFS=':' read -r series_name season_count episodes <<< "$series_info"
    mkdir -p "tv-series/$series_name"
    
    IFS=',' read -ra EPISODES <<< "$episodes"
    for i in $(seq 1 $season_count); do
        season_dir="tv-series/$series_name/season$i"
        mkdir -p "$season_dir"
        
        echo "  Generating $series_name season $i..."
        
        episode_count=${EPISODES[$((i-1))]}
        for j in $(seq 1 $episode_count); do
            episode_num=$(printf "%02d" $j)
            season_num=$(printf "%02d" $i)
            
            # Generate only first 2 episodes per season
            generate_video_with_codec "$series_name S${season_num}E${episode_num}" 5 \
                "$season_dir/S${season_num}E${episode_num}.${ext}" "$vcodec" "$acodec"
            
            $SCRIPT_DIR/generate-subtitle.sh 5 "$series_name S${season_num}E${episode_num}" "$season_dir/S${season_num}E${episode_num}.srt"
        done
    done
done

# Short Clips
vcodec=libx265
acodec=aac
ext=mkv
echo "Creating Short Clips..."
mkdir -p clips
for category in "${clip_categories[@]}"; do
    mkdir -p "clips/$category"
    
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
        generate_video_with_codec "$clip" 10 \
            "clips/$category/${clip}.${ext}" "$vcodec" "$acodec"
    done
done
