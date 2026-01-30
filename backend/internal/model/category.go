package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	NameKey   string    `json:"nameKey,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CategoryInput struct {
	Name string `json:"name"`
}

func (c *CategoryInput) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	return nil
}
