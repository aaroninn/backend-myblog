package password

import "golang.org/x/crypto/bcrypt"

func HashePassword(pwd string) (string, error) {
	hashedpwd, err := bcrypt.GenerateFromPassword([]byte(pwd), 14)
	if err != nil {
		return "", err
	}

	return string(hashedpwd), nil
}

func ComparePassword(pwd, hashedpwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedpwd), []byte(pwd))
}
