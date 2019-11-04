package resolver

import (
	"iQuest/app/graphql/prisma"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	Prisma *prisma.Client
}

// func (r *Resolver) Mutation() schema.MutationResolver {
// 	return &mutationResolver{r}
// }
// func (r *Resolver) Query() schema.QueryResolver {
// 	return &queryResolver{r}
// }
// func (r *Resolver) Subscription() schema.SubscriptionResolver {
// 	return &subscriptionResolver{r}
// }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

type subscriptionResolver struct{ *Resolver }
