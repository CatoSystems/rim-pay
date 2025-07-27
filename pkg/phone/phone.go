package phone

import (
	"fmt"
	"regexp"
	"strings"
)

type Phone struct {
	number   string
}

var mauritanianPattern = regexp.MustCompile(`^(\+222|00222|222)?([234]\d{7})$`)

func NewPhone(number string) (*Phone, error) {
	if number == "" {
		return nil, fmt.Errorf("phone number required")
	}

	cleaned := cleanPhoneNumber(number)
	if !IsValidMauritanianNumber(cleaned) {
		return nil, fmt.Errorf("invalid Mauritanian phone number: %s", number)
	}

	localNumber := extractLocalNumber(cleaned)
	
	return &Phone{
		number:   localNumber,
	}, nil
}

func cleanPhoneNumber(number string) string {
	cleaned := strings.ReplaceAll(number, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")
	return cleaned
}

func extractLocalNumber(number string) string {
	matches := mauritanianPattern.FindStringSubmatch(number)
	if len(matches) >= 3 {
		return matches[2]
	}
	if len(number) == 8 && regexp.MustCompile(`^[2-9]\d{7}$`).MatchString(number) {
		return number
	}
	return ""
}



func IsValidMauritanianNumber(number string) bool {
	cleaned := cleanPhoneNumber(number)
	return mauritanianPattern.MatchString(cleaned)
}

func (mp *Phone) Number() string      { return mp.number }
func (mp *Phone) String() string      { return fmt.Sprintf("+222%s", mp.number) }
func (mp *Phone) LocalFormat() string { return mp.number }
func (mp *Phone) InternationalFormat() string {
	return fmt.Sprintf("+222 %s %s %s", mp.number[:2], mp.number[2:5], mp.number[5:])
}

func (mp *Phone) ForProvider(includeCountryCode bool) string {
	if includeCountryCode {
		return fmt.Sprintf("222%s", mp.number)
	}
	return mp.number
}
