package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLen   = 8
	maxPasswordBytes = 72 // bcrypt игнорирует все байты после 72-го
)

var ErrWeakPassword = errors.New("password must be 8-72 bytes")

func ValidatePassword(p string) error {
	if len(p) < minPasswordLen || len(p) > maxPasswordBytes {
		return ErrWeakPassword
	}
	return nil
}

func HashPassword(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

// CheckPassword возвращает false и для неверного пароля, и для строки,
// которая не является bcrypt-хешем (например, старый пароль в открытом виде).
func CheckPassword(hash, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil
}

// HashOptional валидирует и хеширует пароль, если он передан.
// nil или пустая строка возвращаются как nil: при create это NULL,
// при update COALESCE оставит старый пароль.

func HashOptional(p *string) (*string, error) {
	if p == nil || *p == "" {
		return nil, nil
	}

	if err := ValidatePassword(*p); err != nil {
		return nil, err
	}

	h, err := HashPassword(*p)

	if err != nil {
		return nil, err
	}

	return &h, nil
}
