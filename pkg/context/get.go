package context

import (
	"context"
	"errors"

	"go-ecommerce-backend-api.com/pkg/cache"
)

type InfoUserUUID struct {
	UserId      uint64
	UserAccount string
}

func GetSubjectUUID(ctx context.Context) (string, error) {
	sUUID, ok := ctx.Value("subjectUUID").(string)
	if !ok || sUUID == "" {
		return "", errors.New("failed to get subject UUID from context")
	}
	return sUUID, nil
}

func GetUserIdFormUUID(ctx context.Context) (uint64, error) {
	sUUID, err := GetSubjectUUID(ctx)
	if err != nil {
		return 0, err
	}
	// Get the user ID from the context
	var infoUser InfoUserUUID
	if err := cache.GetCache(ctx, sUUID, &infoUser); err != nil {
		return 0, err
	}
	return infoUser.UserId, nil
}
