package opensubtitles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"my-collection/server/pkg/model"
	"net/http"
	"os"
)

type downloadReq struct {
	FileId string `json:"file_id"`
}

type downloadResp struct {
	Link string `json:"link"`
}

func (s *OpenSubtitiles) Download(subtitle model.SubtitleMetadata, outputFile string) error {
	downloadUrl, err := s.buildUrl("download")
	if err != nil {
		return fmt.Errorf("failed to build download URL: %w", err)
	}

	reqJson, err := s.buildDownloadRequest(subtitle.Id)
	if err != nil {
		return err
	}

	req, err := s.buildRequest("POST", downloadUrl, bytes.NewReader(reqJson))
	if err != nil {
		return fmt.Errorf("failed to build download request: %w", err)
	}

	body, err := s.fetch(req)
	if err != nil {
		return fmt.Errorf("failed to fetch download link: %w", err)
	}

	return s.parseDownloadResponse(body, outputFile)

}

func (s *OpenSubtitiles) buildDownloadRequest(id string) ([]byte, error) {
	req := downloadReq{FileId: id}
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	return reqJson, nil
}

func (s *OpenSubtitiles) parseDownloadResponse(body []byte, outputFile string) error {
	var downloadResp downloadResp
	if err := json.Unmarshal(body, &downloadResp); err != nil {
		return fmt.Errorf("failed to parse download response: %w", err)
	}

	if downloadResp.Link == "" {
		return fmt.Errorf("no download link in response")
	}

	return s.downloadLink(downloadResp.Link, outputFile)
}

func (s *OpenSubtitiles) downloadLink(link string, outputFile string) error {
	fileReq, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return fmt.Errorf("failed to create file download request: %w", err)
	}

	body, err := s.fetch(fileReq)
	if err != nil {
		return fmt.Errorf("failed to download subtitle file: %w", err)
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to write subtitle file: %w", err)
	}

	return nil
}
