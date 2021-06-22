package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
)

func PreparePayload(timestamp int64, clientID string, body interface{}) (string, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	payload := fmt.Sprintf("%s.%s.%s", strconv.FormatInt(timestamp, 10), clientID, string(jsonBody))
	encodedPayload := base64.RawURLEncoding.EncodeToString([]byte(payload))
	return encodedPayload, nil
}

func SignWithPayload(secret string, payload string) (string, error) {
	h := hmac.New(sha256.New, []byte(secret))

	_, err := h.Write([]byte(payload))
	if err != nil {
		return "", err
	}

	signature := hex.EncodeToString(h.Sum(nil))
	return signature, nil
}
