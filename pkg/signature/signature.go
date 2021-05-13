package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"sort"
)

func sortKeys(data map[string]interface{}) map[string]interface{} {
	keys := []string{}
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sortedData := map[string]interface{}{}
	for _, key := range keys {
		switch val := data[key].(type) {
		case map[string]interface{}:
			sortedData[key] = sortKeys(val)
		case []map[string]interface{}:
			arr := []map[string]interface{}{}
			for _, v := range val {
				arr = append(arr, sortKeys(v))
			}
			sortedData[key] = arr
		default:
			sortedData[key] = data[key]
		}
	}
	return sortedData
}

func PreparePayload(data map[string]interface{}) (string, error) {
	sortedData := sortKeys(data)

	buffer, err := json.Marshal(sortedData)
	if err != nil {
		return "", err
	}

	payload := base64.RawURLEncoding.EncodeToString(buffer)
	return payload, nil
}

func SignWithData(secret string, data map[string]interface{}) (string, error) {
	payload, err := PreparePayload(data)
	if err != nil {
		return "", err
	}

	return SignWithPayload(secret, payload)
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
