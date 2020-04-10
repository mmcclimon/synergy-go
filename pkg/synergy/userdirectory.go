package synergy

import (
	"database/sql"
	"log"
)

// Directory represents a directory of users
type Directory struct {
	env   *Environment
	users map[string]User
}

func (ud *Directory) loadUsers() {
	if ud.users == nil {
		ud.users = make(map[string]User)
	}

	rows, err := ud.env.db.Queryx("select * from users")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	type rawUser struct {
		Username  string
		LPID      sql.NullString `db:"lp_id"`
		IsMaster  sql.NullBool   `db:"is_master"`
		IsVirtual sql.NullBool   `db:"is_virtual"`
		IsDeleted sql.NullBool   `db:"is_deleted"`
	}

	for rows.Next() {
		var raw rawUser

		err = rows.StructScan(&raw)
		if err != nil {
			log.Fatal(err)
		}

		username := raw.Username

		ud.users[username] = User{
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
func (ud *Directory) UserNamed(name string) (User, bool) {
	user, ok := ud.users[name]
	return user, ok
}

// UserByChannelAndAddress gives you a user, given a channel name and address
func (ud *Directory) UserByChannelAndAddress(channelName, addr string) *User {
	for _, user := range ud.users {
		ident := user.Identities[channelName]
		if ident == addr {
			return &user
		}
	}

	return nil
}
