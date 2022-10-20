package user

import (
	"context"
	"log"
	"sync"
	"time"

	psql "github.com/hyphengolang/prelude/sql/postgres"
	"github.com/hyphengolang/prelude/types/email"
	"github.com/hyphengolang/prelude/types/password"

	"github.com/jackc/pgx/v5"

	"github.com/rog-golang-buddies/rmx/internal"
	"github.com/rog-golang-buddies/rmx/internal/suid"
)

/*
User schema
users
id uuid primary key,
text,
text unique,
text,
timestamp default::now,
timestamp nullable,
timestamp nullable,
*/

type User struct {
	ID        suid.UUID
	Username  string
	Email     email.Email
	Password  password.PasswordHash
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type Repo struct {
	ctx context.Context

	// connection is using pgx until sqlc is sorted
	c *pgx.Conn
}

func (r *Repo) Context() context.Context {
	if r.ctx != nil {
		return r.ctx
	}
	return context.Background()
}

func (r *Repo) Close(ctx context.Context) error { return r.c.Close(ctx) }

func NewRepo(ctx context.Context, conn *pgx.Conn) internal.UserRepo {
	r := &Repo{ctx, conn}
	return r
}

const (
	qryInsert = `insert into "user" (id, email, username, password) values (@id, @email, @username, @password)`

	qrySelectMany = `select id, email, username, password from "user" order by id`

	qrySelectByID       = `select id, email, username, password from "user" where id = $1`
	qrySelectByEmail    = `select id, email, username, password from "user" where email = $1`
	qrySelectByUsername = `select id, email, username, password from "user" where username = $1`

	qryDeleteByID       = `delete from "user" where id = $1`
	qryDeleteByEmail    = `delete from "user" where email = $1`
	qryDeleteByUsername = `delete from "user" where username = $1`
)

func (r *Repo) Insert(ctx context.Context, u *internal.User) error {
	args := pgx.NamedArgs{
		"id":       u.ID,
		"email":    u.Email,
		"username": u.Username,
		"password": u.Password,
	}

	return psql.Exec(r.c, qryInsert, args)
}

func (r *Repo) SelectMany(ctx context.Context) ([]internal.User, error) {
	return psql.Query(r.c, qrySelectMany, func(r pgx.Rows, u *internal.User) error {
		return r.Scan(&u.ID, &u.Email, &u.Username, &u.Password)
	})
}

func (r *Repo) Select(ctx context.Context, key any) (*internal.User, error) {
	var qry string
	switch key.(type) {
	case suid.UUID:
		qry = qrySelectByID
	case email.Email:
		qry = qrySelectByEmail
	case string:
		qry = qrySelectByUsername
	default:
		return nil, internal.ErrInvalidType
	}
	var u internal.User
	return &u, psql.QueryRow(r.c, qry, func(r pgx.Row) error { return r.Scan(&u.ID, &u.Username, &u.Email, &u.Password) }, key)
}

func (r *Repo) Delete(ctx context.Context, key any) error {
	var qry string
	switch key.(type) {
	case suid.UUID:
		qry = qryDeleteByID
	case email.Email:
		qry = qryDeleteByEmail
	case string:
		qry = qryDeleteByUsername
	default:
		return internal.ErrInvalidType
	}
	return psql.Exec(r.c, qry, key)
}

var DefaultRepo = &repo{
	miu:  make(map[suid.UUID]*User),
	mei:  make(map[string]*User),
	log:  log.Println,
	logf: log.Printf,
}

func (r *repo) Close(ctx context.Context) error { return nil }

func (r *repo) Delete(ctx context.Context, key any) error { return nil }

type repo struct {
	mu  sync.Mutex
	miu map[suid.UUID]*User
	mei map[string]*User

	log  func(v ...any)
	logf func(format string, v ...any)
}

// Remove implements internal.UserRepo
func (r *repo) Remove(ctx context.Context, key any) error {
	panic("unimplemented")
}

func (r *repo) Insert(ctx context.Context, iu *internal.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.mei[iu.Email.String()]; found {
		return internal.ErrAlreadyExists
	}

	u := &User{
		ID:        iu.ID,
		Username:  iu.Username,
		Email:     iu.Email,
		Password:  iu.Password,
		CreatedAt: time.Now(),
	}
	r.mei[iu.Email.String()], r.miu[iu.ID] = u, u

	return nil
}

func (r *repo) SelectMany(ctx context.Context) ([]internal.User, error) {
	panic("not implemented")
}

func (r *repo) Select(ctx context.Context, key any) (*internal.User, error) {
	switch key := key.(type) {
	case suid.UUID:
		return r.selectUUID(key)
	case email.Email:
		return r.selectEmail(key)
	case string:
		return r.selectUsername(key)
	default:
		return nil, internal.ErrInvalidType
	}
}

func (r *repo) selectUUID(uid suid.UUID) (*internal.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if u, ok := r.miu[uid]; ok {
		return &internal.User{
			ID:       u.ID,
			Username: u.Username,
			Email:    email.Email(u.Email),
			Password: password.PasswordHash(u.Password),
		}, nil
	}

	return nil, internal.ErrNotFound
}

func (r *repo) selectUsername(username string) (*internal.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.mei {
		if u.Username == username {
			return &internal.User{
				ID:       u.ID,
				Username: u.Username,
				Email:    email.Email(u.Email),
				Password: password.PasswordHash(u.Password),
			}, nil
		}
	}

	return nil, internal.ErrNotFound
}

func (r *repo) selectEmail(email email.Email) (*internal.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if u, ok := r.mei[email.String()]; ok {
		return &internal.User{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Password: password.PasswordHash(u.Password),
		}, nil
	}

	return nil, internal.ErrNotFound
}
