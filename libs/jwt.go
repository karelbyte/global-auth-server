package libs

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/golang-jwt/jwt/v5"
)

type SecretData struct {
	PrivatePem string `json:"private.pem"`
	PublicPem  string `json:"public.pem"`
}

var privateKey *rsa.PrivateKey

func formatPEMKey(key string) string {
	key = strings.TrimSpace(key)
	key = strings.ReplaceAll(key, "-----BEGIN PRIVATE KEY-----", "")
	key = strings.ReplaceAll(key, "-----END PRIVATE KEY-----", "")
	key = strings.TrimSpace(key)

	var formattedKey strings.Builder
	formattedKey.WriteString("-----BEGIN PRIVATE KEY-----\n")

	for i := 0; i < len(key); i += 64 {
		end := i + 64
		if end > len(key) {
			end = len(key)
		}
		formattedKey.WriteString(key[i:end])
		formattedKey.WriteString("\n")
	}

	formattedKey.WriteString("-----END PRIVATE KEY-----")
	return formattedKey.String()
}

func LoadPrivateKey() (*rsa.PrivateKey, error) {
	if privateKey != nil {
		return privateKey, nil
	}

	secretARN := os.Getenv("AWS_SECRET_ARN")
	region := os.Getenv("AWS_REGION")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if secretARN == "" || region == "" || accessKeyID == "" || secretAccessKey == "" {
		return nil, fmt.Errorf("missing required AWS environment variables")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("error loading AWS config: %w", err)
	}

	svc := secretsmanager.NewFromConfig(cfg)
	result, err := svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretARN),
		VersionStage: aws.String("AWSCURRENT"),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting secret value: %w", err)
	}

	var secretData SecretData
	if err := json.Unmarshal([]byte(*result.SecretString), &secretData); err != nil {
		return nil, fmt.Errorf("error parsing secret JSON: %w", err)
	}

	if secretData.PrivatePem == "" {
		return nil, fmt.Errorf("private key not found in secret")
	}

	formattedKey := formatPEMKey(secretData.PrivatePem)
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(formattedKey))
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}

	privateKey = key
	return privateKey, nil
}


func GenerateJWT(payload map[string]any, duration time.Duration) (string, int64, error) {
	key, err := LoadPrivateKey()
	if err != nil {
		return "", 0, err
	}
	exp := time.Now().Add(duration).Unix()
	claims := jwt.MapClaims{}
	maps.Copy(claims, payload)
	claims["exp"] = exp

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(key)
	if err != nil {
		return "", 0, err
	}
	return signed, exp, nil
}
