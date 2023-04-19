package config

type dataWidth int

type Config interface {
	getDataWidthInBytes() dataWidth
	getColorSimilarityThreshold() int
	getInputFilePath() string
	getOutputFilePath() string
	getFlushAfterEveryWrite() bool
	getFlushAfterEveryCellWrite() bool
	getImageWidth() int
	getImageHeight() int
	getWidthInCells() int
	getHeightInCells() int
	getCellWidth() int
	getCellHeight() int
}

// Data width in bits to bytes
const (
	DataWidth8 dataWidth = 1 << iota
	DataWidth16
	DataWidth32
	DataWidth64
)

func New(
	dataWidthInBytes dataWidth,
	colorSimilarityThreshold int,
	inputFilePath string,
	outputFilePath string,
	flushAfterEveryWrite bool,
	flushAfterEveryCellWrite bool,
) Config {
	imageWidth := int(48 * dataWidthInBytes)
	imageHeight := int(64 * dataWidthInBytes)
	cellWidth := 6
	cellHeight := 8
	return &config{
		dataWidthInBytes:         dataWidthInBytes,
		colorSimilarityThreshold: colorSimilarityThreshold,
		inputFilePath:            inputFilePath,
		outputFilePath:           outputFilePath,
		flushAfterEveryWrite:     flushAfterEveryWrite,
		flushAfterEveryCellWrite: flushAfterEveryCellWrite,
		imageWidth:               imageWidth,
		imageHeight:              imageHeight,
		widthInCells:             imageWidth / cellWidth,
		heightInCells:            imageHeight / cellHeight,
		cellWidth:                cellWidth,
		cellHeight:               cellHeight,
	}
}

type config struct {
	dataWidthInBytes         dataWidth
	colorSimilarityThreshold int
	inputFilePath            string
	outputFilePath           string
	flushAfterEveryWrite     bool
	flushAfterEveryCellWrite bool
	imageWidth               int
	imageHeight              int
	widthInCells             int
	heightInCells            int
	cellWidth                int
	cellHeight               int
}

func (c *config) getDataWidthInBytes() dataWidth {
	return c.dataWidthInBytes
}

func (c *config) getColorSimilarityThreshold() int {
	return c.colorSimilarityThreshold
}

func (c *config) getInputFilePath() string {
	return c.inputFilePath
}

func (c *config) getOutputFilePath() string {
	return c.outputFilePath
}

func (c *config) getFlushAfterEveryWrite() bool {
	return c.flushAfterEveryWrite
}

func (c *config) getFlushAfterEveryCellWrite() bool {
	return c.flushAfterEveryCellWrite
}

func (c *config) getImageWidth() int {
	return c.imageWidth
}

func (c *config) getImageHeight() int {
	return c.imageHeight
}

func (c *config) getWidthInCells() int {
	return c.widthInCells
}
func (c *config) getHeightInCells() int {
	return c.heightInCells
}

func (c *config) getCellWidth() int {
	return c.cellWidth
}

func (c *config) getCellHeight() int {
	return c.cellHeight
}
