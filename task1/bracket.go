package task1

import "fmt"

/*
*
*示例 1：

输入：s = "()"

输出：true

示例 2：

输入：s = "()[]{}"

输出：true

示例 3：

输入：s = "(]"

输出：false

示例 4：

输入：s = "([])"

输出：true

示例 5：

输入：s = "([)]"

输出：false

1 <= s.length <= 10**4
s 仅由括号 '()[]{}' 组成
*/
func IsValid(s string) bool {

	if l := len(s); l%2 != 0 || l == 0 {
		return false
	}
	ss := make([]byte, 0)
	for i, v := range s {
		fmt.Printf("i= %d*****v= %c \n", i, v)
		if v == '(' || v == '[' || v == '{' {
			ss = append(ss, byte(v))
		} else {
			fmt.Println("ss=  \n", ss)
			if v == ')' || v == ']' || v == '}' {
				if len(ss) == 0 {
					return false
				}
				index := len(ss) - 1
				e := ss[index]

				if !((e == '(' && v == ')') || (e == '[' && v == ']') || (e == '{' && v == '}')) {
					return false
				}
				ss = ss[:index]
			}

		}

	}

	return len(ss) == 0

}

// 优化后
func IsValid1(s string) bool {
	if len(s)%2 != 0 {
		return false
	}

	pairs := map[rune]rune{
		'(': ')',
		'[': ']',
		'{': '}',
	}

	stack := make([]rune, 0, len(s)/2)

	for i, char := range s {
		fmt.Println("stack 前  ", stack)
		if closing, isOpen := pairs[char]; isOpen {
			fmt.Printf("i= %d*****v= %c \n", i, closing)
			stack = append(stack, closing)
		} else {
			if len(stack) == 0 || stack[len(stack)-1] != char {
				return false
			}
			stack = stack[:len(stack)-1]
		}
		fmt.Println("stack 后  ", stack)
	}

	return len(stack) == 0

}
