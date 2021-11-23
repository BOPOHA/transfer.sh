package server

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os"
	"testing"
)

func cleanAwsKeyEnviron() {
	os.Setenv("AWS_CONFIG_FILE", "fake_file")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "fake_file")
	os.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	os.Unsetenv("AWS_PROFILE")
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
}

func TestNewSession_NoCredentials(t *testing.T) {
	cleanAwsKeyEnviron()

	staticAccessKey := ""
	staticSecretKey := ""

	s := GetAwsSession(staticAccessKey, staticSecretKey, "", "", false)

	if s.Config.Credentials == nil {
		t.Errorf("expect not nil")
	}

	configCreds, err := s.Config.Credentials.Get()

	if err != credentials.ErrStaticCredentialsEmpty {
		t.Errorf("expect error EmptyStaticCreds, got %+v", err)
	}

	if configCreds.ProviderName != "StaticProvider" {
		t.Errorf("expect not nil")
	}

	anonCreds, _ := credentials.AnonymousCredentials.Get()

	if e, a := anonCreds, configCreds; e != a {
		t.Errorf("expect the same credentials,\nac: %+v\ncc: %+v\n", e, a)
	}

}

func TestNewSession_StaticCredentials(t *testing.T) {
	cleanAwsKeyEnviron()
	staticAccessKey := "AKIA"
	staticSecretKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

	s := GetAwsSession(staticAccessKey, staticSecretKey, "", "", false)

	if s.Config.Credentials == nil {
		t.Errorf("expect not nil")
	}

	configCreds, err := s.Config.Credentials.Get()

	if err != nil {
		t.Errorf("expect nil")
	}

	if configCreds.ProviderName != "StaticProvider" {
		t.Errorf("expect StaticProvider")
	}
	if configCreds.AccessKeyID != staticAccessKey {
		t.Errorf("expect staticAccessKey")
	}
	if configCreds.SecretAccessKey != staticSecretKey {
		t.Errorf("expect staticSecretKey")
	}
	if configCreds.SessionToken != "" {
		t.Errorf("expect empty string")
	}

}

func TestNewSession_SharedCredentials(t *testing.T) {
	cleanAwsKeyEnviron()
	staticAccessKey := "AKIA"
	staticSecretKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	os.Setenv("AWS_ACCESS_KEY_ID", staticAccessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", staticSecretKey)

	s, err := GetSharedConfigSession("", "", false)
	if err != nil {
		t.Errorf("expect nil, %configCreds", err)
	}

	if s.Config.Credentials == nil {
		t.Errorf("expect not nil")
	}

	configCreds, err := s.Config.Credentials.Get()

	if err != nil {
		t.Errorf("expect nil, got error: %s", err)
	}

	if configCreds.ProviderName != "EnvConfigCredentials" {
		t.Errorf("expect StaticProvider, got %s", configCreds.ProviderName)
	}
	if configCreds.AccessKeyID != staticAccessKey {
		t.Errorf("expect staticAccessKey")
	}
	if configCreds.SecretAccessKey != staticSecretKey {
		t.Errorf("expect staticSecretKey")
	}
	if configCreds.SessionToken != "" {
		t.Errorf("expect empty string")
	}

}

//if awsErr, ok := err.(awserr.Error); ok  {
//	if awsErr.Code() != "NoCredentialProviders" {
//		t.Errorf("expect error NoCredentialProviders, got %+v", awsErr.Code())
//	}
//}
