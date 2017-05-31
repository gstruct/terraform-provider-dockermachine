package provider

func ss2is(s []string) []interface{} {
	ret := make([]interface{}, len(s))
	for i := range s {
		ret[i] = s[i]
	}
	return ret
}

func is2ss(s []interface{}) []string {
	ret := make([]string, len(s))
	for i := range s {
		ret[i] = s[i].(string)
	}
	return ret
}
