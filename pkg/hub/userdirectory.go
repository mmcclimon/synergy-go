package hub

import (
	"database/sql"
	"log"

	"github.com/mmcclimon/synergy-go/pkg/user"
)

// Directory represents a directory of users
type Directory struct {
	env   *Environment
	users map[string]user.User
}

type rawUser struct {
	Username  string
	LPID      sql.NullString `db:"lp_id"`
	IsMaster  sql.NullBool   `db:"is_master"`
	IsVirtual sql.NullBool   `db:"is_virtual"`
	IsDeleted sql.NullBool   `db:"is_deleted"`
}

// NewDirectory gives you a new user directory. Note: we have to inject the db
// here!
func NewDirectory(env *Environment) *Directory {
	return &Directory{
		env:   env,
		users: make(map[string]user.User),
	}
}

func (ud *Directory) loadUsers() {
	rows, err := ud.env.db.Queryx("select * from users")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var raw rawUser

		err = rows.StructScan(&raw)
		if err != nil {
			log.Fatal(err)
		}

		username := raw.Username

		// I could make a constructor for this, but also, meh
		ud.users[username] = user.User{
			Username:   raw.Username,
			LPID:       raw.LPID.String,
			IsMaster:   raw.IsMaster.Bool,
			IsVirtual:  raw.IsVirtual.Bool,
			IsDeleted:  raw.IsDeleted.Bool,
			Identities: make(map[string]string),
		}
	}

	ud.loadIdentities()
}

func (ud *Directory) loadIdentities() {
	rows, err := ud.env.db.Query(
		"select username, identity_name, identity_value from user_identities",
	)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var username, name, val string
		err = rows.Scan(&username, &name, &val)

		if err != nil {
			log.Println("error scanning identity row", err)
			continue
		}

		user, ok := ud.users[username]
		if ok {
			user.AddIdentity(name, val)
		} else {
			log.Printf("found identity for %s but no matching user", username)
		}
	}
}

// UserNamed gives you the user for a name (if we have one)
func (ud *Directory) UserNamed(name string) (user.User, bool) {
	user, ok := ud.users[name]
	return user, ok
}
