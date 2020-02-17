package fb

// Profile holds information about Facebook user profile
type Profile struct {
	ID   string // id of user
	Name string // full name of user (name nick surname)
}

// Link gets link to this user profile
func (p *Profile) Link() string {
	return "https://www.facebook.com/" + p.ID
}
