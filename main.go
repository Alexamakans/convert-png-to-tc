package main

import (
	"github.com/Alexamakans/convert-png-to-tc/pkg/config"
	"github.com/Alexamakans/convert-png-to-tc/pkg/converter"
)

const inputFilePath = "./images/test_img_96x128.png"
const outputFilePath = "./outputs/test_img_96x128.tc"

var conv = converter.New()

func main() {
	conv.Config.SetColorSimilarityThreshold(8)
	conv.Config.SetDataWidthInBytes(config.DataWidth16)
	conv.Config.SetFlushAfterEveryWrite(false)
	conv.Config.SetFlushAfterEveryCellWrite(true)

	if err := conv.ConvertFile(inputFilePath, outputFilePath); err != nil {
		panic(err)
	}
}
