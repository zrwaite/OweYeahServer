package database

import (
	"github.com/zrwaite/OweMate/graph/model"
)

func ResolveCycles(user *model.User, connection *model.UserConnection) {
	for {
		// cycleFound, cycleContactUsernamesFound := DetectNestedDebts(5, user, true, user.Username, connection.ID, connection.Debt, []string{})
		// if !cycleFound {
		// 	break
		// }
		// for _, username := range cycleContactUsernamesFound {

		// }
	}
}

func DetectNestedDebts(
	depth int,
	user *model.User,
	root bool,
	rootUsername string,
	parentConnectionId string,
	maxDebt float64,
	cycleUsers []*model.User) (cycleFound bool, cycleContactUsernamesFound []string) {
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
		// userInCycle, usernameIndex := utils.UserBinarySearch(cycleContactUsernames, connection.ContactUsername)
		// if userInCycle {
		// 	continue
		// }
		// newCycleContactUsernames := utils.ArrayInsert(cycleContactUsernames, usernameIndex, connection.ContactUsername)

		if connection.Debt < maxDebt && positiveMinDebt || connection.Debt > maxDebt && !positiveMinDebt {
			// Set new minimum debt for the traversal
			maxDebt = connection.Debt
		}
		if connection.ContactUsername == rootUsername {

		}
		// DetectNestedDebts(depth-1, connection.Contact, false, rootUsername, connection.ID, maxDebt, newCycleContactUsernames)
	}
	return
}
