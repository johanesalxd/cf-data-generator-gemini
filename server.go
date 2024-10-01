package datageneratorgemini

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// init initializes the HTTP function handler for BQRFGemini
func init() {
	initOnce.Do(initAll)
	functions.HTTP("DataGeneratorGemini", DataGeneratorGemini)
}
