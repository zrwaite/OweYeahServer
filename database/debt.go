package database

import (
	"errors"
	"fmt"

	"github.com/zrwaite/OweMate/graph/model"
	"github.com/zrwaite/OweMate/utils"
)

func ResolveCycles(user *model.User, connection *model.UserConnection) (err error) {
	for {
		cycleFound, cycleConenctions, maxCycleDebt := DetectNestedDebts(5, user, true, user, connection, connection.Debt, []*model.User{}, []*model.UserConnection{})
		if !cycleFound {
			break
		}
		for _, connection := range cycleConenctions {
			connection.Debt -= maxCycleDebt
			success := UpdateUserConnection(connection)
			if !success {
				return errors.New("failed to update connection")
			}
		}
	}
	return nil
}

func DetectNestedDebts(
	depth int,
	user *model.User,
	root bool,
	rootUser *model.User,
	parentConnection *model.UserConnection,
	maxDebt float64,
	cycleUsers []*model.User,
	cycleConnections []*model.UserConnection,
) (cycleFound bool, cycleConnectionsFound []*model.UserConnection, maxCycleDebt float64) {
	if depth <= 0 {
		return false, nil, 0
	}
	userConnections, err := GetUserConnections(user.ConnectionIds, user.Username)
	if err != nil {
		fmt.Println("Failed to get user connections:", err)
		return false, nil, 0
	}
	for _, connection := range userConnections {
		newCycleConnections := append(cycleConnections, connection)
		positiveMinDebt := maxDebt > 0
		if connection.ID == parentConnection.ID || // don't go back up tree
			connection.Debt == 0 || // no debt to settle
			(!root && //While not at the root node:
				(connection.Debt < 0 && positiveMinDebt || // Negative debt and positive parent debt
					connection.Debt > 0 && !positiveMinDebt)) { // Positive debt and negative parent debt
			continue
		}

		userInCycle, usernameIndex := utils.UserBinarySearch(cycleUsers, connection.Contact)
		if userInCycle {
			continue
		}
		connectionContact, err := GetUserConnectionContact(connection.ContactUsername)
		if err != nil {
			fmt.Println("Failed to get user connection contact:", err)
			return false, nil, 0
		}
		newCycleUsers := utils.ArrayInsert(cycleUsers, usernameIndex, connectionContact)

		if connection.Debt < maxDebt && positiveMinDebt || connection.Debt > maxDebt && !positiveMinDebt {
			// Set new minimum debt for the traversal
			maxDebt = connection.Debt
		}
		if connection.ContactUsername == rootUser.Username {
			return true, newCycleConnections, maxDebt
		}
		nestedCycleFound, nestedCycleConnections, newMaxDebt := DetectNestedDebts(depth-1, connectionContact, false, rootUser, connection, maxDebt, newCycleUsers, newCycleConnections)
		if nestedCycleFound {
			return true, nestedCycleConnections, newMaxDebt
		}
	}
	return false, nil, 0
}
