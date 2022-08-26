package lib

import (
	"net/url"
	"nft-site/db"
	"reflect"
	"strings"
)

func CheckDiff(u db.User, form url.Values) map[string]string {
	diffMap := make(map[string]string)
	if u.Consumer.Id != 0 {
		setDiffMap(diffMap, form, u.Consumer)
	} else {
		setDiffMap(diffMap, form, u.Seller)
	}
	return diffMap
}

func setDiffMap(diffMap map[string]string, form url.Values, user interface{}) map[string]string {
	rtUser := reflect.TypeOf(user)
	rvUser := reflect.ValueOf(user)
	for i := 1; i < rtUser.NumField()-1; i++ {
		name := rtUser.Field(i).Name
		key := titleLower(name)
		value := rvUser.FieldByName(name).Interface()

		if value != form[key][0] {
			column := toColumn(key)
			diffMap[column] = form[key][0]
		}
	}
	return diffMap
}

func titleLower(str string) string {
	s := string(str[0])
	l := strings.ToLower(s)
	return l + str[1:]
}

func toColumn(key string) string {
	columnMap := map[string]string{
		"familyName":   "family_name",
		"firstName":    "first_name",
		"nickname":     "nickname",
		"mail":         "mail",
		"company":      "company",
		"table":        "table",
		"lotteryUnits": "lottery_units",
	}
	return columnMap[key]
}
