package user

import (
	"time"

	"github.com/inari111/moconv/examples/domain"
)

type User struct {
	ID        domain.UserID
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
