package ffmpeg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/go-errors/errors"
)

type ffprobeShowFormatOutput struct {
	Format ffprobeFormatOutput `json:"format"`
}

type ffprobeFormatOutput struct {
	Duration string `json:"duration"`
}

func execute(name string, arg ...string) ([]byte, error) {
	fmt.Printf("Running %v %v\n", name, arg)

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

func GetDurationInSeconds(videoFile string) (uint64, error) {
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

	return uint64(durationInSeconds), nil
}

func TakeScreenshot(videoFile string, second uint64, targetFile string) error {
	_, err := execute("ffmpeg", "-y", "-ss", fmt.Sprintf("%d", second), "-i", videoFile, "-vframes", "1", targetFile)

	if err != nil {
		return err
	}

	return nil
}
