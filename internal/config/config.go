package config

import "github.com/jinzhu/configor"

// SyncGroup between Nextcloud and Trello.
type SyncGroup struct {
	Name      string
	Type      string
	Nextcloud struct {
		Board int
		Stack int
	}
	Trello struct {
		Board string
		List  string
	}
}

// Config data.
type Config struct {
	Interval  uint64
	Debug     bool
	Log       string
	Nextcloud struct {
		API      string
		Username string
		Password string
	}
	Trello struct {
		Key   string
		Token string
	}
	Sync []SyncGroup
}

// Load config.
func Load() Config {
	var conf Config
	err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&conf, "config.yml")
	if err != nil {
		panic(err)
	}
	return conf
}
