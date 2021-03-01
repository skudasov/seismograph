package main

import (
	"log"
	"os"

	"github.com/f4hrenh9it/seismograph/back/config"
	"github.com/f4hrenh9it/seismograph/back/migrations"
	"github.com/f4hrenh9it/seismograph/back/server"
)

func main() {
	cfgPath := os.Getenv("CFG")
	if cfgPath == "" {
		log.Fatal("provide env flag CFG for seismograph config")
	}
	cfg := config.ViperCfg(cfgPath)
	// TODO: make migrate cli with first release
	migrations.GormMigrateInit(cfg.DB)
	server.NewSeismographService(&cfg)
}
