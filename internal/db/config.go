package db

import (
	"encoding/json"
	"os"
)

const (
	ConfigFile    = "config.json"
	DefaultDbPath = "./.db.json"
)

func getDb() *Db {
	c := getDbConfig()
	f, err := os.Open(c.DbPath)
	if err != nil {
		return &Db{}
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	db := Db{}
	err = dec.Decode(&db)
	if err != nil {
		return &Db{}
	}
	return &db
}

func getDbConfig() Config {
	c := Config{}
	cwd, err := os.Getwd()
	if err != nil {
		c.DbPath = DefaultDbPath
		return c
	}
	f, err := os.Open(cwd + "/" + ConfigFile)
	if err != nil {
		c.DbPath = DefaultDbPath
		return c
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(&c)
	if err != nil {
		c.DbPath = DefaultDbPath
		return c
	} else if c.DbPath == "" {
		c.DbPath = DefaultDbPath
	}
	return c
}
