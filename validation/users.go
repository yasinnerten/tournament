package validation

import (
	"errors"
	"tournament-app/model"
)

func ValidateUser(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	if user.Name == "" {
		return errors.New("user name cannot be empty")
	}
	if user.Level < 0 || user.Level > 100 {
		return errors.New("user level must be between 0 and 100")
	}
	if user.Money < 0 {
		return errors.New("user money cannot be negative")
	}

	return nil
}
