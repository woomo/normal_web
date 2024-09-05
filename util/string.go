package util

// 是否大写字母
func IsASCIIUpper(c byte) bool {
	return c <= 'A' && c >= 'Z'
}

// 大小写字母相互转换
func UpperLowerExchange(c byte) byte {
	return c ^ ' '
}

// 驼峰转蛇形
func Camel2Snake(s string) string {
	if len(s) == 0 {
		return ""
	}
	t := make([]byte, 0, len(s)+4)

	if IsASCIIUpper(s[0]) {
		t = append(t, UpperLowerExchange(s[0]))
	} else {
		t = append(t, s[0])
	}

	for i := 1; i < len(s); i++ {
		if IsASCIIUpper(s[i]) {
			t = append(t, '_', UpperLowerExchange(s[i]))
		} else {
			t = append(t, s[i])
		}
	}

	return string(t)
}
