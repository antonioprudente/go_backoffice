package note

import "time"

// UserFilter rappresenta i filtri opzionali applicabili alle liste utenti.
type NoteFilter struct {
	Search     string     `form:"search"`
	CreateFrom *time.Time `form:"create_from"`
	CreateTo   *time.Time `form:"create_to"`
}
