package ffmpeg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/utils"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
)

var logger = utils.MustGetLogger("ffmpeg")

type ffprobeShowFormatOutput struct {
	Format ffprobeFormatOutput `json:"format"`
}

type ffprobeFormatOutput struct {
	Duration string `json:"duration"`
}

type ffprobeShowStreamsOutput struct {
	Streams []FfprobeShowStreamOutput `json:"streams"`
}

type FfprobeShowStreamOutput struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func execute(ctx context.Context, name string, arg ...string) ([]byte, error) {
	logger.Debugf(ctx, "Running %s \"%s\"", name, strings.Join(arg, "\" \""))

	cmd := exec.CommandContext(ctx, name, arg...)
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if err := cmd.Wait(); err != nil {
		return nil, errors.Errorf("Error running process, exit code: %d, err: %s %v",
			cmd.ProcessState.ExitCode(), stderr.String(), err)
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return nil, errors.Errorf("Invalid exit code %d", cmd.ProcessState.ExitCode())
	}

	return stdout.Bytes(), nil
}

func GetDurationInSeconds(ctx context.Context, videoFile string) (float64, error) {
	output, err := execute(ctx, "ffprobe", videoFile, "-show_format", "-v", "quiet", "-print_format", "json")
	if err != nil {
		return 0, err
	}

	format := ffprobeShowFormatOutput{}
	if err = json.Unmarshal(output, &format); err != nil {
		return 0, errors.Wrap(err, 0)
	}

	durationInSeconds, err := strconv.ParseFloat(format.Format.Duration, 64)
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}

	return durationInSeconds, nil
}

func GetVideoMetadata(ctx context.Context, videoFile string) (FfprobeShowStreamOutput, error) {
	output, err := execute(ctx, "ffprobe", "-show_streams", "-print_format", "json", videoFile)
	if err != nil {
		return FfprobeShowStreamOutput{}, err
	}

	showStreams := ffprobeShowStreamsOutput{}
	if err = json.Unmarshal(output, &showStreams); err != nil {
		return FfprobeShowStreamOutput{}, errors.Wrap(err, 0)
	}

	for _, stream := range showStreams.Streams {
		if stream.CodecType == "video" {
			return stream, nil
		}
	}

	return FfprobeShowStreamOutput{}, errors.Errorf("Video stream not found for %s", videoFile)
}

func GetAudioMetadata(ctx context.Context, videoFile string) (FfprobeShowStreamOutput, error) {
	output, err := execute(ctx, "ffprobe", "-show_streams", "-print_format", "json", videoFile)
	if err != nil {
		return FfprobeShowStreamOutput{}, err
	}

	showStreams := ffprobeShowStreamsOutput{}
	if err = json.Unmarshal(output, &showStreams); err != nil {
		return FfprobeShowStreamOutput{}, errors.Wrap(err, 0)
	}

	for _, stream := range showStreams.Streams {
		if stream.CodecType == "audio" {
			return stream, nil
		}
	}

	return FfprobeShowStreamOutput{}, errors.Errorf("Audio stream not found for %s", videoFile)
}

func TakeScreenshot(ctx context.Context, videoFile string, second float64, targetFile string) error {
	_, err := execute(ctx, "ffmpeg", "-y", "-ss", fmt.Sprintf("%f", second), "-i", videoFile, "-vframes", "1", targetFile)
	if err != nil {
		return err
	}

	return nil
}

func CropScreenshot(ctx context.Context, videoFile string, second float64, rect model.RectFloat, targetFile string) error {
	cropFilter := fmt.Sprintf("crop=%f:%f:%f:%f", rect.W, rect.H, rect.X, rect.Y)
	_, err := execute(ctx, "ffmpeg", "-y", "-ss", fmt.Sprintf("%f", second), "-i", videoFile, "-vf", cropFilter, "-vframes", "1", targetFile)
	if err != nil {
		return err
	}

	return nil
}
func ExtractPartOfVideo(ctx context.Context, videoFile string, second float64, duration int, targetFile string) error {
	_, err := execute(ctx, "ffmpeg", "-ss", fmt.Sprintf("%f", second), "-i", videoFile,
		"-t", fmt.Sprintf("%d", duration), "-vcodec", "copy", "-acodec", "aac", "-ac", "4", targetFile)
	if err != nil {
		return err
	}

	return nil
}

func JoinVideoFiles(ctx context.Context, videoFiles []string, targetFile string) error {
	tempFile, err := os.CreateTemp("", "my-collection-join-video-files-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	str := ""
	for _, file := range videoFiles {
		str += "file '" + file + "'\n"
	}

	if _, err := tempFile.Write([]byte(str)); err != nil {
		return err
	}

	_, err = execute(ctx, "ffmpeg", "-y", "-safe", "0", "-f", "concat", "-i", tempFile.Name(), targetFile)
	if err != nil {
		return err
	}

	return nil
}

func OptimizeVideoForPreview(ctx context.Context, videoFile string, tempFile string) error {
	_, err := execute(ctx, "ffmpeg", "-y", "-i", videoFile, "-b:v", "2M", "-an", tempFile)
	if err != nil {
		return err
	}

	return os.Rename(tempFile, videoFile)
}

func ChangeVideoResolution(ctx context.Context, videoFile string, tempFile string, w int, h int) error {
	logger.Infof(ctx, "Changing video resolution for %s to %d:%d", videoFile, w, h)

	_, err := execute(ctx, "ffmpeg", "-i", videoFile, "-vf", fmt.Sprintf("scale=%d:%d", w, h), "-c:a", "copy", tempFile)
	if err != nil {
		return err
	}

	if err := os.Remove(videoFile); err != nil {
		return err
	}

	return os.Rename(tempFile, videoFile)
}
