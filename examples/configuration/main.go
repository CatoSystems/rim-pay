package main

import (
	"fmt"
	"log"
	"time"

	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

func main() {
	fmt.Println("🏦 RimPay Library - Configuration Examples")
	fmt.Println("==========================================")

	// Example 1: Default Configuration
	fmt.Println("\n⚙️ Example 1: Default Configuration")
	demonstrateDefaultConfig()

	// Example 2: Production Configuration
	fmt.Println("\n🏭 Example 2: Production Configuration")
	demonstrateProductionConfig()

	// Example 3: Custom Timeouts and Limits
	fmt.Println("\n⏱️ Example 3: Custom Timeouts and Limits")
	demonstrateCustomTimeouts()

	// Example 4: Environment-specific Configuration
	fmt.Println("\n🌍 Example 4: Environment-specific Configuration")
	demonstrateEnvironmentConfig()

	fmt.Println("\n💡 Configuration Features Demonstrated:")
	fmt.Println("✅ Default vs custom configurations")
	fmt.Println("✅ Environment-specific settings")
	fmt.Println("✅ Provider-specific configurations")
	fmt.Println("✅ Timeout and connection management")
	fmt.Println("✅ Logging configuration")
}

func demonstrateDefaultConfig() {
	// Get default configuration
	config := rimpay.DefaultConfig()

	fmt.Printf("   Environment: %s\n", config.Environment)
	fmt.Printf("   Default Provider: %s\n", config.DefaultProvider)
	fmt.Printf("   Available Providers: %d\n", len(config.Providers))

	// Show default HTTP settings
	fmt.Printf("   HTTP Timeout: %v\n", config.HTTP.Timeout)
	fmt.Printf("   Max Idle Connections: %d\n", config.HTTP.MaxIdleConns)
	fmt.Printf("   Max Connections per Host: %d\n", config.HTTP.MaxConnsPerHost)

	// Show logging settings
	fmt.Printf("   Log Level: %s\n", config.Logging.Level)
	fmt.Printf("   Log Format: %s\n", config.Logging.Format)
}

func demonstrateProductionConfig() {
	config := &rimpay.Config{
		Environment:     rimpay.EnvironmentProduction,
		DefaultProvider: "bpay",
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://api.bpay.mr/v1", // Production URL
				Timeout: 45 * time.Second,         // Longer timeout for production
				Credentials: map[string]string{
					"username":  "production_user",
					"password":  "secure_production_password",
					"client_id": "prod_client_12345",
				},
			},
			"masrvi": {
				Enabled: true,
				BaseURL: "https://masrviapp.mr/api", // Production URL
				Timeout: 60 * time.Second,           // Even longer for MASRVI
				Credentials: map[string]string{
					"merchant_id": "PROD_MERCHANT_789",
				},
			},
		},
		HTTP: rimpay.HTTPConfig{
			Timeout:         60 * time.Second, // Production timeout
			MaxIdleConns:    50,               // Higher connection pool
			MaxConnsPerHost: 20,               // More connections per host
		},
		Logging: rimpay.LoggingConfig{
			Level:  "warn", // Less verbose in production
			Format: "json", // Structured logging
		},
	}

	fmt.Printf("   ✅ Production configuration created\n")
	fmt.Printf("   Environment: %s\n", config.Environment)
	fmt.Printf("   HTTP Timeout: %v\n", config.HTTP.Timeout)
	fmt.Printf("   Connection Pool: %d max idle, %d per host\n",
		config.HTTP.MaxIdleConns, config.HTTP.MaxConnsPerHost)
	fmt.Printf("   Logging: %s level, %s format\n",
		config.Logging.Level, config.Logging.Format)

	// Create client with production config
	client, err := rimpay.NewClient(config)
	if err != nil {
		fmt.Printf("   ❌ Failed to create production client: %v\n", err)
		return
	}

	fmt.Printf("   ✅ Production client created successfully\n")
	_ = client // Use client for actual operations
}

func demonstrateCustomTimeouts() {
	config := rimpay.DefaultConfig()

	// Customize provider timeouts
	bpayConfig := config.Providers["bpay"]
	bpayConfig.Timeout = 15 * time.Second // Shorter timeout for B-PAY
	config.Providers["bpay"] = bpayConfig

	masrviConfig := config.Providers["masrvi"]
	masrviConfig.Timeout = 90 * time.Second // Longer timeout for MASRVI
	config.Providers["masrvi"] = masrviConfig

	// Customize HTTP settings
	config.HTTP.Timeout = 30 * time.Second
	config.HTTP.MaxIdleConns = 25
	config.HTTP.MaxConnsPerHost = 10

	fmt.Printf("   B-PAY Timeout: %v\n", config.Providers["bpay"].Timeout)
	fmt.Printf("   MASRVI Timeout: %v\n", config.Providers["masrvi"].Timeout)
	fmt.Printf("   HTTP Timeout: %v\n", config.HTTP.Timeout)
	fmt.Printf("   Connection Limits: %d idle, %d per host\n",
		config.HTTP.MaxIdleConns, config.HTTP.MaxConnsPerHost)

	client, err := rimpay.NewClient(config)
	if err != nil {
		fmt.Printf("   ❌ Failed to create client: %v\n", err)
		return
	}

	fmt.Printf("   ✅ Client with custom timeouts created\n")
	_ = client
}

func demonstrateEnvironmentConfig() {
	environments := []struct {
		name string
		env  rimpay.Environment
	}{
		{"Development", rimpay.EnvironmentSandbox},
		{"Production", rimpay.EnvironmentProduction},
	}

	for _, envConfig := range environments {
		fmt.Printf("   🌍 %s Environment:\n", envConfig.name)

		config := createEnvironmentConfig(envConfig.env)

		fmt.Printf("      Environment: %s\n", config.Environment)
		fmt.Printf("      Log Level: %s\n", config.Logging.Level)
		fmt.Printf("      HTTP Timeout: %v\n", config.HTTP.Timeout)

		// Show provider URLs for this environment
		for providerName, providerConfig := range config.Providers {
			fmt.Printf("      %s URL: %s\n", providerName, providerConfig.BaseURL)
		}
		fmt.Println()
	}
}

func createEnvironmentConfig(env rimpay.Environment) *rimpay.Config {
	config := rimpay.DefaultConfig()
	config.Environment = env

	switch env {
	case rimpay.EnvironmentSandbox:
		// Development/Testing settings
		config.Logging.Level = "debug"
		config.HTTP.Timeout = 30 * time.Second

		// Use test URLs
		bpayConfig := config.Providers["bpay"]
		bpayConfig.BaseURL = "https://ebankily-tst.appspot.com"
		config.Providers["bpay"] = bpayConfig

		masrviConfig := config.Providers["masrvi"]
		masrviConfig.BaseURL = "https://test.masrviapp.mr/online"
		config.Providers["masrvi"] = masrviConfig

	case rimpay.EnvironmentProduction:
		// Production settings
		config.Logging.Level = "info"
		config.Logging.Format = "json"
		config.HTTP.Timeout = 60 * time.Second
		config.HTTP.MaxIdleConns = 50
		config.HTTP.MaxConnsPerHost = 20

		// Use production URLs
		bpayConfig := config.Providers["bpay"]
		bpayConfig.BaseURL = "https://api.bpay.mr/v1"
		bpayConfig.Timeout = 45 * time.Second
		config.Providers["bpay"] = bpayConfig

		masrviConfig := config.Providers["masrvi"]
		masrviConfig.BaseURL = "https://masrviapp.mr/api"
		masrviConfig.Timeout = 90 * time.Second
		config.Providers["masrvi"] = masrviConfig
	}

	return config
}

func init() {
	// This example doesn't actually create clients to avoid credential issues
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
