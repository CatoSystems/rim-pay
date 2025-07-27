package main

import (
	"fmt"
	"time"

	"github.com/CatoSystems/rim-pay/pkg/rimpay"

	// Import providers
)

func main() {
	fmt.Println("🏦 RimPay Library - Configuration Examples")
	fmt.Println("==========================================\n")

	// Example 1: Basic configuration
	fmt.Println("📋 Example 1: Basic Configuration")
	basicConfig := createBasicConfig()
	showConfig("Basic", basicConfig)

	// Example 2: Production configuration
	fmt.Println("\n🏭 Example 2: Production Configuration")
	prodConfig := createProductionConfig()
	showConfig("Production", prodConfig)

	// Example 3: Development configuration with debugging
	fmt.Println("\n🔧 Example 3: Development Configuration")
	devConfig := createDevelopmentConfig()
	showConfig("Development", devConfig)

	// Example 4: High-performance configuration
	fmt.Println("\n⚡ Example 4: High-Performance Configuration")
	perfConfig := createHighPerformanceConfig()
	showConfig("High-Performance", perfConfig)

	// Example 5: Configuration validation
	fmt.Println("\n✅ Example 5: Configuration Validation")
	demonstrateConfigValidation()

	// Example 6: Environment-specific settings
	fmt.Println("\n🌍 Example 6: Environment-Specific Settings")
	demonstrateEnvironmentConfigs()

	fmt.Println("\n💡 Configuration Features Demonstrated:")
	fmt.Println("✅ Basic and advanced configuration options")
	fmt.Println("✅ Environment-specific settings")
	fmt.Println("✅ HTTP client tuning for performance")
	fmt.Println("✅ Logging and security configuration")
	fmt.Println("✅ Provider-specific credential management")
	fmt.Println("✅ Configuration validation and error handling")
}

func createBasicConfig() *rimpay.Config {
	return &rimpay.Config{
		Environment:     rimpay.EnvironmentSandbox,
		DefaultProvider: "bpay",
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://ebankily-tst.appspot.com",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{
					"username":  "your_username",
					"password":  "your_password",
					"client_id": "your_client_id",
				},
			},
		},
	}
}

func createProductionConfig() *rimpay.Config {
	return &rimpay.Config{
		Environment:     rimpay.EnvironmentProduction,
		DefaultProvider: "bpay",
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://api.bpay.mr", // Production URL
				Timeout: 60 * time.Second,      // Longer timeout for production
				Credentials: map[string]string{
					"username":  "prod_username",
					"password":  "prod_password",
					"client_id": "prod_client_id",
				},
			},
			"masrvi": {
				Enabled: true,
				BaseURL: "https://masrviapp.mr/online",
				Timeout: 45 * time.Second,
				Credentials: map[string]string{
					"merchant_id": "prod_merchant_id",
				},
			},
		},
		HTTP: rimpay.HTTPConfig{
			Timeout:         60 * time.Second,
			MaxIdleConns:    50,
			MaxConnsPerHost: 20,
		},
		Logging: rimpay.LoggingConfig{
			Level:  "warn", // Less verbose in production
			Format: "json",
		},
		Security: rimpay.SecurityConfig{
			EncryptionKey: "production-encryption-key",
			SigningKey:    "production-signing-key",
			TokenTTL:      2 * time.Hour,
		},
	}
}

func createDevelopmentConfig() *rimpay.Config {
	return &rimpay.Config{
		Environment:     rimpay.EnvironmentSandbox,
		DefaultProvider: "bpay",
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://ebankily-tst.appspot.com",
				Timeout: 10 * time.Second, // Shorter timeout for faster feedback
				Credentials: map[string]string{
					"username":  "dev_username",
					"password":  "dev_password",
					"client_id": "dev_client_id",
				},
			},
			"masrvi": {
				Enabled: true,
				BaseURL: "https://masrviapp.mr/online",
				Timeout: 10 * time.Second,
				Credentials: map[string]string{
					"merchant_id": "dev_merchant_id",
				},
			},
		},
		HTTP: rimpay.HTTPConfig{
			Timeout:         15 * time.Second,
			MaxIdleConns:    5,
			MaxConnsPerHost: 2,
		},
		Logging: rimpay.LoggingConfig{
			Level:  "debug", // Verbose logging for development
			Format: "text",  // Human-readable format
		},
		Security: rimpay.SecurityConfig{
			EncryptionKey: "dev-encryption-key",
			SigningKey:    "dev-signing-key",
			TokenTTL:      30 * time.Minute,
		},
	}
}

func createHighPerformanceConfig() *rimpay.Config {
	return &rimpay.Config{
		Environment:     rimpay.EnvironmentProduction,
		DefaultProvider: "bpay",
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://api.bpay.mr",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{
					"username":  "perf_username",
					"password":  "perf_password",
					"client_id": "perf_client_id",
				},
			},
			"masrvi": {
				Enabled: true,
				BaseURL: "https://masrviapp.mr/online",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{
					"merchant_id": "perf_merchant_id",
				},
			},
		},
		HTTP: rimpay.HTTPConfig{
			Timeout:         30 * time.Second,
			MaxIdleConns:    100, // High connection pooling
			MaxConnsPerHost: 50,  // Many concurrent connections
		},
		Logging: rimpay.LoggingConfig{
			Level:  "error", // Minimal logging for performance
			Format: "json",
		},
		Security: rimpay.SecurityConfig{
			EncryptionKey: "perf-encryption-key",
			SigningKey:    "perf-signing-key",
			TokenTTL:      4 * time.Hour,
		},
	}
}

func showConfig(name string, config *rimpay.Config) {
	fmt.Printf("   Environment: %s\n", config.Environment)
	fmt.Printf("   Default Provider: %s\n", config.DefaultProvider)
	fmt.Printf("   Providers: %v\n", getProviderNames(config))
	
	if config.HTTP.Timeout > 0 {
		fmt.Printf("   HTTP Timeout: %v\n", config.HTTP.Timeout)
		fmt.Printf("   Max Connections: %d\n", config.HTTP.MaxIdleConns)
	}
	
	if config.Logging.Level != "" {
		fmt.Printf("   Log Level: %s\n", config.Logging.Level)
		fmt.Printf("   Log Format: %s\n", config.Logging.Format)
	}
	
	if config.Security.EncryptionKey != "" {
		fmt.Printf("   Security: Encryption enabled, Token TTL: %v\n", 
			config.Security.TokenTTL)
	}

	// Test the configuration
	fmt.Printf("   Validation: ")
	if err := config.Validate(); err != nil {
		fmt.Printf("❌ Invalid - %v\n", err)
	} else {
		fmt.Printf("✅ Valid\n")
	}
}

func getProviderNames(config *rimpay.Config) []string {
	var names []string
	for name, provider := range config.Providers {
		if provider.Enabled {
			names = append(names, name)
		}
	}
	return names
}

func demonstrateConfigValidation() {
	fmt.Printf("   Testing various invalid configurations...\n\n")

	// Test 1: Missing environment
	fmt.Printf("   Test 1: Missing environment\n")
	invalidConfig1 := &rimpay.Config{
		DefaultProvider: "bpay",
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://example.com",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{"username": "test"},
			},
		},
	}
	testConfigValidation(invalidConfig1)

	// Test 2: Invalid environment
	fmt.Printf("\n   Test 2: Invalid environment\n")
	invalidConfig2 := &rimpay.Config{
		Environment:     "invalid",
		DefaultProvider: "bpay",
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://example.com",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{"username": "test"},
			},
		},
	}
	testConfigValidation(invalidConfig2)

	// Test 3: Missing default provider
	fmt.Printf("\n   Test 3: Missing default provider\n")
	invalidConfig3 := &rimpay.Config{
		Environment: rimpay.EnvironmentSandbox,
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://example.com",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{"username": "test"},
			},
		},
	}
	testConfigValidation(invalidConfig3)

	// Test 4: Invalid timeout
	fmt.Printf("\n   Test 4: Invalid timeout\n")
	invalidConfig4 := &rimpay.Config{
		Environment:     rimpay.EnvironmentSandbox,
		DefaultProvider: "bpay",
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://example.com",
				Timeout: -5 * time.Second, // Negative timeout
				Credentials: map[string]string{"username": "test"},
			},
		},
	}
	testConfigValidation(invalidConfig4)
}

func testConfigValidation(config *rimpay.Config) {
	if err := config.Validate(); err != nil {
		fmt.Printf("     ❌ Validation failed: %v\n", err)
	} else {
		fmt.Printf("     ✅ Validation passed\n")
	}
}

func demonstrateEnvironmentConfigs() {
	environments := map[string]*rimpay.Config{
		"Local Development": {
			Environment:     rimpay.EnvironmentSandbox,
			DefaultProvider: "bpay",
			Providers: map[string]rimpay.ProviderConfig{
				"bpay": {
					Enabled: true,
					BaseURL: "http://localhost:8080", // Local mock server
					Timeout: 5 * time.Second,
					Credentials: map[string]string{
						"username":  "local_user",
						"password":  "local_pass",
						"client_id": "local_client",
					},
				},
			},
			Logging: rimpay.LoggingConfig{Level: "debug", Format: "text"},
		},
		"Staging": {
			Environment:     rimpay.EnvironmentSandbox,
			DefaultProvider: "bpay",
			Providers: map[string]rimpay.ProviderConfig{
				"bpay": {
					Enabled: true,
					BaseURL: "https://staging-api.bpay.mr",
					Timeout: 30 * time.Second,
					Credentials: map[string]string{
						"username":  "staging_user",
						"password":  "staging_pass",
						"client_id": "staging_client",
					},
				},
			},
			Logging: rimpay.LoggingConfig{Level: "info", Format: "json"},
		},
		"Production": {
			Environment:     rimpay.EnvironmentProduction,
			DefaultProvider: "bpay",
			Providers: map[string]rimpay.ProviderConfig{
				"bpay": {
					Enabled: true,
					BaseURL: "https://api.bpay.mr",
					Timeout: 60 * time.Second,
					Credentials: map[string]string{
						"username":  "prod_user",
						"password":  "prod_pass",
						"client_id": "prod_client",
					},
				},
			},
			HTTP: rimpay.HTTPConfig{
				Timeout:         60 * time.Second,
				MaxIdleConns:    50,
				MaxConnsPerHost: 20,
			},
			Logging: rimpay.LoggingConfig{Level: "warn", Format: "json"},
			Security: rimpay.SecurityConfig{
				EncryptionKey: "prod-encryption-key",
				SigningKey:    "prod-signing-key",
				TokenTTL:      4 * time.Hour,
			},
		},
	}

	for envName, config := range environments {
		fmt.Printf("\n   %s Environment:\n", envName)
		fmt.Printf("     Base URL: %s\n", config.Providers["bpay"].BaseURL)
		fmt.Printf("     Timeout: %v\n", config.Providers["bpay"].Timeout)
		fmt.Printf("     Log Level: %s\n", config.Logging.Level)
		
		if config.Security.EncryptionKey != "" {
			fmt.Printf("     Security: Enhanced (Encryption + Signing)\n")
		} else {
			fmt.Printf("     Security: Basic\n")
		}
		
		// Create client to test configuration
		client, err := rimpay.NewClient(config)
		if err != nil {
			fmt.Printf("     Status: ❌ Invalid - %v\n", err)
		} else {
			fmt.Printf("     Status: ✅ Valid client created\n")
			// Don't forget to close/cleanup if needed
			_ = client
		}
	}

	fmt.Printf("\n   💡 Environment Best Practices:\n")
	fmt.Printf("   🏠 Local: Short timeouts, debug logging, mock servers\n")
	fmt.Printf("   🎭 Staging: Real APIs, info logging, similar to prod\n")
	fmt.Printf("   🏭 Production: Long timeouts, minimal logging, full security\n")
}
