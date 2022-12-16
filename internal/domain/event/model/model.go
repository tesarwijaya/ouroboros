package model

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/gofrs/uuid"
)

type (
	Event struct {
		ID          uuid.UUID
		Type        string
		ContentType esdb.ContentType
		Data        []byte
		Metadata    []byte
	}
)
