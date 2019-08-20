package main

import (
	"github.com/MakeNowJust/enumcase"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(enumcase.Analyzer) }
