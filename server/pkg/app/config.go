package app

import "my-collection/server/pkg/directorytree"

type MyCollectionConfig struct {
	RootDir                     string
	ListenAddress               string
	FilesFilter                 directorytree.FilesFilter
	AutoMixItemsCount           int
	MixOnDemandItemsCount       int
	ItemsOptimizerMaxResolution int
	ProcessorPaused             bool
	CoversCount                 int
	PreviewSceneCount           int
	PreviewSceneDuration        int
	OpenSubtitleApiKeys         []string
}
