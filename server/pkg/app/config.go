package app

import (
	"my-collection/server/pkg/directorytree"
	"strings"
)

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

func (c *MyCollectionConfig) DebugPrint() {
	logger.Debugf("Configuration:")
	logger.Debugf(strings.Repeat("-", 50))
	logger.Debugf("  %-30s %s", "RootDir:", c.RootDir)
	logger.Debugf("  %-30s %s", "ListenAddress:", c.ListenAddress)
	logger.Debugf("  %-30s %d", "AutoMixItemsCount:", c.AutoMixItemsCount)
	logger.Debugf("  %-30s %d", "MixOnDemandItemsCount:", c.MixOnDemandItemsCount)
	logger.Debugf("  %-30s %d", "ItemsOptimizerMaxResolution:", c.ItemsOptimizerMaxResolution)
	logger.Debugf("  %-30s %t", "ProcessorPaused:", c.ProcessorPaused)
	logger.Debugf("  %-30s %d", "CoversCount:", c.CoversCount)
	logger.Debugf("  %-30s %d", "PreviewSceneCount:", c.PreviewSceneCount)
	logger.Debugf("  %-30s %d", "PreviewSceneDuration:", c.PreviewSceneDuration)
	logger.Debugf("  %-30s %v", "OpenSubtitleApiKeys:", c.OpenSubtitleApiKeys)
	logger.Debugf(strings.Repeat("-", 50))
}
