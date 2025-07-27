package common

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func GenerateTransactionID(prefix string) string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 6)
	rand.Read(randomBytes)
	randomHex := hex.EncodeToString(randomBytes)

	if prefix != "" {
		return fmt.Sprintf("%s_%d_%s", prefix, timestamp, randomHex)
	}
	return fmt.Sprintf("TXN_%d_%s", timestamp, randomHex)
}

func GenerateReference(prefix string) string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomHex := strings.ToUpper(hex.EncodeToString(randomBytes))

	if prefix != "" {
		return fmt.Sprintf("%s_%d_%s", prefix, timestamp, randomHex)
	}
	return fmt.Sprintf("REF_%d_%s", timestamp, randomHex)
}

// SanitizeString removes or replaces potentially dangerous characters
func SanitizeString(input string) string {
	if input == "" {
		return ""
	}

	// Remove control characters
	var result strings.Builder
	for _, r := range input {
		if unicode.IsControl(r) && r != '\t' {
			continue // Skip control characters except tab
		}
		result.WriteRune(r)
	}

	sanitized := result.String()

	// Remove newlines and carriage returns
	sanitized = strings.ReplaceAll(sanitized, "\n", " ")
	sanitized = strings.ReplaceAll(sanitized, "\r", " ")
	sanitized = strings.ReplaceAll(sanitized, "\t", " ")

	// Trim and collapse multiple spaces
	sanitized = strings.TrimSpace(sanitized)
	spaceRegex := regexp.MustCompile(`\s+`)
	sanitized = spaceRegex.ReplaceAllString(sanitized, " ")

	return sanitized
}

// FormatAmount formats monetary amount for display
func FormatAmount(amount float64, currency string) string {
	return fmt.Sprintf("%.2f %s", amount, currency)
}

// FormatAmountWithSeparators formats amount with thousands separators
func FormatAmountWithSeparators(amount float64, currency string) string {
	// Convert to string with 2 decimal places
	amountStr := fmt.Sprintf("%.2f", amount)

	// Split into integer and decimal parts
	parts := strings.Split(amountStr, ".")
	integerPart := parts[0]
	decimalPart := parts[1]

	// Add thousands separators
	if len(integerPart) > 3 {
		var result []string
		for i, char := range reverse(integerPart) {
			if i > 0 && i%3 == 0 {
				result = append(result, ",")
			}
			result = append(result, string(char))
		}
		integerPart = reverse(strings.Join(result, ""))
	}

	return fmt.Sprintf("%s.%s %s", integerPart, decimalPart, currency)
}

// reverse reverses a string
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ParseDuration parses duration with fallback
func ParseDuration(duration string, fallback time.Duration) time.Duration {
	if duration == "" {
		return fallback
	}

	if d, err := time.ParseDuration(duration); err == nil {
		return d
	}

	// Try parsing as seconds
	if seconds, err := strconv.Atoi(duration); err == nil {
		return time.Duration(seconds) * time.Second
	}

	return fallback
}

// ParseAmount parses amount string to float64
func ParseAmount(amountStr string) (float64, error) {
	if amountStr == "" {
		return 0, fmt.Errorf("amount string is empty")
	}

	// Remove any currency symbols and spaces
	cleaned := strings.TrimSpace(amountStr)
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	// Remove common currency symbols
	currencySymbols := []string{"$", "€", "£", "¥", "₹", "MRU", "MRO"}
	for _, symbol := range currencySymbols {
		cleaned = strings.ReplaceAll(cleaned, symbol, "")
	}

	return strconv.ParseFloat(cleaned, 64)
}

// Hash generates SHA256 hash of input
func Hash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// HashWithSalt generates SHA256 hash with salt
func HashWithSalt(input, salt string) string {
	return Hash(input + salt)
}

// GenerateSalt generates a random salt
func GenerateSalt(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

// ToJSON converts interface to JSON string
func ToJSON(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

// FromJSON parses JSON string to interface
func FromJSON(jsonStr string, v interface{}) error {
	return json.Unmarshal([]byte(jsonStr), v)
}

// IsValidUUID checks if string is valid UUID
func IsValidUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	return uuidRegex.MatchString(uuid)

}

func SliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// SliceUnique removes duplicates from slice
func SliceUnique(slice []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// MapKeys returns keys from string map
func MapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// ContainsAny checks if string contains any of the substrings
func ContainsAny(s string, substrings []string) bool {
	for _, substring := range substrings {
		if strings.Contains(strings.ToLower(s), strings.ToLower(substring)) {
			return true
		}
	}
	return false
}

// IsNumeric checks if string contains only numeric characters
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}

	for _, char := range s {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

// IsAlphaNumeric checks if string contains only alphanumeric characters
func IsAlphaNumeric(s string) bool {
	if s == "" {
		return false
	}

	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

// Capitalize capitalizes first letter of string
func Capitalize(s string) string {
	if s == "" {
		return ""
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// CamelToSnake converts camelCase to snake_case
func CamelToSnake(s string) string {
	var result []rune

	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}

	return string(result)
}

// SnakeToCamel converts snake_case to camelCase
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return s
	}

	result := parts[0]
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			result += Capitalize(parts[i])
		}
	}

	return result
}

// GetMapValue safely gets value from map with default
func GetMapValue(m map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}

// GetMapString safely gets string value from map
func GetMapString(m map[string]interface{}, key string) string {
	if value, exists := m[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetMapInt safely gets int value from map
func GetMapInt(m map[string]interface{}, key string) int {
	if value, exists := m[key]; exists {
		switch v := value.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

// GetMapFloat safely gets float64 value from map
func GetMapFloat(m map[string]interface{}, key string) float64 {
	if value, exists := m[key]; exists {
		switch v := value.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
	}
	return 0.0
}

// GetMapBool safely gets bool value from map
func GetMapBool(m map[string]interface{}, key string) bool {
	if value, exists := m[key]; exists {
		switch v := value.(type) {
		case bool:
			return v
		case string:
			return v == "true" || v == "1" || v == "yes"
		case int:
			return v != 0
		case float64:
			return v != 0
		}
	}
	return false
}

// TimeElapsed returns human readable time elapsed
func TimeElapsed(start time.Time) string {
	elapsed := time.Since(start)

	if elapsed < time.Minute {
		return fmt.Sprintf("%.0f seconds", elapsed.Seconds())
	} else if elapsed < time.Hour {
		return fmt.Sprintf("%.0f minutes", elapsed.Minutes())
	} else if elapsed < 24*time.Hour {
		return fmt.Sprintf("%.1f hours", elapsed.Hours())
	} else {
		days := elapsed.Hours() / 24
		return fmt.Sprintf("%.1f days", days)
	}
}

// FormatFileSize formats file size in human readable format
func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.1f %s", float64(size)/float64(div), units[exp])
}
