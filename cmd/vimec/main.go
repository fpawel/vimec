package main

import (
	"os"
	"path/filepath"
	"time"
)

func main() {
	_ = os.RemoveAll(filepath.Join(filepath.Dir(os.Args[0]), "pdf"))
	modelTableActs.SelectByDate(time.Now().Year(), int(time.Now().Month()))
	runMainWindow()
}
