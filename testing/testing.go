package testing

import (
	"os"
	"path"
	"runtime"

	"github.com/transientvariable/config"
	"github.com/transientvariable/log"
)

func init() {
	_, f, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(f), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	if err := config.Load(); err != nil {
		panic(err)
	}

	if err := log.SetDefault(log.New(log.WithLevel("debug"))); err != nil {
		panic(err)
	}
}
