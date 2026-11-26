package controller

import (
	"errors"
	"regexp"
)

var pricePattern = regexp.MustCompile(`^\d{1,8}(\.\d{1,2})?$`)

func validateProduct(price string, stock int) (string, error) {
	if price == "" {
		price = "0.00"
	}

	if !pricePattern.MatchString(price) {
		return "", errors.New("price must be a number with up to 2 decimals (max 99999999.99)")
	}

	if stock < 0 {
		return "", errors.New("stock must be non-negative")
	}

	return price, nil
}
