package ffmpeg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"my-collection/server/pkg/model"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("ffmpeg")

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

func execute(name string, arg ...string) ([]byte, error) {
	logger.Debugf("Running %s \"%s\"", name, strings.Join(arg, "\" \""))

	cmd := exec.Command(name, arg...)
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

func GetDurationInSeconds(videoFile string) (float64, error) {
	output, err := execute("ffprobe", videoFile, "-show_format", "-v", "quiet", "-print_format", "json")
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

func GetVideoMetadata(videoFile string) (FfprobeShowStreamOutput, error) {
	output, err := execute("ffprobe", "-show_streams", "-print_format", "json", videoFile)
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

func GetAudioMetadata(videoFile string) (FfprobeShowStreamOutput, error) {
	output, err := execute("ffprobe", "-show_streams", "-print_format", "json", videoFile)
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

func TakeScreenshot(videoFile string, second float64, targetFile string) error {
	_, err := execute("ffmpeg", "-y", "-ss", fmt.Sprintf("%f", second), "-i", videoFile, "-vframes", "1", targetFile)
	if err != nil {
		return err
	}

	return nil
}

func CropScreenshot(videoFile string, second float64, rect model.RectFloat, targetFile string) error {
	cropFilter := fmt.Sprintf("crop=%f:%f:%f:%f", rect.W, rect.H, rect.X, rect.Y)
	_, err := execute("ffmpeg", "-y", "-ss", fmt.Sprintf("%f", second), "-i", videoFile, "-vf", cropFilter, "-vframes", "1", targetFile)
	if err != nil {
		return err
	}

	return nil
}
func ExtractPartOfVideo(videoFile string, second float64, duration int, targetFile string) error {
	_, err := execute("ffmpeg", "-ss", fmt.Sprintf("%f", second), "-i", videoFile,
		"-t", fmt.Sprintf("%d", duration), "-vcodec", "copy", "-acodec", "aac", "-ac", "4", targetFile)
	if err != nil {
		return err
	}

	return nil
}

func JoinVideoFiles(videoFiles []string, targetFile string) error {
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

	_, err = execute("ffmpeg", "-y", "-safe", "0", "-f", "concat", "-i", tempFile.Name(), targetFile)
	if err != nil {
		return err
	}

	return nil
}

func OptimizeVideoForPreview(videoFile string, tempFile string) error {
	_, err := execute("ffmpeg", "-y", "-i", videoFile, "-b:v", "2M", "-an", tempFile)
	if err != nil {
		return err
	}

	return os.Rename(tempFile, videoFile)
}

func ChangeVideoResolution(videoFile string, tempFile string, newResolution string) error {
	logger.Infof("Changing video resolution for %s to %s", videoFile, newResolution)

	_, err := execute("ffmpeg", "-i", videoFile, "-vf", fmt.Sprintf("scale=%s", newResolution), "-c:a", "copy", tempFile)
	if err != nil {
		return err
	}

	if err := os.Remove(videoFile); err != nil {
		return err
	}

	return os.Rename(tempFile, videoFile)
}
