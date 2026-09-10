package user

// UserFilter rappresenta i filtri opzionali applicabili alle liste utenti.
type UserFilter struct {
	Status string `form:"status"`
	Search string `form:"search"` // cerca in first_name, last_name, username, email
}
