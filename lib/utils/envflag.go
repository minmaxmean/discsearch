package utils

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var (
	logLevel = Level("loglevel", slog.LevelInfo, "slog level")
)

func LoadFlagsFromEnv() error {
	flag.Parse()
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Printf("could not load dotenv: %+v", os.IsNotExist(err))
		return err
	}
	alreadySet := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		alreadySet[f.Name] = true
	})
	var err error
	flag.VisitAll(func(f *flag.Flag) {
		if alreadySet[f.Name] || err != nil {
			return
		}
		if envar := os.Getenv(f.Name); envar != "" {
			if err = f.Value.Set(envar); err != nil {
				return
			}
		}
	})
	slog.SetLogLoggerLevel(*logLevel)
	return err
}
