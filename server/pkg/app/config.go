package app

import (
	"context"
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

func (c *MyCollectionConfig) DebugPrint(ctx context.Context) {
	logger.Debugf(ctx, "Configuration:")
	logger.Debugf(ctx, strings.Repeat("-", 50))
	logger.Debugf(ctx, "  %-30s %s", "RootDir:", c.RootDir)
	logger.Debugf(ctx, "  %-30s %s", "ListenAddress:", c.ListenAddress)
	logger.Debugf(ctx, "  %-30s %d", "AutoMixItemsCount:", c.AutoMixItemsCount)
	logger.Debugf(ctx, "  %-30s %d", "MixOnDemandItemsCount:", c.MixOnDemandItemsCount)
	logger.Debugf(ctx, "  %-30s %d", "ItemsOptimizerMaxResolution:", c.ItemsOptimizerMaxResolution)
	logger.Debugf(ctx, "  %-30s %t", "ProcessorPaused:", c.ProcessorPaused)
	logger.Debugf(ctx, "  %-30s %d", "CoversCount:", c.CoversCount)
	logger.Debugf(ctx, "  %-30s %d", "PreviewSceneCount:", c.PreviewSceneCount)
	logger.Debugf(ctx, "  %-30s %d", "PreviewSceneDuration:", c.PreviewSceneDuration)
	logger.Debugf(ctx, "  %-30s %v", "OpenSubtitleApiKeys:", c.OpenSubtitleApiKeys)
	logger.Debugf(ctx, strings.Repeat("-", 50))
}
