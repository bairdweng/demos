package resolver

import (
	"context"
	"iQuest/app/Api"
	"iQuest/app/graphql/model"
	m_user "iQuest/app/model/user"
)

func (r *mutationResolver) UserVerified(ctx context.Context, verifiedData model.VerifiedInput) (bool, error) {

	user := ctx.Value("user").(m_user.SessionUser)

	_, err := Api.UserVerifiedAndCreate(m_user.UserServiceVerfiedInput{
		CompanyId:  int(user.CompanyID),
		RealName:   verifiedData.Name,
		IdCardNo:   verifiedData.IDCardNo,
		BankCardNo: verifiedData.BankCardNo,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
