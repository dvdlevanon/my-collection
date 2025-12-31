# Quick Start Guide

## TL;DR

```bash
# Generate complete test library
./build-test-library.sh

# Result: test-library/ directory with all edge cases
```

## What You Get

```
test-library/
├── by_issue/              # Edge cases organized by problem
│   ├── codec_variations/       # 11 codec combinations
│   ├── audio_tracks/           # 6 multi-audio scenarios
│   ├── subtitle_variations/    # 7 subtitle edge cases
│   ├── resolution_variations/  # 16 resolution tests
│   ├── broken_files/           # 10 corrupted file scenarios
│   └── metadata_issues/        # 10 encoding problems
│
├── lib1/                  # Movies (category hierarchy)
├── lib2/                  # TV Series
├── lib3/                  # Short clips
└── by_codec/              # Same content, different codecs
    ├── h264/
    ├── h265/
    ├── vp9/
    └── av1/
```

## Common Commands

### Generate Everything
```bash
./build-test-library.sh
```

### Only Edge Cases (Fast)
```bash
./build-test-library.sh --issue-only
```

### Only Categories (For Filesystem Testing)
```bash
./build-test-library.sh --category-only
```

### Clean and Rebuild
```bash
./build-test-library.sh --clean
```

### Custom Location
```bash
./build-test-library.sh /path/to/output
```

## Individual Test Generators

### Codec Variations
```bash
./generate-codec-variations.sh ./output/codecs
```
Creates: H.264, H.265, VP9, AV1, DivX, XviD, etc.

### Multi-Audio Tracks
```bash
./generate-audio-tracks.sh ./output/audio
```
Creates: Dual language, 5.1 surround, 7.1, mixed codecs

### Subtitle Problems
```bash
./generate-subtitle-variations.sh ./output/subs
```
Creates: External, embedded, multi-language, Hebrew encoding

### Resolution Edge Cases
```bash
./generate-resolution-variations.sh ./output/resolutions
```
Creates: 480p, 720p, 1080p, 4K, portrait, ultra-wide

### Broken Files
```bash
./generate-broken-files.sh ./output/broken
```
Creates: Truncated, zero-byte, corrupted, wrong extensions

### Metadata Issues
```bash
./generate-metadata-issues.sh ./output/metadata
```
Creates: VFR, 23.976fps, interlaced, wrong color space

## Testing Your Media Library

### 1. Generate Test Library
```bash
./build-test-library.sh
```

### 2. Scan With Your Tool
```bash
my-collection scan test-library/by_issue
my-collection scan test-library/lib1
```

### 3. Check Specific Edge Cases
```bash
# Test broken file handling
my-collection scan test-library/by_issue/broken_files

# Test multi-audio selection
my-collection scan test-library/by_issue/audio_tracks

# Test codec conversion
my-collection convert test-library/by_codec/h264 --target h265
```

## Expected Results

After generation:
- **60+ test cases** covering major edge cases
- **100+ video files** (most are 10-20 seconds)
- **~500MB-2GB** total size (depends on placeholders)
- **10-30 minutes** generation time

## Files You Need

### Required
- `build-test-library.sh` - Master script
- All `generate-*.sh` scripts
- `mux-container.sh`, `create-media-library.sh`
- Your `audio.mp3` file (for audio source)

### Optional
- `generate-subtitle.sh` (can use existing if you have it)

## Troubleshooting

### "ffmpeg: command not found"
```bash
# Install ffmpeg first
brew install ffmpeg              # macOS
sudo apt install ffmpeg          # Ubuntu/Debian
```

### "AV1 encoding too slow"
Edit `generate-video.sh` line with AV1 codec:
```bash
# Change from libsvtav1 to libaom-av1
# Or skip AV1 tests entirely
```

### "Out of disk space"
Generate selectively:
```bash
./build-test-library.sh --issue-only  # Skip category hierarchy
# Or
./generate-codec-variations.sh ./output  # Just one generator
```

### "Permission denied"
```bash
chmod +x *.sh  # Make scripts executable
```

## Next Steps

1. ✅ Generate test library
2. ✅ Run your media scanner on it
3. ✅ Fix bugs found
4. ✅ Add new edge cases you discover
5. ✅ Regenerate and retest

## Need Help?

See `OVERVIEW.md` for:
- Complete architecture documentation
- Individual script usage
- Customization guide
- All test cases explained

## Quick Reference

| Task | Command |
|------|---------|
| Generate everything | `./build-test-library.sh` |
| Only edge cases | `./build-test-library.sh --issue-only` |
| Clean rebuild | `./build-test-library.sh --clean` |
| View structure | `tree -L 2 test-library` |
| Check size | `du -sh test-library` |
| List all videos | `find test-library -name '*.mp4'` |
| Remove test library | `rm -rf test-library` |

---

**Pro Tip:** Start with `--issue-only` for faster iteration during development. Add category hierarchy later for full filesystem testing.
