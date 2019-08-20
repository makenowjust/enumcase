package enumcase_test

import (
	"testing"

	"github.com/MakeNowJust/enumcase"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, enumcase.Analyzer, "a")
}
