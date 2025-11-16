package opensubtitles

import (
	"encoding/json"
	"fmt"
	"my-collection/server/pkg/model"
)

type listResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Release string `json:"release"`
		} `json:"attributes"`
	} `json:"data"`
}

func (s *OpenSubtitiles) List(imdbId string) ([]model.SubtitleMetadata, error) {
	url, err := s.buildUrl("subtitles")
	if err != nil {
		return nil, err
	}

	query := url.Query()
	query.Set("imdb_id", imdbId)
	query.Set("languages", s.language)
	if !s.aiTranslated {
		query.Set("ai_translated", "exclude")
	}
	url.RawQuery = query.Encode()

	req, err := s.buildRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := s.fetch(req)
	if err != nil {
		return nil, err
	}
	return s.parseListResponse(body)
}

func (s *OpenSubtitiles) parseListResponse(body []byte) ([]model.SubtitleMetadata, error) {
	var apiResp listResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	results := make([]model.SubtitleMetadata, 0, len(apiResp.Data))
	for _, item := range apiResp.Data {
		results = append(results, model.SubtitleMetadata{
			Id:      item.ID,
			Release: item.Attributes.Release,
		})
	}

	return results, nil
}
