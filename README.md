# moconv

[WIP] moconv is generate DB model from domain model.

```
go get -u github.com/inari111/moconv
```

```
moconv gen
```

## example

domain model  
examples/domain/user/user.go
```go
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
```

```
moconv gen
```

generate examples/infra/postgres/user.go
```go
package postgres

import (
	domain "github.com/inari111/moconv/examples/domain"
	user "github.com/inari111/moconv/examples/domain/user"
	"time"
	u "u"
)

type User struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *User) ToDomain() *user.User {
	return &user.User{
		ID:        domain.UserID(u.ID),
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
func NewUser(u *user.User) *User {
	return &User{
		ID:        u.ID.String(),
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
```
