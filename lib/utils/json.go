package utils

import (
	"encoding/json"
	"flag"
	"os"
	"path"
)

var (
	cacheFolder  = flag.String("cache_folder", ".cache", "cache folder")
	cacheEnabled = flag.Bool("cache_enabled", true, "is cached_exec enabled")
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
	dataBuf, err := json.MarshalIndent(data, "  ", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, dataBuf, 0644)
}

func CachedExec[R any](key string, f func() (R, error)) (R, error) {
	if !*cacheEnabled {
		return f()
	}
	var val R
	if err := CJsonLoad(key+".json", &val); err == nil {
		return val, nil
	}
	val, err := f()
	if err != nil {
		return val, err
	}
	if err := CJsonSave(key+".json", val); err != nil {
		return val, err
	}
	return val, nil
}
