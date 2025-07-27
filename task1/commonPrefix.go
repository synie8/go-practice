package task1

import "fmt"

func LongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	reslut := ""
	for i := 0; i < len(strs[0]); i++ {
		for _, v := range strs {
			if i >= len(v) || v[i] != strs[0][i] {
				return reslut
			}
		}
		reslut += string(strs[0][i])
		fmt.Println("result ***** ", reslut)
	}

	return reslut
}
