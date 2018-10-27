package repository

import (
	"github.com/larwef/ki/internal/adding"
	"github.com/larwef/ki/internal/listing"
)

// Repository has to satisfy adding and listing repository interfaces.
type Repository interface {
	adding.Repository
	listing.Repository
}
