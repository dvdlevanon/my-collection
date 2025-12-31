# Test Library Generator for my-collection

Complete test suite for media library edge case testing - generates 60+ test cases covering codecs, audio tracks, subtitles, resolutions, broken files, and metadata issues.

## 🚀 Quick Start

```bash
# Generate complete test library
./build-test-library.sh

# Result: test-library/ with all edge cases
```

See [QUICKSTART.md](QUICKSTART.md) for more.

## 📁 Script Inventory

### Master Scripts
- **`build-test-library.sh`** - Main orchestrator, runs everything
- **`build-by-issue-hierarchy.sh`** - Builds by_issue/ hierarchy
- **`create-media-library.sh`** - Builds category hierarchy (lib1/lib2/lib3)

### Building Blocks (Layer 1)
- **`generate-video.sh`** - Create videos with specific codecs/settings
- **`generate-audio.sh`** - Generate audio files with various formats
- **`generate-subtitle.sh`** - Create SRT subtitle files

### Container Muxer (Layer 2)
- **`mux-container.sh`** - Combine video/audio/subs into containers (MKV, MP4, etc.)

### Issue Generators (Layer 3)
- **`generate-codec-variations.sh`** - Different codec combinations (H.264, H.265, VP9, AV1, etc.)
- **`generate-audio-tracks.sh`** - Multi-audio scenarios (dual language, 5.1, 7.1)
- **`generate-subtitle-variations.sh`** - Subtitle edge cases (external, embedded, multi-lang)
- **`generate-resolution-variations.sh`** - Resolution/aspect ratios (480p-4K, portrait, etc.)
- **`generate-broken-files.sh`** - Corrupted/problematic files
- **`generate-metadata-issues.sh`** - Frame rates, color spaces, interlacing

## 📚 Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - Get started in 2 minutes
- **[OVERVIEW.md](OVERVIEW.md)** - Complete architecture and usage guide

## 🎯 What You Get

```
test-library/
├── by_issue/              # 60+ edge cases organized by problem type
│   ├── codec_variations/       # 11 test cases
│   ├── audio_tracks/           # 6 test cases
│   ├── subtitle_variations/    # 7 test cases
│   ├── resolution_variations/  # 16 test cases
│   ├── broken_files/           # 10 test cases
│   └── metadata_issues/        # 10 test cases
│
├── lib1/                  # Movies (standard H.264/AAC)
├── lib2/                  # TV Series
├── lib3/                  # Short clips
└── by_codec/              # Same content in H.264, H.265, VP9, AV1
```

## 🔥 Common Use Cases

### Test Codec Detection
```bash
./build-test-library.sh --issue-only
my-collection scan test-library/by_issue/codec_variations
```

### Test Multi-Audio Handling
```bash
./generate-audio-tracks.sh ./test-output
my-collection scan test-output
```

### Test Subtitle Loading
```bash
./generate-subtitle-variations.sh ./test-output
my-collection scan test-output
```

### Test Error Handling
```bash
./generate-broken-files.sh ./test-output
my-collection scan test-output  # Should handle gracefully
```

### Test Codec Conversion
```bash
./build-test-library.sh
my-collection convert test-library/by_codec/h264 --target h265
```

## ✅ Test Coverage

| Category | Test Cases | Coverage |
|----------|------------|----------|
| **Codecs** | 11 | H.264, H.265, VP9, AV1, DivX, XviD, MPEG2, various audio codecs |
| **Audio** | 6 | Multi-language, stereo, 5.1, 7.1, mixed codecs, no audio |
| **Subtitles** | 7 | External, embedded, multi-lang, forced, Hebrew encoding |
| **Resolution** | 16 | 480p-4K, DVD, portrait, ultra-wide, wrong metadata |
| **Broken Files** | 10 | Truncated, zero-byte, corrupted, wrong extensions |
| **Metadata** | 10 | VFR, 23.976fps, interlaced, color spaces, HDR |

## 🛠️ Requirements

### Essential
- `ffmpeg` with libx264, libx265, libvpx-vp9
- `bash` 4.0+

### Recommended  
- `libsvtav1` (5-10x faster AV1 encoding)
- `ffprobe` (metadata display)

### Optional
- `tree` (view directory structure)
- `mediainfo` (detailed file inspection)

## 📊 Performance

- **Generation Time:** 10-30 minutes (depends on system)
- **Total Size:** ~500MB-2GB
- **Video Count:** 100+ files
- **Video Length:** 10-20 seconds each (fast testing)

## 🎨 Architecture Diagram

```
build-test-library.sh (Master)
    │
    ├─→ build-by-issue-hierarchy.sh
    │       │
    │       ├─→ generate-codec-variations.sh
    │       ├─→ generate-audio-tracks.sh
    │       ├─→ generate-subtitle-variations.sh
    │       ├─→ generate-resolution-variations.sh
    │       ├─→ generate-broken-files.sh
    │       └─→ generate-metadata-issues.sh
    │               │
    │               └─→ Uses: generate-video.sh
    │                         generate-audio.sh
    │                         generate-subtitle.sh
    │                         mux-container.sh
    │
    └─→ create-media-library.sh
            │
            └─→ Uses: generate-video.sh
                      generate-audio.sh
                      generate-subtitle.sh
```

## 💡 Tips

### Speed Up Generation
```bash
# Only edge cases (skip category hierarchy)
./build-test-library.sh --issue-only

# Or run individual generators
./generate-codec-variations.sh ./output
```

### Reduce Disk Usage
Edit generator scripts to create fewer test cases, or generate selectively.

### Custom Test Cases
1. Create new `generate-my-test.sh` following existing patterns
2. Add to `build-by-issue-hierarchy.sh`
3. Rebuild

## 🐛 Troubleshooting

| Problem | Solution |
|---------|----------|
| "ffmpeg not found" | Install: `brew install ffmpeg` or `apt install ffmpeg` |
| AV1 too slow | Use libsvtav1 or skip AV1 tests |
| Out of disk | Use `--issue-only` or selective generation |
| Permission denied | `chmod +x *.sh` |

## 📖 Examples

### Generate Only Broken Files Tests
```bash
./generate-broken-files.sh ./my-broken-tests
```

### Create Custom Video
```bash
./generate-video.sh "My Test" 30 output.mp4 \
    --video-codec libx265 \
    --audio-codec ac3 \
    --resolution 1920x1080
```

### Mux Multi-Audio Container
```bash
./mux-container.sh --video video.mp4 \
    --audio eng.aac --audio-lang eng \
    --audio heb.ac3 --audio-lang heb \
    --subtitle eng.srt --sub-lang eng \
    --output movie.mkv
```

## 🔄 Workflow

1. **Generate:** `./build-test-library.sh`
2. **Scan:** Run your media library scanner
3. **Fix:** Address issues found
4. **Add:** Create new test cases for discovered edge cases
5. **Regenerate:** `./build-test-library.sh --clean`
6. **Repeat:** Until all edge cases handled

## 📝 Files Overview

| File | Purpose | Type |
|------|---------|------|
| build-test-library.sh | Main orchestrator | Master |
| build-by-issue-hierarchy.sh | Edge case organizer | Builder |
| create-media-library.sh | Category organizer | Builder |
| generate-video.sh | Video generator | Primitive |
| generate-audio.sh | Audio generator | Primitive |
| generate-subtitle.sh | Subtitle generator | Primitive |
| mux-container.sh | Container muxer | Utility |
| generate-codec-variations.sh | Codec tests | Issue Gen |
| generate-audio-tracks.sh | Audio tests | Issue Gen |
| generate-subtitle-variations.sh | Subtitle tests | Issue Gen |
| generate-resolution-variations.sh | Resolution tests | Issue Gen |
| generate-broken-files.sh | Broken file tests | Issue Gen |
| generate-metadata-issues.sh | Metadata tests | Issue Gen |

## 🚦 Status Indicators

When you run the scripts, you'll see:
- ✓ = Successfully created
- ✗ = Failed to create
- Warning = Non-critical issue

## 🎓 Learning Resources

- Read [OVERVIEW.md](OVERVIEW.md) for complete architecture
- Read [QUICKSTART.md](QUICKSTART.md) for immediate usage
- Check individual script `--help` for detailed options

## 🤝 Contributing

Found a new edge case in your 20TB library? 
1. Document it
2. Create generator script
3. Add to `build-by-issue-hierarchy.sh`
4. Test and commit

## 📜 License

Scripts provided as-is for testing my-collection.
Generated media is synthetic and contains no copyrighted content.

---

**Made for testing my-collection** - A media library that pre-optimizes everything for streaming.

**Get Started:** `./build-test-library.sh`
