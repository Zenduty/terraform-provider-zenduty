package zenduty

import (
	"encoding/json"
	"regexp"
)

func isJSONString(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// func isElementExist(s []int, str int) bool {
// 	for _, v := range s {
// 		if v == str {
// 			return true
// 		}
// 	}
// 	return false
// }

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
