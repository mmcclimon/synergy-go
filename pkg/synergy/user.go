package synergy

// User represents a user.
type User struct {
	Username   string
	LPID       string
	IsMaster   bool
	IsVirtual  bool
	IsDeleted  bool
	Identities map[string]string
}

// AddIdentity sets the identity for channel "name" to "val"
func (u *User) AddIdentity(name, val string) {
	u.Identities[name] = val
}
