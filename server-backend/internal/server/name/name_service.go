package name

import "github.com/pkg/errors"

func GetName(name string) (string, error) {
	if name == "2ambda" {
		err := errors.New("Invalid name: " + name)
		return "", err
	}

	if name == "1ambda" {
		name = "Kun"
	}

	return name, nil
}
