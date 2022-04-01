package util

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(passwd string) (string, error) {
	bytePasswd, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytePasswd), nil
}

func VerifyPasswd(encrypted, origin string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(origin)) == nil
}
