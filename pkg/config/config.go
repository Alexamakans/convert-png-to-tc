package config

const CellWidth = 6
const CellHeight = 8

type dataWidth int

type Config interface {
	GetDataWidthInBytesAsInt() int
	GetDataWidthInBytes() dataWidth
	GetColorSimilarityThreshold() int
	GetFlushAfterEveryWrite() bool
	GetFlushAfterEveryCellWrite() bool
	GetImageWidth() int
	GetImageHeight() int
	GetWidthInCells() int
	GetHeightInCells() int

	SetDataWidthInBytes(dataWidth)
	SetColorSimilarityThreshold(int)
	SetFlushAfterEveryWrite(bool)
	SetFlushAfterEveryCellWrite(bool)
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
	flushAfterEveryWrite bool,
	flushAfterEveryCellWrite bool,
) Config {
	imageWidth := int(48 * dataWidthInBytes)
	imageHeight := int(64 * dataWidthInBytes)
	return &config{
		dataWidthInBytes:         dataWidthInBytes,
		colorSimilarityThreshold: colorSimilarityThreshold,
		flushAfterEveryWrite:     flushAfterEveryWrite,
		flushAfterEveryCellWrite: flushAfterEveryCellWrite,
		imageWidth:               imageWidth,
		imageHeight:              imageHeight,
		widthInCells:             imageWidth / CellWidth,
		heightInCells:            imageHeight / CellHeight,
	}
}

type config struct {
	dataWidthInBytes         dataWidth
	colorSimilarityThreshold int
	flushAfterEveryWrite     bool
	flushAfterEveryCellWrite bool
	imageWidth               int
	imageHeight              int
	widthInCells             int
	heightInCells            int
}

func (c *config) GetDataWidthInBytesAsInt() int {
	return int(c.dataWidthInBytes)
}

func (c *config) GetDataWidthInBytes() dataWidth {
	return c.dataWidthInBytes
}

func (c *config) GetColorSimilarityThreshold() int {
	return c.colorSimilarityThreshold
}

func (c *config) GetFlushAfterEveryWrite() bool {
	return c.flushAfterEveryWrite
}

func (c *config) GetFlushAfterEveryCellWrite() bool {
	return c.flushAfterEveryCellWrite
}

func (c *config) GetImageWidth() int {
	return c.imageWidth
}

func (c *config) GetImageHeight() int {
	return c.imageHeight
}

func (c *config) GetWidthInCells() int {
	return c.widthInCells
}

func (c *config) GetHeightInCells() int {
	return c.heightInCells
}

func (c *config) SetDataWidthInBytes(dataWidthInBytes dataWidth) {
	c.dataWidthInBytes = dataWidthInBytes
	c.imageWidth = int(48 * dataWidthInBytes)
	c.imageHeight = int(64 * dataWidthInBytes)
	c.widthInCells = c.imageWidth / CellWidth
	c.heightInCells = c.imageHeight / CellHeight
}

func (c *config) SetColorSimilarityThreshold(colorSimilarityThreshold int) {
	c.colorSimilarityThreshold = colorSimilarityThreshold
}

func (c *config) SetFlushAfterEveryWrite(flushAfterEveryWrite bool) {
	c.flushAfterEveryWrite = flushAfterEveryWrite
}

func (c *config) SetFlushAfterEveryCellWrite(flushAfterEveryCellWrite bool) {
	c.flushAfterEveryCellWrite = flushAfterEveryCellWrite
}
