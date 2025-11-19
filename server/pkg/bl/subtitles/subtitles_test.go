package subtitles

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookForAvailableSubtitles(t *testing.T) {
	t.Run("finds single srt file", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "subtitles-test-*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Create a test .srt file
		srtFile := filepath.Join(tempDir, "test.srt")
		err = os.WriteFile(srtFile, []byte("test content"), 0644)
		assert.NoError(t, err)

		names, err := lookForAvailableSubtitles(tempDir)

		assert.NoError(t, err)
		assert.Len(t, names, 1)
		assert.Equal(t, "test.srt", names[0].Title)
	})

	t.Run("finds multiple srt files", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "subtitles-test-*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Create multiple .srt files
		files := []string{"sub1.srt", "sub2.srt", "sub3.srt"}
		for _, f := range files {
			srtFile := filepath.Join(tempDir, f)
			err = os.WriteFile(srtFile, []byte("test content"), 0644)
			assert.NoError(t, err)
		}

		names, err := lookForAvailableSubtitles(tempDir)

		assert.NoError(t, err)
		assert.Len(t, names, 3)
		// Verify all expected files are found
		for _, expectedFile := range files {
			found := false
			for _, name := range names {
				if name.Title == expectedFile {
					found = true
					break
				}
			}
			assert.True(t, found, "Expected to find %s", expectedFile)
		}
	})

	t.Run("finds srt file in directory with spaces", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "subtitles test with spaces-*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Create a test .srt file in a directory with spaces
		srtFile := filepath.Join(tempDir, "my subtitle.srt")
		err = os.WriteFile(srtFile, []byte("test content"), 0644)
		assert.NoError(t, err)

		names, err := lookForAvailableSubtitles(tempDir)

		assert.NoError(t, err)
		assert.Len(t, names, 1)
		assert.Equal(t, "my subtitle.srt", names[0].Title)
	})

	t.Run("finds srt file when directory path has spaces", func(t *testing.T) {
		// Create a temp directory with spaces in its name
		tempBase, err := os.MkdirTemp("", "base-*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempBase)

		// Create a subdirectory with spaces
		subDir := filepath.Join(tempBase, "my video folder")
		err = os.MkdirAll(subDir, 0755)
		assert.NoError(t, err)

		// Create a test .srt file in the directory with spaces
		srtFile := filepath.Join(subDir, "subtitle.srt")
		err = os.WriteFile(srtFile, []byte("test content"), 0644)
		assert.NoError(t, err)

		names, err := lookForAvailableSubtitles(subDir)

		assert.NoError(t, err)
		assert.Len(t, names, 1)
		assert.Equal(t, "subtitle.srt", names[0].Title)
	})

	t.Run("finds multiple srt files in directory with spaces", func(t *testing.T) {
		// Create a temp directory with spaces in its name
		tempBase, err := os.MkdirTemp("", "base-*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempBase)

		// Create a subdirectory with spaces
		subDir := filepath.Join(tempBase, "my collection folder")
		err = os.MkdirAll(subDir, 0755)
		assert.NoError(t, err)

		// Create multiple .srt files
		files := []string{"english.srt", "spanish subtitle.srt", "french.srt"}
		for _, f := range files {
			srtFile := filepath.Join(subDir, f)
			err = os.WriteFile(srtFile, []byte("test content"), 0644)
			assert.NoError(t, err)
		}

		names, err := lookForAvailableSubtitles(subDir)

		assert.NoError(t, err)
		assert.Len(t, names, 3)
		// Verify all expected files are found
		for _, expectedFile := range files {
			found := false
			for _, name := range names {
				if name.Title == expectedFile {
					found = true
					break
				}
			}
			assert.True(t, found, "Expected to find %s", expectedFile)
		}
	})

	t.Run("returns empty slice when no srt files found", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "subtitles-test-*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		names, err := lookForAvailableSubtitles(tempDir)

		assert.NoError(t, err)
		assert.NotNil(t, names)
		assert.Len(t, names, 0)
	})

	t.Run("ignores non-srt files", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "subtitles-test-*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Create non-srt files
		otherFiles := []string{"video.mp4", "audio.mp3", "text.txt"}
		for _, f := range otherFiles {
			file := filepath.Join(tempDir, f)
			err = os.WriteFile(file, []byte("test content"), 0644)
			assert.NoError(t, err)
		}

		// Create one .srt file
		srtFile := filepath.Join(tempDir, "subtitle.srt")
		err = os.WriteFile(srtFile, []byte("test content"), 0644)
		assert.NoError(t, err)

		names, err := lookForAvailableSubtitles(tempDir)

		assert.NoError(t, err)
		assert.Len(t, names, 1)
		assert.Equal(t, "subtitle.srt", names[0].Title)
	})

	t.Run("returns only base names, not full paths", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "subtitles-test-*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Create a test .srt file
		srtFile := filepath.Join(tempDir, "test.srt")
		err = os.WriteFile(srtFile, []byte("test content"), 0644)
		assert.NoError(t, err)

		names, err := lookForAvailableSubtitles(tempDir)

		assert.NoError(t, err)
		assert.Len(t, names, 1)
		// Verify it's just the filename, not the full path
		assert.Equal(t, "test.srt", names[0].Title)
		assert.NotContains(t, names[0].Title, tempDir)
	})

	t.Run("finds srt file in directory with spaces, parentheses and brackets - real world example", func(t *testing.T) {
		// Create a temp base directory
		tempBase, err := os.MkdirTemp("", "base-*")
		assert.NoError(t, err)
		defer os.RemoveAll(tempBase)

		// Create a subdirectory matching the real-world example: "Inside Llewyn Davis (2013) [imdbid-tt2042568]"
		subDir := filepath.Join(tempBase, "Inside Llewyn Davis (2013) [imdbid-tt2042568]")
		err = os.MkdirAll(subDir, 0755)
		assert.NoError(t, err)

		// Create the exact .srt file from the real-world example
		srtFile := filepath.Join(subDir, "Inside.Llewyn.Davis.2013.1080p.BluRay.x264.YIFY.srt")
		err = os.WriteFile(srtFile, []byte("test content"), 0644)
		assert.NoError(t, err)

		names, err := lookForAvailableSubtitles(subDir)

		assert.NoError(t, err)
		assert.Len(t, names, 1, "Expected to find 1 subtitle file")
		if len(names) > 0 {
			assert.Equal(t, "Inside.Llewyn.Davis.2013.1080p.BluRay.x264.YIFY.srt", names[0].Title)
		}
	})
}
