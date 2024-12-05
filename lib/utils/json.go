package utils

import (
	"encoding/json"
	"flag"
	"os"
	"path"
)

var (
	cacheFolder = flag.String("cache_folder", ".cache", "cache folder")
)

func CJsonLoad(filename string, data any) error {
	filename = path.Join(*cacheFolder, filename)
	dataBuf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(dataBuf, data); err != nil {
		return err
	}
	return nil
}

func CJsonSave(filename string, data any) error {
	filename = path.Join(*cacheFolder, filename)
	if err := os.MkdirAll(path.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	dataBuf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, dataBuf, 0644)
}
