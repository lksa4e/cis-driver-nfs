package token

import (
	"os"
	"strings"
	"testing"
)

func ConfigFromEnv() Config {
	var cfg Config

	cfg.Global.AuthURL = os.Getenv("OS_AUTH_URL")
	cfg.Global.UserID = os.Getenv("OS_USER_ID")
	cfg.Global.Password = os.Getenv("OS_PASSWORD")

	cfg.Global.TrustID = os.Getenv("OS_TRUST_ID")
	cfg.Global.Region = os.Getenv("OS_REGION_NAME")

	return cfg
}

func TestReadClouds(t *testing.T) {

	cfg, err := ReadConfig(strings.NewReader(`
[Global]
auth-url=http://169.254.169.241/v3
user-id=b67f4184f9874cfcae0bf96049f498d0
password=VjmF9M5A3urWHksVqP
trust-id=b220e21b4d624358ad731b88176f2eaa
region=KR2-T1
ca-file=/etc/kubernetes/ca-bundle.crt
`))

	if err != nil {
		t.Fatalf("Should succeed when a valid config is provided: %s", err)
	}

	// config has priority
	if cfg.Global.AuthURL != "http://169.254.169.241/v3" {
		t.Errorf("incorrect IdentityEndpoint: %s", cfg.Global.AuthURL)
	}

	if cfg.Global.UserID != "b67f4184f9874cfcae0bf96049f498d0" {
		t.Errorf("incorrect user-id: %s", cfg.Global.UserID)
	}

	if cfg.Global.Password != "VjmF9M5A3urWHksVqP" {
		t.Errorf("incorrect password: %s", cfg.Global.Password)
	}

	if cfg.Global.Region != "KR2-T1" {
		t.Errorf("incorrect region: %s", cfg.Global.Region)
	}

	if cfg.Global.TrustID != "b220e21b4d624358ad731b88176f2eaa" {
		t.Errorf("incorrect tenant name: %s", cfg.Global.TrustID)
	}

}
