package hw10programoptimization

import (
	"bufio"
	"io"
	"log"
	"strings"

	"github.com/mailru/easyjson" //nolint:depguard
)

//go:generate easyjson -all stats.go
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := DomainStat{}
	var user User

	for scanner.Scan() {
		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return result, err
		}
		if strings.HasSuffix(user.Email, domain) {
			if b, emailDomain := compareDomain(user.Email, domain); b {
				result[emailDomain]++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result, nil
}

func compareDomain(email string, domain string) (bool, string) {
	index := strings.IndexRune(email, '@')
	if index != -1 && index < len(email)-1 {
		emailDomain := strings.ToLower(email[index+1:])
		if strings.HasSuffix(emailDomain, domain) {
			return true, emailDomain
		}
	}
	return false, ""
}
