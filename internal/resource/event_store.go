package resource

import (
	"fmt"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/tesarwijaya/ouroboros/internal/config"
)

func NewEventStoreConnection(cfg *config.Config) (*esdb.Client, error) {
	conf, err := esdb.ParseConnectionString(fmt.Sprintf("esdb://%s:%d?tls=false", cfg.EventStoreDBHost, cfg.EventStoreDBPort))
	if err != nil {
		return nil, err
	}

	return esdb.NewClient(conf)
}
