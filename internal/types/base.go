package types

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Base struct {
	ID        ulid.ULID `gorm:"index:idx_brin,type:brin"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
}

func NewBase() (*Base, error) {
	id, err := ulid.New(ulid.Now(), nil)
	if err != nil {
		return nil, err
	}

	return &Base{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Time{},
	}, nil
}
