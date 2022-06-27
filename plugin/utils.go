package gaccauth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"golang.org/x/oauth2"
)

type GenericMap map[string]interface{}

func strSliceEquals(a, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)

	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func strSliceHasIntersection(a, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)

	for i, j := 0, 0; i < len(a) && j < len(b); {
		if a[i] == b[j] {
			return true
		}
		if a[i] < b[j] {
			i++
		} else {
			j++
		}
	}

	return false
}

func encodeToken(token *oauth2.Token) (string, error) {
	buf, err := json.Marshal(token)

	if err != nil {
		return "", err
	}

	return string(buf), err
}

func decodeToken(encoded string) (*oauth2.Token, error) {
	var token oauth2.Token

	if err := json.Unmarshal([]byte(encoded), &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func getPositiveIntData(data *framework.FieldData, prop string) (*int, error) {
	if v, ok := data.GetOk(prop); ok {
		value := v.(int)
		if value < 1 {
			return nil, fmt.Errorf("value cannot be negative or zero")
		}

		return &value, nil
	}

	return nil, nil
}

func getRequiredStringData(data *framework.FieldData, prop string) (*string, error) {
	if v, ok := data.GetOk(prop); ok {
		value := strings.TrimSpace(v.(string))
		if len(value) > 0 {
			return &value, nil
		}

		return nil, fmt.Errorf("property '%s' cannot be empty", prop)
	}

	return nil, fmt.Errorf("missing property '%s' in configuration", prop)
}

func getFilteredStringSliceData(data *framework.FieldData, prop string) *[]string {
	if v, ok := data.GetOk(prop); ok {
		filteredValues := []string{}

		for _, value := range v.([]string) {
			tv := strings.TrimSpace(value)
			if len(tv) > 0 {
				filteredValues = append(filteredValues, strings.TrimSpace(value))
			}
		}

		return &filteredValues
	}

	return nil
}

func isValidUrl(addr string) bool {
	u, err := url.Parse(addr)

	if err != nil {
		return false
	}

	protocol := strings.ToLower(u.Scheme)
	if !(protocol == "http" || protocol == "https") {
		return false
	}

	return true
}
