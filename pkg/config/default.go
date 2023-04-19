package config

const colorSimilarityThreshold = 6
const flushAfterEveryWrite = false
const flushAfterEveryCellWrite = true

var DefaultConfig = New(
	DataWidth16,
	colorSimilarityThreshold,
	flushAfterEveryWrite,
	flushAfterEveryCellWrite,
)
