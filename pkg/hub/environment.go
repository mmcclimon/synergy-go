package hub

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // sql driver
	"github.com/mmcclimon/synergy-go/internal/config"
)

// Environment represents global synergy things
type Environment struct {
	db            *sqlx.DB
	UserDirectory *Directory
}

// NewEnvironment gives you a new environment (handle onto the db, basically
func NewEnvironment(config config.Config) *Environment {
	db, err := sqlx.Open("sqlite3", config.StateDBFile)
	if err != nil {
		log.Fatalf("could not open sqlite db: %s", err)
	}

	env := Environment{db: db}

	// circular reference here is useful, and they're singletons so won't leak
	env.UserDirectory = NewDirectory(&env)

	env.maybeCreateStateTables()
	env.UserDirectory.loadUsers()

	return &env
}

func (env *Environment) maybeCreateStateTables() {
	_, err := env.db.Exec(`
		CREATE TABLE IF NOT EXISTS synergy_state (
			reactor_name TEXT PRIMARY KEY,
			stored_at INTEGER NOT NULL,
			json TEXT NOT NULL
		);
	`)

	if err != nil {
		log.Fatalf("could not create table: %s", err)
	}

	_, err = env.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			lp_id TEXT,
			is_master INTEGER DEFAULT 0,
			is_virtual INTEGER DEFAULT 0,
			is_deleted INTEGER DEFAULT 0
		);
	`)

	if err != nil {
		log.Fatalf("could not create table: %s", err)
	}

	_, err = env.db.Exec(`
		CREATE TABLE IF NOT EXISTS user_identities (
			id INTEGER PRIMARY KEY,
			username TEXT NOT NULL,
			identity_name TEXT NOT NULL,
			identity_value TEXT NOT NULL,
			FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE,
			CONSTRAINT constraint_username_identity UNIQUE (username, identity_name),
			UNIQUE (identity_name, identity_value)
		);
	`)

	if err != nil {
		log.Fatalf("could not create table: %s", err)
	}

}
