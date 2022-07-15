package database

import "github.com/zrwaite/OweMate/graph/model"

func DetectNestedDebts(
	depth int,
	user *model.User,
	root bool,
	parentConnectionId string,
	maxDebt float64,
	checkedContactUsernames []string) {
	if depth <= 0 {
		return
	}
	for _, connection := range user.Connections {
		positiveMinDebt := maxDebt > 0
		if connection.ID == parentConnectionId || // don't go back up tree
			connection.Debt == 0 || // no debt to settle
			(!root && //While not at the root node:
				(connection.Debt < 0 && positiveMinDebt || // Negative debt and positive parent debt
					connection.Debt > 0 && !positiveMinDebt)) { // Positive debt and negative parent debt
			continue
		}
		// if connection.ContactUsername {
		// }
		if connection.Debt < maxDebt && positiveMinDebt || connection.Debt > maxDebt && !positiveMinDebt {
			// Set new minimum debt for the traversal
			maxDebt = connection.Debt
		}

	}
}
