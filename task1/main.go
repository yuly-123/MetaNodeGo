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
	//d([]string{"as中", "asdf", "as123"})
	//d([]string{"asdf1", "asdf1", "asdf1"})
	//e([4]uint64{9, 9, 9, 9})
	//e([4]uint64{4, 3, 2, 1})
	//f([10]int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4})
	//g([5][2]int{{1, 3}, {15, 18}, {4, 7}, {10, 12}, {2, 6}})
	//g([4][2]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}})
	//g([2][2]int{{4, 7}, {1, 4}})
	//h([4]int{2, 7, 11, 15}, 9)
}

// 只出现一次的数字：给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
// 给你一个 非空 整数数组 nums ，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
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

// 回文数：判断一个整数是否是回文数
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

// 有效的括号：给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效
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

// 最长公共前缀：查找字符串数组中的最长公共前缀
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

// 加一：给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
func e(digits [4]uint64) {
	// 数组转字符串
	var s string = ""
	for _, v := range digits {
		s += "" + strconv.FormatUint(v, 10)
	}
	fmt.Println(s)

	// 字符串转数字+1后再转回字符串
	i, _ := strconv.ParseUint(s, 10, 64)
	i++
	s = strconv.FormatUint(i, 10)
	fmt.Println(s)

	// 字符串转[]rune
	var r []uint64 = make([]uint64, len(s))
	for i := 0; i < len(s); i++ {
		v, _ := strconv.ParseUint(string(s[i]), 10, 64)
		r[i] = v

	}
	fmt.Println(r)
}

// 删除有序数组中的重复项：给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
// 不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。
// 给你一个 非严格递增排列 的数组 nums ，请你 原地 删除重复出现的元素，使每个元素 只出现一次 ，返回删除后数组的新长度。
// 元素的 相对顺序 应该保持 一致 。然后返回 nums 中唯一元素的个数。
func f(nums [10]int) {
	fmt.Println(nums)

	var c int = 0       // 计数
	var n int = nums[0] // 对比数
	for i := 1; i < len(nums); i++ {
		if n == nums[i] {
			nums[i] = 9999
			c++
		} else {
			n = nums[i]
		}
	}
	fmt.Println(nums)

	// 排序，从小到大，
	for i := 0; i < len(nums)-1; i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] > nums[j] {
				n = nums[i]
				nums[i] = nums[j]
				nums[j] = n
			}
		}
		//fmt.Println(nums)
	}
	fmt.Println(nums)
	fmt.Println(len(nums) - c)
}

// 合并区间：以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。可以先对区间数组按照区间的起始位置进行排序。
// 示例 1：
// 输入：intervals = [[1,3],[2,6],[8,10],[15,18]]
// 输出：[[1,6],[8,10],[15,18]]
// 解释：区间 [1,3] 和 [2,6] 重叠, 将它们合并为 [1,6].
// 示例 2：
// 输入：intervals = [[1,4],[4,5]]
// 输出：[[1,5]]
// 解释：区间 [1,4] 和 [4,5] 可被视为重叠区间。
// 示例 3：
// 输入：intervals = [[4,7],[1,4]]
// 输出：[[1,7]]
// 解释：区间 [1,4] 和 [4,7] 可被视为重叠区间。
func g(nums [2][2]int) {
	fmt.Println("原始输出：", nums)

	// 排序，从小到大，
	var temp int
	for i := 0; i < len(nums)-1; i++ {
		for j := i + 1; j < len(nums); j++ {
			// 第1位数相等的，第2位数排序，从小到大
			if nums[i][0] == nums[j][0] && nums[i][1] > nums[j][1] {
				temp = nums[i][1]
				nums[i][1] = nums[j][1]
				nums[j][1] = temp
			} else if nums[i][0] > nums[j][0] {
				temp = nums[i][0]
				nums[i][0] = nums[j][0]
				nums[j][0] = temp

				temp = nums[i][1]
				nums[i][1] = nums[j][1]
				nums[j][1] = temp
			}
		}
	}
	fmt.Println("顺序输出：", nums)

	var c int = 0
	var m map[int][2]int = make(map[int][2]int)
	for i := 0; i < len(nums)-1; i++ {
		// if 无重合，直接进入下一个元素，再往后对比
		// else 有重合，当前元素需要继续和下下个元素对比，
		if nums[i][1] < nums[i+1][0] {
			m[c] = [2]int{nums[i][0], nums[i][1]} // 收集
			c++
			m[c] = [2]int{nums[i+1][0], nums[i+1][1]} // 把后一个暂时也收集进来
		} else {
			nums[i+1][0] = nums[i][0] // 改变下一个元素的值，并收集
			m[c] = nums[i+1]          // 收集
		}
	}

	fmt.Println("合并输出：", m)
}

// 两数之和：给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数，并返回它们的数组下标。
// 你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。
// 你可以按任意顺序返回答案。
// 示例 1：
// 输入：nums = [2,7,11,15], target = 9
// 输出：[0,1]
// 解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1] 。
// 示例 2：
// 输入：nums = [3,2,4], target = 6
// 输出：[1,2]
// 示例 3：
// 输入：nums = [3,3], target = 6
// 输出：[0,1]
func h(nums [4]int, target int) {
	fmt.Println(nums, target)

a:
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				fmt.Println(i, j)
				break a
			}
		}
	}
}
