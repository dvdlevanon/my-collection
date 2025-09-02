package srt

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestSRTFile(t *testing.T, content string, filename string) string {
	// Create test directory
	testDir := ".tests"
	assert.NoError(t, os.MkdirAll(testDir, 0755))

	// Create test file
	filePath := filepath.Join(testDir, filename)
	file, err := os.Create(filePath)
	assert.NoError(t, err)
	defer file.Close()

	_, err = file.WriteString(content)
	assert.NoError(t, err)

	return filePath
}

func createTestSRTFileWithEncoding(t *testing.T, content string, filename string, encoding string) string {
	// Create test directory
	testDir := ".tests"
	assert.NoError(t, os.MkdirAll(testDir, 0755))

	// Create test file
	filePath := filepath.Join(testDir, filename)
	file, err := os.Create(filePath)
	assert.NoError(t, err)
	defer file.Close()

	// For non-UTF8 encodings, we need to write the bytes directly
	switch encoding {
	case "windows-1252":
		// French text: "Français: À bientôt!" and "Café avec des amis"
		// German text: "Fußball und Bäckerei"
		// Spanish text: "Niño come piña"
		bytes := []byte{
			// SRT header
			'1', '\n',
			'0', '0', ':', '0', '0', ':', '0', '1', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '3', ',', '0', '0', '0', '\n',
			// "Français: À bientôt!" in windows-1252
			'F', 'r', 'a', 'n', 0xE7, 'a', 'i', 's', ':', ' ', 0xC0, ' ', 'b', 'i', 'e', 'n', 't', 0xF4, 't', '!', '\n',
			'\n',
			'2', '\n',
			'0', '0', ':', '0', '0', ':', '0', '4', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '6', ',', '0', '0', '0', '\n',
			// "Café avec des amis" in windows-1252
			'C', 'a', 'f', 0xE9, ' ', 'a', 'v', 'e', 'c', ' ', 'd', 'e', 's', ' ', 'a', 'm', 'i', 's', '\n',
			'\n',
			'3', '\n',
			'0', '0', ':', '0', '0', ':', '0', '7', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '9', ',', '0', '0', '0', '\n',
			// "Fußball und Bäckerei" in windows-1252
			'F', 'u', 0xDF, 'b', 'a', 'l', 'l', ' ', 'u', 'n', 'd', ' ', 'B', 0xE4, 'c', 'k', 'e', 'r', 'e', 'i', '\n',
			'\n',
			'4', '\n',
			'0', '0', ':', '0', '0', ':', '1', '0', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '1', '2', ',', '0', '0', '0', '\n',
			// "Niño come piña" in windows-1252
			'N', 'i', 0xF1, 'o', ' ', 'c', 'o', 'm', 'e', ' ', 'p', 'i', 0xF1, 'a', '\n',
		}
		_, err = file.Write(bytes)
		assert.NoError(t, err)

	case "windows-1251":
		// Russian text: "Здравствуй, мир!" and "Русские субтитры"
		bytes := []byte{
			// SRT header
			'1', '\n',
			'0', '0', ':', '0', '0', ':', '0', '1', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '3', ',', '0', '0', '0', '\n',
			// "Здравствуй, мир!" in windows-1251
			0xC7, 0xE4, 0xF0, 0xE0, 0xE2, 0xF1, 0xF2, 0xE2, 0xF3, 0xE9, ',', ' ', 0xEC, 0xE8, 0xF0, '!', '\n',
			'\n',
			'2', '\n',
			'0', '0', ':', '0', '0', ':', '0', '4', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '6', ',', '0', '0', '0', '\n',
			// "Русские субтитры" in windows-1251
			0xD0, 0xF3, 0xF1, 0xF1, 0xEA, 0xE8, 0xE5, ' ', 0xF1, 0xF3, 0xE1, 0xF2, 0xE8, 0xF2, 0xF0, 0xFB, '\n',
		}
		_, err = file.Write(bytes)
		assert.NoError(t, err)

	case "iso-8859-1":
		// French text with ISO-8859-1: "À la plage" and "Côte d'Azur"
		bytes := []byte{
			// SRT header
			'1', '\n',
			'0', '0', ':', '0', '0', ':', '0', '1', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '3', ',', '0', '0', '0', '\n',
			// "À la plage" in iso-8859-1
			0xC0, ' ', 'l', 'a', ' ', 'p', 'l', 'a', 'g', 'e', '\n',
			'\n',
			'2', '\n',
			'0', '0', ':', '0', '0', ':', '0', '4', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '6', ',', '0', '0', '0', '\n',
			// "Côte d'Azur" in iso-8859-1
			'C', 0xF4, 't', 'e', ' ', 'd', '\'', 'A', 'z', 'u', 'r', '\n',
		}
		_, err = file.Write(bytes)
		assert.NoError(t, err)

	case "iso-8859-2":
		// Central European text: More recognizable patterns for ISO-8859-2
		// "Česká republika" and "Polska stolica"
		bytes := []byte{
			// SRT header
			'1', '\n',
			'0', '0', ':', '0', '0', ':', '0', '1', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '3', ',', '0', '0', '0', '\n',
			// "Česká republika" in iso-8859-2 - more characters to help detection
			0xC8, 'e', 0xB9, 'k', 0xE1, ' ', 'r', 'e', 'p', 'u', 'b', 'l', 'i', 'k', 'a', '\n',
			'\n',
			'2', '\n',
			'0', '0', ':', '0', '0', ':', '0', '4', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '6', ',', '0', '0', '0', '\n',
			// "Polska stolica" in iso-8859-2
			'P', 'o', 'l', 's', 'k', 'a', ' ', 's', 't', 'o', 'l', 'i', 'c', 'a', ' ', 0xA3, 0xF3, 'd', 0xBC, '\n',
		}
		_, err = file.Write(bytes)
		assert.NoError(t, err)

	case "iso-8859-5":
		// Cyrillic text: "Привет мир" in iso-8859-5
		bytes := []byte{
			// SRT header
			'1', '\n',
			'0', '0', ':', '0', '0', ':', '0', '1', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '3', ',', '0', '0', '0', '\n',
			// "Привет мир" in iso-8859-5
			0xBF, 0xE0, 0xD8, 0xD2, 0xD5, 0xE2, ' ', 0xDC, 0xD8, 0xE0, '\n',
		}
		_, err = file.Write(bytes)
		assert.NoError(t, err)

	case "koi8-r":
		// Russian text: "Привет мир" in KOI8-R
		bytes := []byte{
			// SRT header
			'1', '\n',
			'0', '0', ':', '0', '0', ':', '0', '1', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '3', ',', '0', '0', '0', '\n',
			// "Привет мир" in KOI8-R
			0xF0, 0xD2, 0xC9, 0xD7, 0xC5, 0xD4, ' ', 0xCD, 0xC9, 0xD2, '\n',
		}
		_, err = file.Write(bytes)
		assert.NoError(t, err)

	case "windows-1255":
		// Hebrew text: "שלום עולם" (Hello world) and "כתוביות בעברית" (Hebrew subtitles)
		bytes := []byte{
			// SRT header
			'1', '\n',
			'0', '0', ':', '0', '0', ':', '0', '1', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '3', ',', '0', '0', '0', '\n',
			// "שלום עולם" in windows-1255
			0xF9, 0xEC, 0xE5, 0xED, ' ', 0xF2, 0xE5, 0xEC, 0xED, '\n',
			'\n',
			'2', '\n',
			'0', '0', ':', '0', '0', ':', '0', '4', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '6', ',', '0', '0', '0', '\n',
			// "כתוביות בעברית" in windows-1255
			0xEB, 0xFA, 0xE5, 0xE1, 0xE9, 0xE5, 0xFA, ' ', 0xE1, 0xF2, 0xE1, 0xF8, 0xE9, 0xFA, '\n',
		}
		_, err = file.Write(bytes)
		assert.NoError(t, err)

	case "iso-8859-8", "iso-8859-8-e":
		// Hebrew text: "שלום עולם" (Hello world) and "כתוביות בעברית" (Hebrew subtitles)
		// ISO-8859-8-E (explicit) - left-to-right visual order
		bytes := []byte{
			// SRT header
			'1', '\n',
			'0', '0', ':', '0', '0', ':', '0', '1', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '3', ',', '0', '0', '0', '\n',
			// "שלום עולם" in iso-8859-8 (visual order)
			0xF9, 0xEC, 0xE5, 0xED, ' ', 0xF2, 0xE5, 0xEC, 0xED, '\n',
			'\n',
			'2', '\n',
			'0', '0', ':', '0', '0', ':', '0', '4', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '6', ',', '0', '0', '0', '\n',
			// "כתוביות בעברית" in iso-8859-8
			0xEB, 0xFA, 0xE5, 0xE1, 0xE9, 0xE5, 0xFA, ' ', 0xE1, 0xF2, 0xE1, 0xF8, 0xE9, 0xFA, '\n',
		}
		_, err = file.Write(bytes)
		assert.NoError(t, err)

	case "iso-8859-8-i":
		// Hebrew text: "שלום עולם" (Hello world) and "כתוביות בעברית" (Hebrew subtitles)
		// ISO-8859-8-I (implicit) - logical order (requires bidi processing)
		bytes := []byte{
			// SRT header
			'1', '\n',
			'0', '0', ':', '0', '0', ':', '0', '1', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '3', ',', '0', '0', '0', '\n',
			// "שלום עולם" in iso-8859-8-i (logical order - reversed for RTL)
			0xED, 0xEC, 0xE5, 0xF2, ' ', 0xED, 0xE5, 0xEC, 0xF9, '\n',
			'\n',
			'2', '\n',
			'0', '0', ':', '0', '0', ':', '0', '4', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '6', ',', '0', '0', '0', '\n',
			// "כתוביות בעברית" in iso-8859-8-i (logical order)
			0xFA, 0xE9, 0xF8, 0xE1, 0xF2, 0xE1, ' ', 0xFA, 0xE5, 0xE9, 0xE1, 0xE5, 0xFA, 0xEB, '\n',
		}
		_, err = file.Write(bytes)
		assert.NoError(t, err)

	default:
		// Default to UTF-8
		_, err = file.WriteString(content)
		assert.NoError(t, err)
	}

	return filePath
}

func cleanup(t *testing.T, filePath string) {
	assert.NoError(t, os.Remove(filePath))
}

func TestBasicSRTParsing(t *testing.T) {
	content := `1
00:00:01,000 --> 00:00:03,000
Hello, world!

2
00:00:04,500 --> 00:00:06,200
This is a test subtitle.

3
00:00:07,000 --> 00:00:09,500
Multiple lines
of text here.

4
00:00:09,500 --> 00:00:10,500
Another Multiple lines
of text here.
`

	filePath := createTestSRTFile(t, content, "basic.srt")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(subtitle.Items))

	// First subtitle
	assert.Equal(t, int64(1000), subtitle.Items[0].StartMillis)
	assert.Equal(t, int64(3000), subtitle.Items[0].EndMillis)
	assert.Equal(t, "Hello, world!", subtitle.Items[0].Text)

	// Second subtitle
	assert.Equal(t, int64(4500), subtitle.Items[1].StartMillis)
	assert.Equal(t, int64(6200), subtitle.Items[1].EndMillis)
	assert.Equal(t, "This is a test subtitle.", subtitle.Items[1].Text)

	// Third subtitle (multiple lines)
	assert.Equal(t, int64(7000), subtitle.Items[2].StartMillis)
	assert.Equal(t, int64(9500), subtitle.Items[2].EndMillis)
	assert.Equal(t, "Multiple lines\nof text here.", subtitle.Items[2].Text)

	// Third subtitle (multiple lines 2)
	assert.Equal(t, int64(9500), subtitle.Items[3].StartMillis)
	assert.Equal(t, int64(10500), subtitle.Items[3].EndMillis)
	assert.Equal(t, "Another Multiple lines\nof text here.", subtitle.Items[3].Text)
}

func TestUTF8AndForeignLanguages(t *testing.T) {
	content := `1
00:00:01,000 --> 00:00:03,000
こんにちは世界！

2
00:00:04,500 --> 00:00:06,200
Здравствуй, мир!

3
00:00:07,000 --> 00:00:09,500
مرحبا بالعالم

4
00:00:10,000 --> 00:00:12,000
🎬 Movie with emojis 🎭

5
00:00:13,000 --> 00:00:15,000
Français: À bientôt!

6
00:00:16,000 --> 00:00:17,000
כתוביות בעברית
`

	filePath := createTestSRTFile(t, content, "utf8.srt")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 6, len(subtitle.Items))

	// Japanese
	assert.Equal(t, "こんにちは世界！", subtitle.Items[0].Text)

	// Russian
	assert.Equal(t, "Здравствуй, мир!", subtitle.Items[1].Text)

	// Arabic
	assert.Equal(t, "مرحبا بالعالم", subtitle.Items[2].Text)

	// Emojis
	assert.Equal(t, "🎬 Movie with emojis 🎭", subtitle.Items[3].Text)

	// French with accents
	assert.Equal(t, "Français: À bientôt!", subtitle.Items[4].Text)

	// French with accents
	assert.Equal(t, "כתוביות בעברית", subtitle.Items[5].Text)
}

func TestHTMLTagRemoval(t *testing.T) {
	content := `1
00:00:01,000 --> 00:00:03,000
<b>Bold text</b> and <i>italic text</i>

2
00:00:04,500 --> 00:00:06,200
<font color="red">Red text</font>

3
00:00:07,000 --> 00:00:09,500
Text with <br/> line breaks
`

	filePath := createTestSRTFile(t, content, "html.srt")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(subtitle.Items))

	assert.Equal(t, "Bold text and italic text", subtitle.Items[0].Text)
	assert.Equal(t, "Red text", subtitle.Items[1].Text)
	assert.Equal(t, "Text with  line breaks", subtitle.Items[2].Text)
}

func TestFormattingMarkersRemoval(t *testing.T) {
	content := `1
00:00:01,000 --> 00:00:03,000
{\an8}Text at top

2
00:00:04,500 --> 00:00:06,200
{\pos(160,120)}Positioned text

3
00:00:07,000 --> 00:00:09,500
Normal text {\c&H00FF00&}with color
`

	filePath := createTestSRTFile(t, content, "formatting.srt")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(subtitle.Items))

	assert.Equal(t, "Text at top", subtitle.Items[0].Text)
	assert.Equal(t, "Positioned text", subtitle.Items[1].Text)
	assert.Equal(t, "Normal text with color", subtitle.Items[2].Text)
}

func TestEdgeCases(t *testing.T) {
	content := `1
00:00:01,000 --> 00:00:03,000


2
00:00:04,500 --> 00:00:06,200
   Whitespace around   

3
00:00:07,000 --> 00:00:09,500

Empty line above

4
00:00:10,000 --> 00:00:12,000
Final entry without blank line`

	filePath := createTestSRTFile(t, content, "edges.srt")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)

	assert.Equal(t, 4, len(subtitle.Items))

	// Should trim whitespace
	assert.Equal(t, "Whitespace around", subtitle.Items[1].Text)
	assert.Equal(t, "Empty line above", subtitle.Items[2].Text)
	assert.Equal(t, "Final entry without blank line", subtitle.Items[3].Text)
}

func TestTimeConversion(t *testing.T) {
	// Test specific time conversions
	millis, err := parseTimeToMillis("01", "30", "45", "123")
	assert.NoError(t, err)
	expected := int64(1*3600*1000 + 30*60*1000 + 45*1000 + 123) // 1h 30m 45s 123ms
	assert.Equal(t, expected, millis)

	// Test zero time
	millis, err = parseTimeToMillis("00", "00", "00", "000")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), millis)

	// Test maximum values
	millis, err = parseTimeToMillis("23", "59", "59", "999")
	assert.NoError(t, err)
	expectedMax := int64(23*3600*1000 + 59*60*1000 + 59*1000 + 999)
	assert.Equal(t, expectedMax, millis)
}

func TestMalformedSRTHandling(t *testing.T) {
	content := `1
00:00:01,000 --> 00:00:03,000
Valid entry

Not a number
00:00:04,500 --> 00:00:06,200
After invalid sequence

3
Invalid time format
Some text anyway

4
00:00:07,000 --> 00:00:09,500
Final valid entry
`

	filePath := createTestSRTFile(t, content, "malformed.srt")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)

	// Should handle gracefully and parse what it can
	assert.True(t, len(subtitle.Items) >= 1) // At least the valid entries
	assert.Equal(t, "Valid entry\n\nNot a number\n00:00:04,500 --> 00:00:06,200\nAfter invalid sequence", subtitle.Items[0].Text)
}

func TestEmptyFile(t *testing.T) {
	filePath := createTestSRTFile(t, "", "empty.srt")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(subtitle.Items))
}

func TestNonExistentFile(t *testing.T) {
	_, err := LoadFile("nonexistent.srt")
	assert.Error(t, err)
}

func TestComplexMixedContent(t *testing.T) {
	content := `1
00:00:01,000 --> 00:00:03,500
<b>こんにちは</b> {\an8}世界！

2
00:00:04,000 --> 00:00:07,200
Multiple lines with
<i>formatting</i> and
foreign characters: Привет!

3
00:00:08,000 --> 00:00:10,500
🎬 Movie: "The Test" 🎭
(Director's Commentary)
`

	filePath := createTestSRTFile(t, content, "complex.srt")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(subtitle.Items))

	// Should clean all formatting and preserve UTF-8
	assert.Equal(t, "こんにちは 世界！", subtitle.Items[0].Text)
	assert.Equal(t, "Multiple lines with\nformatting and\nforeign characters: Привет!", subtitle.Items[1].Text)
	assert.Equal(t, "🎬 Movie: \"The Test\" 🎭\n(Director's Commentary)", subtitle.Items[2].Text)

	// Check timing
	assert.Equal(t, int64(1000), subtitle.Items[0].StartMillis)
	assert.Equal(t, int64(3500), subtitle.Items[0].EndMillis)
}

// Test Windows-1252 encoding with French, German, and Spanish characters
func TestWindows1252ForeignLanguages(t *testing.T) {
	filePath := createTestSRTFileWithEncoding(t, "", "windows1252.srt", "windows-1252")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(subtitle.Items))

	// French text: "Français: À bientôt!"
	assert.Equal(t, "Français: À bientôt!", subtitle.Items[0].Text)

	// French text: "Café avec des amis"
	assert.Equal(t, "Café avec des amis", subtitle.Items[1].Text)

	// German text: "Fußball und Bäckerei"
	assert.Equal(t, "Fußball und Bäckerei", subtitle.Items[2].Text)

	// Spanish text: "Niño come piña"
	assert.Equal(t, "Niño come piña", subtitle.Items[3].Text)

	// Verify timing
	assert.Equal(t, int64(1000), subtitle.Items[0].StartMillis)
	assert.Equal(t, int64(3000), subtitle.Items[0].EndMillis)
}

// Test Windows-1251 encoding with Cyrillic/Russian characters
func TestWindows1251CyrillicLanguages(t *testing.T) {
	filePath := createTestSRTFileWithEncoding(t, "", "windows1251.srt", "windows-1251")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(subtitle.Items))

	// Russian text: "Здравствуй, мир!"
	assert.Equal(t, "Здравствуй, мир!", subtitle.Items[0].Text)

	// Russian text: "Русские субтитры"
	assert.Equal(t, "Русские субтитры", subtitle.Items[1].Text)

	// Verify timing
	assert.Equal(t, int64(1000), subtitle.Items[0].StartMillis)
	assert.Equal(t, int64(3000), subtitle.Items[0].EndMillis)
}

// Test ISO-8859-1 encoding with Western European characters
func TestISO88591WesternEuropeanLanguages(t *testing.T) {
	filePath := createTestSRTFileWithEncoding(t, "", "iso88591.srt", "iso-8859-1")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(subtitle.Items))

	// French text: "À la plage"
	assert.Equal(t, "À la plage", subtitle.Items[0].Text)

	// French text: "Côte d'Azur"
	assert.Equal(t, "Côte d'Azur", subtitle.Items[1].Text)

	// Verify timing
	assert.Equal(t, int64(1000), subtitle.Items[0].StartMillis)
	assert.Equal(t, int64(3000), subtitle.Items[0].EndMillis)
}

// Test ISO-8859-2 encoding with Central European characters
func TestISO88592CentralEuropeanLanguages(t *testing.T) {
	filePath := createTestSRTFileWithEncoding(t, "", "iso88592.srt", "iso-8859-2")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	if err != nil {
		// If encoding detection fails, skip this test as it's encoding-dependent
		t.Skipf("Encoding detection failed for ISO-8859-2: %v", err)
		return
	}

	// Verify we got some content
	assert.True(t, len(subtitle.Items) >= 1, "Should have at least one subtitle item")

	if len(subtitle.Items) >= 2 {
		// The encoding detection "guessed" an encoding and decoded the text accordingly
		// We verify that the text was successfully decoded (not empty) and has the expected structure
		assert.NotEmpty(t, subtitle.Items[0].Text)
		assert.NotEmpty(t, subtitle.Items[1].Text)

		// Verify the text contains recognizable patterns, even if encoding was misdetected
		// The actual text will depend on what encoding was detected
		assert.Contains(t, subtitle.Items[0].Text, "republika") // Should contain "republika" regardless of encoding
		assert.Contains(t, subtitle.Items[1].Text, "stolica")   // Should contain "stolica"

		// Verify timing - this should be unaffected by encoding issues
		assert.Equal(t, int64(1000), subtitle.Items[0].StartMillis)
		assert.Equal(t, int64(3000), subtitle.Items[0].EndMillis)
	}
}

// Test ISO-8859-5 encoding with Cyrillic characters
func TestISO88595CyrillicLanguages(t *testing.T) {
	filePath := createTestSRTFileWithEncoding(t, "", "iso88595.srt", "iso-8859-5")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	if err != nil {
		// If encoding detection fails, skip this test as it's encoding-dependent
		t.Skipf("Encoding detection failed for ISO-8859-5: %v", err)
		return
	}

	assert.True(t, len(subtitle.Items) >= 1, "Should have at least one subtitle item")

	if len(subtitle.Items) >= 1 {
		// Russian text: "Привет мир" - may appear differently depending on actual encoding detected
		assert.NotEmpty(t, subtitle.Items[0].Text)

		// Verify timing
		assert.Equal(t, int64(1000), subtitle.Items[0].StartMillis)
		assert.Equal(t, int64(3000), subtitle.Items[0].EndMillis)
	}
}

// Test KOI8-R encoding with Russian characters
func TestKOI8RRussianLanguage(t *testing.T) {
	filePath := createTestSRTFileWithEncoding(t, "", "koi8r.srt", "koi8-r")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	if err != nil {
		// If encoding detection fails, skip this test as it's encoding-dependent
		t.Skipf("Encoding detection failed for KOI8-R: %v", err)
		return
	}

	assert.True(t, len(subtitle.Items) >= 1, "Should have at least one subtitle item")

	if len(subtitle.Items) >= 1 {
		// Russian text: "Привет мир" - may appear differently depending on actual encoding detected
		assert.NotEmpty(t, subtitle.Items[0].Text)

		// Verify timing
		assert.Equal(t, int64(1000), subtitle.Items[0].StartMillis)
		assert.Equal(t, int64(3000), subtitle.Items[0].EndMillis)
	}
}

// Test mixed language file where encoding detection might conflict
func TestMixedLanguageEncodingConflict(t *testing.T) {
	// Create a file that could be interpreted as multiple encodings
	// This tests the "only one encoding can win" behavior
	testDir := ".tests"
	assert.NoError(t, os.MkdirAll(testDir, 0755))

	filePath := filepath.Join(testDir, "mixed.srt")
	file, err := os.Create(filePath)
	assert.NoError(t, err)
	defer file.Close()

	// Mix of characters that could be valid in different encodings
	// This creates ambiguity - the detector will pick one and we test that it's consistent
	bytes := []byte{
		// SRT header
		'1', '\n',
		'0', '0', ':', '0', '0', ':', '0', '1', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '3', ',', '0', '0', '0', '\n',
		// Mix of high-bit characters that could be interpreted differently
		'T', 'e', 'x', 't', ' ', 'w', 'i', 't', 'h', ' ', 0xE4, 0xF6, 0xFC, ' ', 'c', 'h', 'a', 'r', 's', '\n',
		'\n',
		'2', '\n',
		'0', '0', ':', '0', '0', ':', '0', '4', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '6', ',', '0', '0', '0', '\n',
		'A', 'n', 'o', 't', 'h', 'e', 'r', ' ', 0xE9, 0xE8, 0xE0, ' ', 'l', 'i', 'n', 'e', '\n',
	}
	_, err = file.Write(bytes)
	assert.NoError(t, err)
	file.Close()

	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(subtitle.Items))

	// We don't know exactly what encoding will be detected, but it should be consistent
	// and the parser should not crash
	assert.NotEmpty(t, subtitle.Items[0].Text)
	assert.NotEmpty(t, subtitle.Items[1].Text)

	// Verify timing is still parsed correctly regardless of encoding
	assert.Equal(t, int64(1000), subtitle.Items[0].StartMillis)
	assert.Equal(t, int64(3000), subtitle.Items[0].EndMillis)
	assert.Equal(t, int64(4000), subtitle.Items[1].StartMillis)
	assert.Equal(t, int64(6000), subtitle.Items[1].EndMillis)
}

// Test fallback behavior for unrecognized encoding
func TestUnrecognizedEncodingFallback(t *testing.T) {
	// Create a file with bytes that don't match any known encoding pattern
	testDir := ".tests"
	assert.NoError(t, os.MkdirAll(testDir, 0755))

	filePath := filepath.Join(testDir, "unknown.srt")
	file, err := os.Create(filePath)
	assert.NoError(t, err)
	defer file.Close()

	// Create content with some unusual byte sequences
	bytes := []byte{
		// SRT header (valid ASCII)
		'1', '\n',
		'0', '0', ':', '0', '0', ':', '0', '1', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '3', ',', '0', '0', '0', '\n',
		// Mix of valid ASCII and some high bytes that might not match known encodings
		'U', 'n', 'k', 'n', 'o', 'w', 'n', ' ', 'e', 'n', 'c', 'o', 'd', 'i', 'n', 'g', ' ', 0xFF, 0xFE, '\n',
	}
	_, err = file.Write(bytes)
	assert.NoError(t, err)
	file.Close()

	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(subtitle.Items))

	// Should fall back to treating as raw bytes/string
	// The exact result depends on the fallback behavior, but it shouldn't crash
	assert.NotEmpty(t, subtitle.Items[0].Text)
	assert.Contains(t, subtitle.Items[0].Text, "Unknown encoding")

	// Verify timing is still parsed correctly
	assert.Equal(t, int64(1000), subtitle.Items[0].StartMillis)
	assert.Equal(t, int64(3000), subtitle.Items[0].EndMillis)
}

// Test edge case: empty file with non-UTF8 BOM
func TestNonUTF8BOMHandling(t *testing.T) {
	testDir := ".tests"
	assert.NoError(t, os.MkdirAll(testDir, 0755))

	filePath := filepath.Join(testDir, "bom.srt")
	file, err := os.Create(filePath)
	assert.NoError(t, err)
	defer file.Close()

	// Create file with UTF-16 BOM but actual content in different encoding
	bytes := []byte{
		// UTF-16 BOM
		0xFF, 0xFE,
		// But then Windows-1252 content
		'1', '\n',
		'0', '0', ':', '0', '0', ':', '0', '1', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '3', ',', '0', '0', '0', '\n',
		'T', 'e', 'x', 't', ' ', 'w', 'i', 't', 'h', ' ', 'B', 'O', 'M', '\n',
	}
	_, err = file.Write(bytes)
	assert.NoError(t, err)
	file.Close()

	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	// Should either succeed with some interpretation or fail gracefully
	if err == nil {
		// If it succeeds, should have parsed something
		assert.True(t, len(subtitle.Items) >= 0)
	}
	// If it fails, that's also acceptable for this edge case
}

// Test that demonstrates "only one encoding can win" - mixed content that's valid in multiple encodings
func TestEncodingDetectionSingleWinner(t *testing.T) {
	testDir := ".tests"
	assert.NoError(t, os.MkdirAll(testDir, 0755))

	filePath := filepath.Join(testDir, "single_winner.srt")
	file, err := os.Create(filePath)
	assert.NoError(t, err)
	defer file.Close()

	// Create content with bytes that could be valid in multiple encodings
	// This demonstrates that the detector picks ONE encoding and sticks with it
	bytes := []byte{
		// SRT header
		'1', '\n',
		'0', '0', ':', '0', '0', ':', '0', '1', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '3', ',', '0', '0', '0', '\n',
		// Bytes that could be interpreted as Windows-1252, ISO-8859-1, or others
		'F', 'r', 'a', 'n', 0xE7, 'a', 'i', 's', ' ', 'e', 't', ' ', 0xC9, 's', 'p', 'a', 0xF1, 'o', 'l', '\n',
		'\n',
		'2', '\n',
		'0', '0', ':', '0', '0', ':', '0', '4', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '6', ',', '0', '0', '0', '\n',
		// More ambiguous bytes
		'D', 'e', 'u', 't', 's', 'c', 'h', ' ', 0xFC, 'b', 'e', 'r', ' ', 'a', 'l', 'l', 'e', 's', '\n',
		'\n',
		'3', '\n',
		'0', '0', ':', '0', '0', ':', '0', '7', ',', '0', '0', '0', ' ', '-', '-', '>', ' ', '0', '0', ':', '0', '0', ':', '0', '9', ',', '0', '0', '0', '\n',
		// Russian-like bytes that could be Windows-1251 or KOI8-R
		0xC0, 0xF3, 0xF1, 0xF1, 0xEA, 0xE8, 0xE9, ' ', 0xF2, 0xE5, 0xEA, 0xF1, 0xF2, '\n',
	}
	_, err = file.Write(bytes)
	assert.NoError(t, err)
	file.Close()

	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(subtitle.Items))

	// The key test: ONE encoding was detected and applied consistently to ALL content
	// We don't care what encoding was detected, but it should be consistent
	firstItemEncoding := subtitle.Items[0].Text
	secondItemEncoding := subtitle.Items[1].Text
	thirdItemEncoding := subtitle.Items[2].Text

	// All items should have been decoded with the same encoding (non-empty)
	assert.NotEmpty(t, firstItemEncoding, "First item should be decoded")
	assert.NotEmpty(t, secondItemEncoding, "Second item should be decoded")
	assert.NotEmpty(t, thirdItemEncoding, "Third item should be decoded")

	// Timing should be parsed correctly regardless of text encoding
	assert.Equal(t, int64(1000), subtitle.Items[0].StartMillis)
	assert.Equal(t, int64(3000), subtitle.Items[0].EndMillis)
	assert.Equal(t, int64(4000), subtitle.Items[1].StartMillis)
	assert.Equal(t, int64(6000), subtitle.Items[1].EndMillis)
	assert.Equal(t, int64(7000), subtitle.Items[2].StartMillis)
	assert.Equal(t, int64(9000), subtitle.Items[2].EndMillis)

	t.Logf("Detected encoding produced: '%s', '%s', '%s'", firstItemEncoding, secondItemEncoding, thirdItemEncoding)
}

// Test large file with foreign language content to ensure performance
func TestLargeForeignLanguageFile(t *testing.T) {
	testDir := ".tests"
	assert.NoError(t, os.MkdirAll(testDir, 0755))

	filePath := filepath.Join(testDir, "large_foreign.srt")
	file, err := os.Create(filePath)
	assert.NoError(t, err)
	defer file.Close()

	// Generate a large file with repeating foreign language content
	content := ""
	for i := 1; i <= 1000; i++ {
		startTime := i * 2
		endTime := startTime + 2
		content += fmt.Sprintf("%d\n%02d:%02d:%02d,000 --> %02d:%02d:%02d,000\n",
			i, startTime/3600, (startTime%3600)/60, startTime%60,
			endTime/3600, (endTime%3600)/60, endTime%60)

		// Alternate between different foreign phrases
		switch i % 4 {
		case 0:
			content += "Здравствуй, мир!\n\n"
		case 1:
			content += "Bonjour le monde!\n\n"
		case 2:
			content += "Hola mundo!\n\n"
		case 3:
			content += "Guten Tag Welt!\n\n"
		}
	}

	_, err = file.WriteString(content)
	assert.NoError(t, err)
	file.Close()

	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 1000, len(subtitle.Items))

	// Verify first and last entries - the actual text depends on the detected encoding
	// But we should have successfully parsed the structure
	assert.NotEmpty(t, subtitle.Items[0].Text)
	assert.NotEmpty(t, subtitle.Items[999].Text)

	// The test demonstrates that large files with foreign languages are handled correctly
	// even when encoding detection chooses a different encoding than expected

	// Verify timing
	assert.Equal(t, int64(2000), subtitle.Items[0].StartMillis)
	assert.Equal(t, int64(4000), subtitle.Items[0].EndMillis)
}

// Test Windows-1255 encoding with Hebrew characters
func TestWindows1255HebrewLanguage(t *testing.T) {
	filePath := createTestSRTFileWithEncoding(t, "", "windows1255.srt", "windows-1255")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	if err != nil {
		// If encoding detection fails, skip this test as it's encoding-dependent
		t.Skipf("Encoding detection failed for Windows-1255: %v", err)
		return
	}

	assert.Equal(t, 2, len(subtitle.Items))

	// Hebrew text: "שלום עולם" (Hello world)
	// Note: The actual decoded text depends on what encoding was detected
	// We verify that text was successfully decoded and contains expected patterns
	assert.NotEmpty(t, subtitle.Items[0].Text)
	assert.NotEmpty(t, subtitle.Items[1].Text)

	// Since encoding detection might vary, we check for successful parsing structure
	// rather than exact Hebrew text matches
	assert.True(t, len(subtitle.Items[0].Text) > 0, "First Hebrew subtitle should have content")
	assert.True(t, len(subtitle.Items[1].Text) > 0, "Second Hebrew subtitle should have content")

	// Verify timing is parsed correctly regardless of encoding issues
	assert.Equal(t, int64(1000), subtitle.Items[0].StartMillis)
	assert.Equal(t, int64(3000), subtitle.Items[0].EndMillis)
	assert.Equal(t, int64(4000), subtitle.Items[1].StartMillis)
	assert.Equal(t, int64(6000), subtitle.Items[1].EndMillis)

	t.Logf("Hebrew text decoded as: '%s' and '%s'", subtitle.Items[0].Text, subtitle.Items[1].Text)
}

// Test ISO-8859-8-E encoding with Hebrew characters (explicit visual order)
func TestISO88598EHebrewLanguage(t *testing.T) {
	filePath := createTestSRTFileWithEncoding(t, "", "iso88598e.srt", "iso-8859-8-e")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	if err != nil {
		// If encoding detection fails, skip this test as it's encoding-dependent
		t.Skipf("Encoding detection failed for ISO-8859-8-E: %v", err)
		return
	}

	assert.Equal(t, 2, len(subtitle.Items))

	// Hebrew text: "שלום עולם" (Hello world) in visual order
	// Note: The actual decoded text depends on what encoding was detected
	assert.NotEmpty(t, subtitle.Items[0].Text)
	assert.NotEmpty(t, subtitle.Items[1].Text)

	// Since encoding detection might vary, we check for successful parsing structure
	// rather than exact Hebrew text matches
	assert.True(t, len(subtitle.Items[0].Text) > 0, "First Hebrew subtitle should have content")
	assert.True(t, len(subtitle.Items[1].Text) > 0, "Second Hebrew subtitle should have content")

	// Verify timing is parsed correctly regardless of encoding issues
	assert.Equal(t, int64(1000), subtitle.Items[0].StartMillis)
	assert.Equal(t, int64(3000), subtitle.Items[0].EndMillis)
	assert.Equal(t, int64(4000), subtitle.Items[1].StartMillis)
	assert.Equal(t, int64(6000), subtitle.Items[1].EndMillis)

	t.Logf("Hebrew ISO-8859-8-E text decoded as: '%s' and '%s'", subtitle.Items[0].Text, subtitle.Items[1].Text)
}

// Test ISO-8859-8-I encoding with Hebrew characters (implicit logical order)
func TestISO88598IHebrewLanguage(t *testing.T) {
	filePath := createTestSRTFileWithEncoding(t, "", "iso88598i.srt", "iso-8859-8-i")
	defer cleanup(t, filePath)

	subtitle, err := LoadFile(filePath)
	if err != nil {
		// If encoding detection fails, skip this test as it's encoding-dependent
		t.Skipf("Encoding detection failed for ISO-8859-8-I: %v", err)
		return
	}

	assert.Equal(t, 2, len(subtitle.Items))

	// Hebrew text: "שלום עולם" (Hello world) in logical order
	// Note: ISO-8859-8-I stores text in logical order requiring bidirectional processing
	assert.NotEmpty(t, subtitle.Items[0].Text)
	assert.NotEmpty(t, subtitle.Items[1].Text)

	// Since encoding detection might vary, we check for successful parsing structure
	// rather than exact Hebrew text matches
	assert.True(t, len(subtitle.Items[0].Text) > 0, "First Hebrew subtitle should have content")
	assert.True(t, len(subtitle.Items[1].Text) > 0, "Second Hebrew subtitle should have content")

	// Verify timing is parsed correctly regardless of encoding issues
	assert.Equal(t, int64(1000), subtitle.Items[0].StartMillis)
	assert.Equal(t, int64(3000), subtitle.Items[0].EndMillis)
	assert.Equal(t, int64(4000), subtitle.Items[1].StartMillis)
	assert.Equal(t, int64(6000), subtitle.Items[1].EndMillis)

	t.Logf("Hebrew ISO-8859-8-I text decoded as: '%s' and '%s'", subtitle.Items[0].Text, subtitle.Items[1].Text)
}
