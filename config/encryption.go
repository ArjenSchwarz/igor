package config

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

// DecryptString decrypts a string if KMS is enabled, otherwise it will return
// the string as is
func DecryptString(toDecrypt string) (string, error) {
	generalConfig, _ := GeneralConfig()
	return decryptValue(generalConfig.Kms, toDecrypt)
}

func decryptValue(useKms bool, toDecrypt string) (string, error) {
	if !useKms {
		return toDecrypt, nil
	}
	sess, err := session.NewSession()
	if err != nil {
		return toDecrypt, err
	}
	svc := kms.New(sess)

	decodedString, err := base64.StdEncoding.DecodeString(toDecrypt)
	if err != nil {
		return toDecrypt, err
	}

	params := &kms.DecryptInput{
		CiphertextBlob: []byte(decodedString), // Required
	}
	resp, err := svc.Decrypt(params)

	if err != nil {
		return toDecrypt, err
	}

	return string(resp.Plaintext[:]), nil
}
