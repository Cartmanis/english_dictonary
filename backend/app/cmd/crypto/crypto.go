package crypto

import (
	"crypto/md5"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"math/rand"
	"os"
	"time"
)

const (
	letter = 1
	bigLetter
	digital
	symbol
)

var (
	salt = os.Getenv("SALT")
)

func GetHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashPassword(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

func CreateToken() (string, error) {
	fmt.Println("salt: ", salt)
	h := md5.New()
	_, err := io.WriteString(h, salt+time.Now().String())
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func RandPassword(count byte) string {
	var allRunes = []rune(`abcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*,.?<>\|/)(~-_=+[]{}'":;ABCDEFGHIJKLMNOPRSTUVWXYZ`)
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	var letterBigRunes = []rune("ABCDEFGHIJKLMNOPRSTUVWXYZ")
	var digitalRunes = []rune("0123456789")
	var sybmulRunes = []rune(`!@#$%^&*,.?<>\|/)(~-_=+[]{}'":;`)

	rand.Seed(time.Now().UnixNano())
	password := make([]rune, count)

	mask := make([]int, count)

	//Для значений пароля от 17 и выше задется просто случайные значения из все возможных символов. Теоритически тут могут быть только цифры или буквы, но маловероятно
	if count > 16 {
		for i := range password {
			password[i] = allRunes[rand.Intn(len(allRunes))]
		}

	}

	for i := range mask {
		for {
			rand.Seed(time.Now().UnixNano())
			mask[i] = rand.Intn(4) + 1
			if validMask(mask, count) {
				break
			}
		}
		switch mask[i] {
		case 1:
			password[i] = letterRunes[rand.Intn(len(letterRunes))]
		case 2:
			password[i] = letterBigRunes[rand.Intn(len(letterBigRunes))]
		case 3:
			password[i] = digitalRunes[rand.Intn(len(digitalRunes))]
		case 4:
			password[i] = sybmulRunes[rand.Intn(len(sybmulRunes))]
		}
	}
	return string(password)
}

//функция обеспечивает правильное количество букв, букв в верхнем регистре цифр и символов для паролей разной длинны
func validMask(mask []int, count byte) bool {
	if len(mask) < 3 {
		return true
	}
	countLetter := 0
	countBigLetter := 0
	countDigital := 0
	countSymbol := 0
	for _, v := range mask {
		switch v {
		case 1:
			countLetter++
			break
		case 2:
			countBigLetter++
			break
		case 3:
			countDigital++
			break
		case 4:
			countSymbol++
			break
		}
	}

	if count < 6 {
		if countLetter > 2 || countBigLetter > 1 || countDigital > 1 || countSymbol > 1 {
			return false
		}
		return true
	}

	switch count {
	case 6:
		if countLetter > 2 || countBigLetter > 2 || countDigital > 1 || countSymbol > 1 {
			return false
		}
		return true
	case 7:
		if countLetter > 2 || countBigLetter > 2 || countDigital > 2 || countSymbol > 1 {
			return false
		}
		return true
	case 8:
		if countLetter > 2 || countBigLetter > 2 || countDigital > 2 || countSymbol > 2 {
			return false
		}
		return true
	case 9:
		if countLetter > 3 || countBigLetter > 2 || countDigital > 2 || countSymbol > 2 {
			return false
		}
		return true
	case 10:
		if countLetter > 4 || countBigLetter > 2 || countDigital > 2 || countSymbol > 2 {
			return false
		}
		return true
	case 11:
		if countLetter > 4 || countBigLetter > 3 || countDigital > 2 || countSymbol > 2 {
			return false
		}
		return true
	case 12:
		if countLetter > 4 || countBigLetter > 3 || countDigital > 3 || countSymbol > 2 {
			return false
		}
		return true
	case 13:
		if countLetter > 5 || countBigLetter > 3 || countDigital > 3 || countSymbol > 2 {
			return false
		}
		return true
	case 14:
		if countLetter > 5 || countBigLetter > 3 || countDigital > 3 || countSymbol > 3 {
			return false
		}
		return true
	case 15:
		if countLetter > 5 || countBigLetter > 4 || countDigital > 3 || countSymbol > 3 {
			return false
		}
		return true
	case 16:
		if countLetter > 6 || countBigLetter > 4 || countDigital > 3 || countSymbol > 3 {
			return false
		}
		return true
	}

	return true
}

//func getMask(count int) []int {
//	mask := make([]int, count)
//	for i := range mask {
//		for {
//			rand.Seed(time.Now().UnixNano())
//			mask[i] = rand.Intn(4) + 1
//			if validMask(mask, count) {
//				break
//			}
//		}
//	}
//	return mask
//}
