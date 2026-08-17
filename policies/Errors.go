package policies

import "errors"

// Errori di autorizzazione: denial legittimo di business, non un guasto tecnico.
// Il controller li mappa tipicamente su HTTP 403.
var (
	// ErrForbidden è il denial generico ("non hai i permessi per questa operazione").
	// Già usato da AgentPolicy/AgencyPolicy per i casi di ruolo non ammesso.
	ErrForbidden = errors.New("Policy: Non hai i permessi per eseguire questa azione")
)

// Errori tecnici/di configurazione: indicano uno stato che la policy non
// dovrebbe mai incontrare in condizioni normali (ruolo sconosciuto, dato
// mancante che il chiamante avrebbe dovuto precaricare). Il controller li
// mappa tipicamente su HTTP 500, non 403: non è colpa dell'utente.
var (
	// ErrUnknownRole indica un valore di Role non previsto dallo switch
	// (né ADMIN, OPERATOR, AGENT, AGENCY, USER). Segnala una regressione
	// o un enum esteso senza aggiornare la policy corrispondente.
	ErrUnknownRole = errors.New("Policy: ruolo non gestito")

	// ErrMissingRelation indica che una relazione necessaria alla decisione
	// (es. target.Foreign per un target di Role USER) non è stata precaricata
	// dal repository. Non è un problema di permessi, ma di dati mancanti.
	ErrMissingRelation = errors.New("Policy: relazione necessaria non caricata")

	// ErrNotImplemented segnala uno stub non ancora implementato.
	ErrNotImplemented = errors.New("Policy: operazione non ancora implementata")
)
