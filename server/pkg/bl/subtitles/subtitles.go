package subtitles

import (
	"context"
	"fmt"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/srt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("subtitles")

const onlineSubsDir = ".online-subs"

type SubtitlesLister interface {
	List(imdbId string, lang string, aiTranslated bool) ([]model.SubtitleMetadata, error)
}

type SubtitlesDownloader interface {
	Download(subtitle model.SubtitleMetadata, outputFile string) error
}

var ErrSubtitileNotFound = fmt.Errorf("subtitle not found")

func lookForAvailableSubtitles(dir string) ([]model.SubtitleMetadata, error) {
	names := make([]model.SubtitleMetadata, 0)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".srt") {
			names = append(names, model.SubtitleMetadata{
				Id:    "local",
				Title: entry.Name(),
				Url:   relativasor.GetRelativePath(filepath.Join(dir, entry.Name())),
			})
		}
	}

	return names, nil
}

func GetSubtitle(ctx context.Context, url string) (model.Subtitle, error) {
	return srt.LoadFile(relativasor.GetAbsoluteFile(url))
}

func GetAvailableNames(ctx context.Context, ir model.ItemReader, itemId uint64) ([]model.SubtitleMetadata, error) {
	item, err := ir.GetItem(ctx, itemId)
	if err != nil {
		return nil, err
	}

	videoFile := relativasor.GetAbsoluteFile(item.Url)
	videoDir := filepath.Dir(videoFile)
	return lookForAvailableSubtitles(videoDir)
}

func extractIMDbID(path string) string {
	re := regexp.MustCompile(`\[imdbid-(tt\d+)\]`)
	matches := re.FindStringSubmatch(path)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

func GetOnlineNames(ctx context.Context, ir model.ItemReader, l SubtitlesLister, itemId uint64, lang string, aiTranslated bool) ([]model.SubtitleMetadata, error) {
	item, err := ir.GetItem(ctx, itemId)
	if err != nil {
		return nil, err
	}

	imdbId := extractIMDbID(item.Url)

	subtitles, err := l.List(imdbId, lang, aiTranslated)
	if err != nil {
		return nil, err
	}

	subtitles = addUrls(item, subtitles)
	return subtitles, nil
}

func getDownloadedSubUrl(item *model.Item, subtitle model.SubtitleMetadata) string {
	videoFile := relativasor.GetAbsoluteFile(item.Url)
	videoDir := filepath.Dir(videoFile)

	return filepath.Join(videoDir, onlineSubsDir, subtitle.Id, fmt.Sprintf("%s.srt", subtitle.Title))
}

func addUrls(item *model.Item, subtitles []model.SubtitleMetadata) []model.SubtitleMetadata {
	for i := range subtitles {
		url := getDownloadedSubUrl(item, subtitles[i])

		_, err := os.Stat(url)
		if err == nil {
			subtitles[i].Url = relativasor.GetRelativePath(url)
		}
	}

	return subtitles
}

func Download(ctx context.Context, ir model.ItemReader, d SubtitlesDownloader, tp model.TempFileProvider, itemId uint64, subtitle model.SubtitleMetadata) error {
	item, err := ir.GetItem(ctx, itemId)
	if err != nil {
		return err
	}

	tempFile := tp.GetTempFile() // TODO: In case of error the temp file stays
	err = d.Download(subtitle, tempFile)
	if err != nil {
		return err
	}

	targetFile := getDownloadedSubUrl(item, subtitle)
	if err := os.MkdirAll(path.Dir(targetFile), 0755); err != nil {
		return err
	}

	return os.Rename(tempFile, targetFile)
}

func Delete(ctx context.Context, url string) error {
	subtitleFile := relativasor.GetAbsoluteFile(url)
	return os.Remove(subtitleFile)
}
