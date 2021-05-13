package signature_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	pkg_signature "tiniapp-backend-oauth-sample/pkg/signature"
)

func TestSignature_PrepareData_OK(t *testing.T) {
	const float64Value = float64(9007199254740991)
	const stringValue = "RLCKb7Ae9kx4DXtXsCWjnDXtggFnM43W"
	const int64Value = int64(1620901354516)

	data := map[string]interface{}{
		"float64_field": float64Value,
		"string_field":  stringValue,
		"int64_field":   int64Value,
		"map_field": map[string]interface{}{
			"float64_field": float64Value,
			"string_field":  stringValue,
			"int64_field":   int64Value,
		},
		"arr_field": []map[string]interface{}{
			{
				"float64_field": float64Value,
				"string_field":  stringValue,
				"int64_field":   int64Value,
			},
			{
				"float64_field": float64Value,
				"string_field":  stringValue,
				"int64_field":   int64Value,
			},
		},
	}

	payload, err := pkg_signature.PreparePayload(data)
	require.NoError(t, err)

	jsonData, err := base64.RawURLEncoding.DecodeString(payload)
	require.NoError(t, err)

	type Data struct {
		Float64Field float64 `json:"float64_field"`
		StringField  string  `json:"string_field"`
		Int64Field   int64   `json:"int64_field"`
		MapField     *struct {
			Float64Field float64 `json:"float64_field"`
			StringField  string  `json:"string_field"`
			Int64Field   int64   `json:"int64_field"`
		} `json:"map_field"`
		ArrField []*struct {
			Float64Field float64 `json:"float64_field"`
			StringField  string  `json:"string_field"`
			Int64Field   int64   `json:"int64_field"`
		} `json:"arr_field"`
	}

	var decodedData Data
	err = json.Unmarshal(jsonData, &decodedData)
	require.NoError(t, err)

	require.Equal(t, float64Value, decodedData.Float64Field)
	require.Equal(t, stringValue, decodedData.StringField)
	require.Equal(t, int64Value, decodedData.Int64Field)

	require.Equal(t, float64Value, decodedData.MapField.Float64Field)
	require.Equal(t, stringValue, decodedData.MapField.StringField)
	require.Equal(t, int64Value, decodedData.MapField.Int64Field)

	for _, item := range decodedData.ArrField {
		require.Equal(t, float64Value, item.Float64Field)
		require.Equal(t, stringValue, item.StringField)
		require.Equal(t, int64Value, item.Int64Field)
	}
}

func TestSignature_Sign_OK(t *testing.T) {
	sign, err := pkg_signature.SignWithData(
		"EhjGcsUUuRSJTHiYPbW5fxzyaKEx0JuAZIKRQ4HnIfNFidB2kMg6locQbTIEz3Vf",
		map[string]interface{}{
			"code":      "Y24CwtKFz7aPXv8wmghwZpnqHhpHaeA0THJ0qIY3BXjwm744zZ0JY6SkhRiNFVJT",
			"client_id": "RLCKb7Ae9kx4DXtXsCWjnDXtggFnM43W",
			"timestamp": int64(1620901354516),
		},
	)
	require.NoError(t, err)
	require.Equal(t, "f94fbc88a110c565b85904e963f5b82c7e9c1423480cb6dd0b0387d7d956533d", sign)
}
