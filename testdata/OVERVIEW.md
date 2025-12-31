# Test Library Generator for my-collection

Complete test suite generator for comprehensive media library edge case testing.

## Architecture

The test suite follows a layered architecture:

```
Layer 1: Building Blocks (Primitives)
    ├── generate-video.sh        - Generate videos with specific codecs
    ├── generate-audio.sh        - Generate audio with various formats
    └── generate-subtitle.sh     - Generate SRT subtitle files

Layer 2: Container Muxer
    └── mux-container.sh         - Combine streams into containers (MKV, MP4, AVI, etc.)

Layer 3: Issue Generators (Domain-specific)
    ├── generate-codec-variations.sh       - Different codec combinations
    ├── generate-audio-tracks.sh           - Multi-audio scenarios
    ├── generate-subtitle-variations.sh    - Subtitle edge cases
    ├── generate-resolution-variations.sh  - Resolution/aspect ratios
    ├── generate-broken-files.sh           - Corrupted files
    └── generate-metadata-issues.sh        - Frame rates, color spaces

Layer 4: Hierarchy Builders
    ├── build-by-issue-hierarchy.sh   - Organizes by problem type
    └── create-media-library.sh       - Organizes by category

Layer 5: Master Orchestrator
    └── build-test-library.sh         - Runs everything
```

## Quick Start

### Generate Complete Test Library

```bash
./build-test-library.sh
```

This creates:
```
test-library/
├── by_issue/           # Edge case tests organized by problem
│   ├── codec_variations/
│   ├── audio_tracks/
│   ├── subtitle_variations/
│   ├── resolution_variations/
│   ├── broken_files/
│   └── metadata_issues/
├── lib1/               # Movies (category-based)
├── lib2/               # TV Series
├── lib3/               # Short clips
└── by_codec/           # Same content, different codecs
```

### Generate Only Edge Cases

```bash
./build-test-library.sh --issue-only
```

### Generate Only Category Hierarchy

```bash
./build-test-library.sh --category-only
```

### Clean and Rebuild

```bash
./build-test-library.sh --clean
```

## Individual Script Usage

### Building Blocks

#### Generate Video
```bash
# Basic usage
./generate-video.sh "My Video" 20 output.mp4

# With specific codecs
./generate-video.sh "Test" 30 output.mp4 \
    --video-codec libx265 \
    --audio-codec ac3 \
    --resolution 1920x1080

# Silent video
./generate-video.sh "Silent" 15 output.mp4 --no-audio

# With custom audio
./generate-video.sh "Custom" 20 output.mp4 my-audio.mp3
```

#### Generate Audio
```bash
# From source (transcode)
./generate-audio.sh 20 output.ac3 --codec ac3 --channels 6

# Generate synthetic
./generate-audio.sh 30 beep.mp3 --generate --tone 1000

# Different formats
./generate-audio.sh 15 audio.dts --codec dts --sample-rate 48000
```

#### Generate Subtitles
```bash
./generate-subtitle.sh 30 "My subtitle text" output.srt
```

### Container Muxer

```bash
# Simple mux
./mux-container.sh --video video.mp4 --audio audio.ac3 --output movie.mkv

# Multiple audio tracks
./mux-container.sh --video video.mp4 \
    --audio eng.aac --audio-lang eng --audio-title "English" \
    --audio heb.ac3 --audio-lang heb --audio-title "Hebrew" \
    --output movie.mkv

# With subtitles
./mux-container.sh --video video.mp4 --audio audio.aac \
    --subtitle eng.srt --sub-lang eng \
    --subtitle heb.srt --sub-lang heb --sub-forced \
    --output movie.mkv

# Fast copy (no re-encoding)
./mux-container.sh --video video.mp4 --audio audio.aac \
    --subtitle subs.srt --output movie.mkv --copy-codecs
```

### Issue Generators

```bash
# Run individual generators
./generate-codec-variations.sh ./output/codecs
./generate-audio-tracks.sh ./output/audio
./generate-subtitle-variations.sh ./output/subs
./generate-resolution-variations.sh ./output/resolutions
./generate-broken-files.sh ./output/broken
./generate-metadata-issues.sh ./output/metadata
```

### Hierarchy Builders

```bash
# Build by-issue hierarchy
./build-by-issue-hierarchy.sh ./test-library

# Build category hierarchy (already integrated in create-media-library.sh)
./create-media-library.sh ./test-library
```

## Test Cases Covered

### Codec Variations
- Old codecs (DivX, XviD, MPEG2)
- Modern codecs (H.264, H.265, VP9, AV1)
- Various containers (MKV, AVI, MP4, MOV, WebM)
- Audio codecs (AAC, MP3, AC3, DTS, Opus)

### Audio Tracks
- Multiple languages (English, Hebrew)
- Channel configurations (mono, stereo, 5.1, 7.1)
- Mixed codecs in same file
- No audio tracks

### Subtitles
- External SRT files
- Embedded subtitles
- Multiple languages
- Forced subtitles
- ASS/SSA styled subtitles
- Hebrew encoding issues

### Resolutions
- Standard (480p, 720p, 1080p, 4K)
- Non-standard (DVD PAL/NTSC, VCD)
- Portrait orientation (phone videos)
- Ultra-wide aspect ratios
- Wrong aspect ratio metadata

### Broken Files
- Zero-byte files
- Truncated videos
- Wrong extensions
- Corrupted headers
- Very short clips (<1s)
- Symlinks to non-existent files
- Special characters in names
- Hidden system files

### Metadata Issues
- Variable frame rate (VFR)
- Non-standard frame rates (23.976, 29.97, 50, 60)
- Interlaced vs progressive
- Different color spaces (BT.709, BT.601)
- Different pixel formats (yuv420p, yuv422p)
- Missing creation dates
- HDR metadata

## Directory Structure

```
.
├── OVERVIEW.md                           # This file
├── build-test-library.sh                 # Master orchestrator
│
├── Building Blocks/
│   ├── generate-video.sh
│   ├── generate-audio.sh
│   └── generate-subtitle.sh
│
├── Container Muxer/
│   └── mux-container.sh
│
├── Issue Generators/
│   ├── generate-codec-variations.sh
│   ├── generate-audio-tracks.sh
│   ├── generate-subtitle-variations.sh
│   ├── generate-resolution-variations.sh
│   ├── generate-broken-files.sh
│   └── generate-metadata-issues.sh
│
└── Hierarchy Builders/
    ├── build-by-issue-hierarchy.sh
    └── create-media-library.sh
```

## Performance Notes

- Videos are optimized for speed (ultrafast presets)
- SVT-AV1 used instead of libaom-av1 (5-10x faster)
- Most videos are 10-20 seconds for quick testing
- Placeholders used for repetitive content
- Full generation takes 10-30 minutes depending on system

## Requirements

### Essential
- ffmpeg (with libx264, libx265, libvpx-vp9)
- bash 4.0+

### Recommended
- ffprobe (for metadata display)
- libsvtav1 (for faster AV1 encoding)
- sox (for advanced audio generation)

### Optional
- tree (for viewing directory structure)
- mediainfo (for detailed file inspection)

## Customization

### Add New Test Case

1. Create issue generator script:
```bash
./generate-my-issue.sh <output_dir>
```

2. Add to `build-by-issue-hierarchy.sh`:
```bash
generators+=(
    "generate-my-issue.sh:my_issue:My Issue Description"
)
```

3. Rebuild:
```bash
./build-test-library.sh --clean
```

### Modify Existing Test Case

Edit the appropriate generator script and rebuild:
```bash
./build-test-library.sh --issue-only --clean
```

## Testing Workflow

### 1. Generate Test Library
```bash
./build-test-library.sh
```

### 2. Run my-collection Against It
```bash
my-collection scan test-library/by_issue/codec_variations
my-collection scan test-library/lib1
```

### 3. Verify Edge Case Handling
Check logs for:
- Proper codec detection
- Error recovery on broken files
- Correct metadata extraction
- Audio track selection
- Subtitle loading

### 4. Test Specific Scenarios
```bash
# Test codec conversion
my-collection convert test-library/by_codec/h264 --target-codec h265

# Test multi-audio handling
my-collection scan test-library/by_issue/audio_tracks

# Test subtitle handling
my-collection scan test-library/by_issue/subtitle_variations
```

## Troubleshooting

### ffmpeg not found
```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt install ffmpeg

# Arch
sudo pacman -S ffmpeg
```

### SVT-AV1 not available
Edit `generate-video.sh` and change:
```bash
["av1"]="libsvtav1:libopus:mp4"
```
to:
```bash
["av1"]="libaom-av1:libopus:mp4"
```

Note: libaom-av1 is much slower.

### Generation too slow
Use `--issue-only` or run individual generators:
```bash
./generate-codec-variations.sh ./test-output
```

### Out of disk space
Videos are small but numerous. To reduce:
1. Edit generator scripts to create fewer test cases
2. Or generate only what you need:
```bash
./generate-codec-variations.sh ./output
```

## Contributing

To add new test scenarios:

1. Identify edge case from real 20TB library
2. Create generator script following existing patterns
3. Add to `build-by-issue-hierarchy.sh`
4. Update this OVERVIEW.md
5. Test generation and cleanup

## License

Scripts are provided as-is for testing my-collection.
Generated media files are synthetic and contain no copyrighted content.
