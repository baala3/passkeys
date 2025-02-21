package concerns

import (
	"github.com/baala3/passkeys/model"
	"github.com/baala3/passkeys/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func CurrentUser(ctx echo.Context, userRepository repository.UserRepository) *model.User {
	userID := ctx.Get("userID").(string)
		if userID == "" {
			return nil
		}

		parsedUUID, err := uuid.Parse(userID)
		if err != nil {
			return nil
		}

		userIDBytes, err := parsedUUID.MarshalBinary()
		if err != nil {
			return nil
		}

	user, err := userRepository.FindUserById(ctx.Request().Context(), userIDBytes)
	if err != nil {
		return nil
	}
		return user
}
