package policies

import "example/go_backoffice/models"

// belongsTo verifica, in modo nil-safe, che il ForeignID di u punti a ownerID
// (es. una AGENCY o uno USER agganciato a un dato AGENT/AGENCY). ForeignID nil
// qui è uno stato di business legittimo (es. nodo radice), non un errore.
func belongsTo(u *models.User, ownerID uint) bool {
	return u.ForeignID != nil && *u.ForeignID == ownerID
}
