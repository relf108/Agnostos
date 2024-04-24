package config

import (
	"context"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/apple/pkl-go/pkl"
)

func ReadConfig(path string) Config {
	evaluator, err := pkl.NewEvaluator(context.Background(), pkl.PreconfiguredOptions)
	if err != nil {
		panic(err)
	}
	defer evaluator.Close()
	var cfg Config
	if err = evaluator.EvaluateModule(context.Background(), pkl.FileSource(path), &cfg); err != nil {
		panic(err)
	}
	return cfg
}

func FindConfig() string {
	var confPath string
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".pkl") {
			return nil
		}
		if strings.HasSuffix(d.Name(), ".pkl") {
			confPath = path
			return io.EOF
		}
		return nil
	})
  if err != nil && err != io.EOF {
    panic(err)
  }
	if confPath == "" {
		panic("No config file found")
	}
	return confPath
}
