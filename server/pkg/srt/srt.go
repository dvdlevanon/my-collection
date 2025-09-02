package srt

import (
	"bufio"
	"fmt"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var timeRegex = regexp.MustCompile(`^(\d{2}):(\d{2}):(\d{2}),(\d{3})\s*-->\s*(\d{2}):(\d{2}):(\d{2}),(\d{3})$`)

func parseTimeToMillis(hours, minutes, seconds, millis string) (int64, error) {
	h, err := strconv.Atoi(hours)
	if err != nil {
		return 0, err
	}
	m, err := strconv.Atoi(minutes)
	if err != nil {
		return 0, err
	}
	s, err := strconv.Atoi(seconds)
	if err != nil {
		return 0, err
	}
	ms, err := strconv.Atoi(millis)
	if err != nil {
		return 0, err
	}

	totalMillis := int64(h)*time.Hour.Milliseconds() +
		int64(m)*time.Minute.Milliseconds() +
		int64(s)*time.Second.Milliseconds() +
		int64(ms)

	return totalMillis, nil
}

func cleanText(text string) string {
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	text = htmlTagRegex.ReplaceAllString(text, "")

	formatRegex := regexp.MustCompile(`\{[^}]*\}`)
	text = formatRegex.ReplaceAllString(text, "")

	text = strings.TrimSpace(text)
	return text
}

func isSequenceLine(line string) bool {
	_, err := strconv.Atoi(line)
	return err == nil
}

func parseTimeLine(line string) (bool, []string) {
	matches := timeRegex.FindStringSubmatch(line)
	return len(matches) == 9, matches
}

func isStartOfSubtitle(lines []string, i int) (bool, []string) {
	curLine := strings.TrimSpace(lines[i])
	if !isSequenceLine(curLine) {
		return false, nil
	}
	if i+1 >= len(lines) {
		return false, nil
	}

	timeLine := strings.TrimSpace(lines[i+1])
	return parseTimeLine(timeLine)
}

func isEndOfSubtitle(lines []string, i int) bool {
	curLine := strings.TrimSpace(lines[i])
	nextLine := strings.TrimSpace(lines[i+1])
	return curLine == "" && isSequenceLine(nextLine)
}

func readSubtitle(lines []string, i int) (string, int) {
	textLines := make([]string, 0)
	for i < len(lines) {
		if i+1 >= len(lines) {
			textLines = append(textLines, strings.TrimSpace(lines[i]))
			i++
			break
		}

		if isEndOfSubtitle(lines, i) {
			i++
			break
		}

		textLines = append(textLines, strings.TrimSpace(lines[i]))
		i++
	}

	return cleanText(strings.Join(textLines, "\n")), i
}

func createSubtitle(parsedTimeLine []string, text string) (model.SubtitleItem, error) {
	startMillis, err := parseTimeToMillis(parsedTimeLine[1], parsedTimeLine[2], parsedTimeLine[3], parsedTimeLine[4])
	if err != nil {
		return model.SubtitleItem{}, fmt.Errorf("invalid time %s", parsedTimeLine)
	}

	endMillis, err := parseTimeToMillis(parsedTimeLine[5], parsedTimeLine[6], parsedTimeLine[7], parsedTimeLine[8])
	if err != nil {
		return model.SubtitleItem{}, fmt.Errorf("invalid time %s", parsedTimeLine)
	}

	return model.SubtitleItem{
		StartMillis: startMillis,
		EndMillis:   endMillis,
		Text:        text,
	}, nil
}

func readLines(path string) ([]string, error) {
	content, err := utils.DetectEncodingAndRead(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func LoadFile(path string) (model.Subtitle, error) {
	lines, err := readLines(path)
	if err != nil {
		return model.Subtitle{}, err
	}

	result := model.Subtitle{Items: make([]model.SubtitleItem, 0)}
	i := 0
	for i < len(lines) {
		isStart, parsedTimeLine := isStartOfSubtitle(lines, i)
		if !isStart {
			i = i + 1
			continue
		}

		var subtitleText string
		subtitleText, i = readSubtitle(lines, i+2)

		subtitle, err := createSubtitle(parsedTimeLine, subtitleText)
		if err != nil {
			return result, err
		}

		result.Items = append(result.Items, subtitle)
	}

	return result, nil
}
