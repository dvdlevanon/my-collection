#!/bin/bash

SCRIPT_DIR=$(dirname "$(realpath "${BASH_SOURCE[0]}")")

# Media Library Directory Structure Generator
# Creates a comprehensive media library with movies, series, and short clips

# Get current timestamp for unique directory name
TIMESTAMP=$(date +%Y-%m-%d-%H%M%S)
BASE_DIR=$1
[ -z "$BASE_DIR" ] && BASE_DIR="media-library-$TIMESTAMP"

echo "Creating media library structure in: $BASE_DIR"

# Create base directory
mkdir -p "$BASE_DIR"
cd "$BASE_DIR"

# LIB1 - Movies
echo "Creating lib1 (Movies)..."
mkdir -p lib1

# Movie directories with various genres
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

for movie in "${movies[@]}"; do
    mkdir -p "lib1/$movie"
    $SCRIPT_DIR/generate-movie.sh $movie 20 "lib1/$movie/movie.mp4"
    $SCRIPT_DIR/generate-subtitle.sh 20 $movie "lib1/$movie/movie.srt"
    # Add additional subtitle files for some movies
    if [[ $movie == *"action"* ]] || [[ $movie == *"sci_fi"* ]]; then
        $SCRIPT_DIR/generate-subtitle.sh 20 "$movie - es" "lib1/$movie/movie.es.srt"
        $SCRIPT_DIR/generate-subtitle.sh 20 "$movie - fr" "lib1/$movie/movie.fr.srt"
    fi
    # Add extras for some movies
    if [[ $movie == *"drama"* ]] || [[ $movie == *"documentary"* ]]; then
        mkdir -p "lib1/$movie/extras"
        $SCRIPT_DIR/generate-movie.sh "$movie - behind the scenes" 20 "lib1/$movie/extras/behind_the_scenes.mp4"
        $SCRIPT_DIR/generate-movie.sh "$movie - deleted scenes" 20 "lib1/$movie/extras/deleted_scenes.mp4"
    fi
done

# LIB2 - TV Series
echo "Creating lib2 (TV Series)..."
mkdir -p lib2

# Series with multiple seasons
series_data=(
    "drama_series:3:8,10,6"        # 3 seasons with 8,10,6 episodes
    "comedy_series:4:12,12,12,10"  # 4 seasons
    "action_series:2:10,12"        # 2 seasons
    "sci_fi_series:5:8,8,10,10,12" # 5 seasons
    "crime_series:3:6,8,8"         # 3 seasons
    "fantasy_series:2:10,10"       # 2 seasons
)

for series_info in "${series_data[@]}"; do
    IFS=':' read -r series_name season_count episodes <<< "$series_info"
    mkdir -p "lib2/$series_name"
    
    IFS=',' read -ra EPISODES <<< "$episodes"
    for i in $(seq 1 $season_count); do
        season_dir="lib2/$series_name/season$i"
        mkdir -p "$season_dir"
        
        episode_count=${EPISODES[$((i-1))]}
        for j in $(seq 1 $episode_count); do
            episode_num=$(printf "%02d" $j)
            touch "$season_dir/S$(printf "%02d" $i)E$episode_num.mp4"
            touch "$season_dir/S$(printf "%02d" $i)E$episode_num.srt"
        done
        
        # Add season extras for some series
        if [[ $series_name == *"drama"* ]] && [ $i -eq 1 ]; then
            mkdir -p "$season_dir/extras"
            touch "$season_dir/extras/gag_reel.mp4"
            touch "$season_dir/extras/cast_interviews.mp4"
        fi
    done
done

# LIB3 - Short Clips
echo "Creating lib3 (Short Clips)..."
mkdir -p lib3

# Categories of short clips
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

for category in "${clip_categories[@]}"; do
    mkdir -p "lib3/$category"
    
    # Generate different numbers of clips per category
    case $category in
        "music_videos")
            clips=("pop_hit_2024.mp4" "rock_anthem.mp4" "indie_track.mp4" "jazz_session.mp4" "electronic_mix.mp4")
            ;;
        "trailers")
            clips=("action_movie_trailer.mp4" "horror_teaser.mp4" "comedy_preview.mp4" "drama_trailer.mp4")
            ;;
        "comedy_skits")
            clips=("office_parody.mp4" "stand_up_clip.mp4" "sketch_01.mp4" "sketch_02.mp4" "sketch_03.mp4" "blooper_reel.mp4")
            ;;
        "tutorials")
            clips=("cooking_basics.mp4" "guitar_lesson_1.mp4" "photography_tips.mp4" "coding_intro.mp4")
            ;;
        "vlogs")
            clips=("day_in_life_01.mp4" "travel_vlog_paris.mp4" "travel_vlog_tokyo.mp4" "morning_routine.mp4")
            ;;
        "gaming_highlights")
            clips=("epic_win_compilation.mp4" "funny_moments.mp4" "speedrun_record.mp4" "pvp_highlights.mp4")
            ;;
        "nature_clips")
            clips=("sunset_timelapse.mp4" "ocean_waves.mp4" "forest_ambience.mp4" "wildlife_footage.mp4")
            ;;
        "sports_highlights")
            clips=("goal_compilation.mp4" "best_saves.mp4" "tournament_recap.mp4")
            ;;
        "news_clips")
            clips=("tech_news_summary.mp4" "weather_update.mp4" "breaking_news.mp4")
            ;;
        "animations")
            clips=("2d_short.mp4" "3d_render.mp4" "stop_motion.mp4" "motion_graphics.mp4")
            ;;
    esac
    
    for clip in "${clips[@]}"; do
        touch "lib3/$category/$clip"
    done
    
    # Add metadata files for some categories
    if [[ $category == "music_videos" ]] || [[ $category == "animations" ]]; then
        touch "lib3/$category/playlist.m3u"
        touch "lib3/$category/README.txt"
    fi
done

# Add some mixed/uncategorized clips directly in lib3
mixed_clips=("viral_clip_001.mp4" "viral_clip_002.mp4" "random_funny.mp4" "old_commercial.mp4" "meme_compilation.mp4")
for clip in "${mixed_clips[@]}"; do
    touch "lib3/$clip"
done

# Create some additional organizational files
touch "lib1/movies_catalog.txt"
touch "lib2/series_watchlist.txt" 
touch "lib3/clips_index.txt"

# Create a main README
cat > README.md << EOF
# Media Library

Generated on: $(date)

## Structure:
- **lib1/**: Movie collection with subtitles and extras
- **lib2/**: TV series organized by seasons
- **lib3/**: Short clips categorized by type

## Statistics:
- Movies: $(find lib1 -name "*.mp4" | wc -l)
- TV Episodes: $(find lib2 -name "*.mp4" | wc -l) 
- Short Clips: $(find lib3 -name "*.mp4" | wc -l)

Total video files: $(find . -name "*.mp4" | wc -l)
EOF

echo "Media library structure created successfully!"
echo "Directory: $BASE_DIR"
echo ""
echo "To view the structure, run:"
echo "tree -s $BASE_DIR"
echo ""
echo "Statistics:"
echo "- Movies: $(find lib1 -name "*.mp4" 2>/dev/null | wc -l)"
echo "- TV Episodes: $(find lib2 -name "*.mp4" 2>/dev/null | wc -l)" 
echo "- Short Clips: $(find lib3 -name "*.mp4" 2>/dev/null | wc -l)"
