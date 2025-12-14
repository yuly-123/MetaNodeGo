package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	//a()
	//b(12345321)
	//c("(]")
	d([]string{"as中", "asdf", "as123"})
	//d([]string{"asdf1", "asdf1", "asdf1"})
}

func a() {
	var numArry [7]int = [7]int{3, 6, 9, 5, 3, 6, 9}
	fmt.Println(numArry)

	var numMap map[int]int = make(map[int]int, 4)
	fmt.Println(numMap[0])

	// 把数字存入map，key是数字，value是出现次数
	for _, v := range numArry {
		numMap[v] += 1 // 从0开始，key值每覆盖一次，值加1
	}
	fmt.Println(numMap)

	for k, v := range numMap {
		if v == 1 { // 出现一次的数字，value值等于1
			fmt.Println(k)
			break
		}
	}
}

func b(num uint64) {
	// 最左和最右依次对比是否相等
	var s = strconv.FormatUint(num, 10)
	var l = 0          // 最左下标
	var r = len(s) - 1 // 最右下标
	var c = len(s) / 2 // 对比次数
	fmt.Println(s, l, r, c)

	var flag bool = true
	for i := 0; i < c; i++ {
		fmt.Println(l, r, s[l], s[r])
		if s[l] != s[r] {
			flag = false
		}
		l++
		r--
	}
	fmt.Println(flag)
}

func c(s string) {
	// 把给定的字符串，拆分单个字符，判断字符在规则中是否存在
	var regx string = "(){}[]"
	fmt.Println(regx)

	var flag bool = true
	for _, v := range s {
		str := string(v)
		fmt.Println(str)
		if !strings.Contains(regx, str) {
			flag = false
			break
		}
	}
	fmt.Println(flag)
}

func d(s []string) {
	// 先判断完全相等，
	var b bool = true
	for i := 0; i < len(s)-1; i++ {
		if s[i] != s[i+1] {
			b = false
		}
	}
	if b {
		fmt.Println("完全相等：", s[0])
		return
	}

	// 转换集合，
	var m map[int][]rune = make(map[int][]rune, len(s))
	for i, v := range s {
		m[i] = []rune(v)
	}
	fmt.Println(m)

	// 循环判断，
	var same map[int]rune = make(map[int]rune, 0) // 存储公共前缀，
	j := 0
	for {
		// 依次取出所有字符串的相同位置的字符
		var r map[int]rune = make(map[int]rune, 0)
		for i := 0; i < len(m); i++ {
			if j < len(m[i]) {
				r[i] = m[i][j]
			}
		}
		fmt.Println(r)

		// r中存放的字符数量，和m的元素个数不相同，表示某些字符串已经取完了，
		if len(r) != len(m) {
			break
		}

		// 判断是否相同
		var result rune = 0
		for _, v := range r {
			result += v
		}
		// 每个值加起来除以数量等于这个值，表示所有值完全一样
		if result/rune(len(r)) == r[0] {
			same[j] = r[0]
		} else {
			break // 已经不相同了，结束对比，后面的字符不用在对比，
		}
		j++ // 对比下一位的字符，
	}

	fmt.Println(same) // 存储公共前缀，

	// 转换成字符串输出，
	var sameRune []rune = make([]rune, len(same))
	for i := 0; i < len(same); i++ {
		sameRune[i] = same[i]
	}
	fmt.Println("公共前缀：", string(sameRune))
}
