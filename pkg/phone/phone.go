package phone

import (
	"fmt"
	"regexp"
	"strings"
)

type Operator string

const (
	Mauritel   Operator = "mauritel"   // 2-5
	Mattel     Operator = "mattel"     // 6-7
	Chinguitel Operator = "chinguitel" // 8-9
)

type Phone struct {
	number   string
	operator Operator
}

var mauritanianPattern = regexp.MustCompile(`^(\+222|00222|222)?([2-9]\d{7})$`)

func NewPhone(number string) (*Phone, error) {
	if number == "" {
		return nil, fmt.Errorf("phone number required")
	}

	cleaned := cleanPhoneNumber(number)
	if !IsValidMauritanianNumber(cleaned) {
		return nil, fmt.Errorf("invalid Mauritanian phone number: %s", number)
	}

	localNumber := extractLocalNumber(cleaned)
	operator := determineOperator(localNumber)
	if operator == "" {
		return nil, fmt.Errorf("unknown operator: %s", number)
	}

	return &Phone{
		number:   localNumber,
		operator: operator,
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

func determineOperator(localNumber string) Operator {
	if len(localNumber) != 8 {
		return ""
	}

	switch localNumber[0] {
	case '2', '3', '4', '5':
		return Mauritel
	case '6', '7':
		return Mattel
	case '8', '9':
		return Chinguitel
	default:
		return ""
	}
}

func IsValidMauritanianNumber(number string) bool {
	cleaned := cleanPhoneNumber(number)
	return mauritanianPattern.MatchString(cleaned)
}

func (mp *Phone) Number() string      { return mp.number }
func (mp *Phone) Operator() Operator  { return mp.operator }
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
