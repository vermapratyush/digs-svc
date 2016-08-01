package mapper

import (
	"digs/domain"
	"digs/socket"
	"digs/common"
	"digs/models"
)

func MapUserAccountToPersonResponse(userAccounts []models.UserAccount) []domain.PersonResponse {
	res := make([]domain.PersonResponse, len(userAccounts))

	for idx, userAccount := range(userAccounts) {
		_, present := socket.GetLookUp(userAccount.UID)
		activeState := "active"
		if !present {
			activeState = "inactive"
		}
		res[idx] = domain.PersonResponse{
			Name: common.GetName(userAccount.FirstName, userAccount.LastName),
			UID: userAccount.UID,
			About: userAccount.About,
			Verified: userAccount.Verified,
			ActiveState: activeState,
			ProfilePicture: userAccount.ProfilePicture,
		}
	}
	return res
}
