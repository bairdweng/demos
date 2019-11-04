package directive

import (
	"context"
	// m_user "iQuest/app/model/user"

	"github.com/99designs/gqlgen/graphql"
)

// IsAuthenticated 判断登录状态
func IsAuthenticated(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	// user := ctx.Value("user")
	// if user != nil && user.(m_user.SessionUser).UserID != 0 {
	// 	return next(ctx)
	// }
	return next(ctx)

	// return nil, &gqlerror.Error{
	// 	Message: "Unauthorised",
	// 	Extensions: map[string]interface{}{
	// 		"code": 401,
	// 	},
	// }
}
