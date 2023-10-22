package util

import "sort"

// 去除某个元素
func RemoveElement(arr []string, ele string) (newArr []string) {
	for _, v := range arr {
		if v != ele {
			newArr = append(newArr, v)
		}
	}
	return newArr
}

// 泛型数组去重
func RemoveDuplicates[T comparable](arr []T) []T {
	encountered := map[T]bool{}
	result := []T{}
	for _, v := range arr {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

// 替换指定元素
func ReplaceElement(slice []string, old, new string) []string {
	newSlice := []string{}
	for _, v := range slice {
		if v != old {
			newSlice = append(newSlice, v)
		} else {
			newSlice = append(newSlice, new)
		}
	}
	return newSlice
}

// 泛型判断数组中是否存在某个元素类似python in
// 如果数组长度为0那么也返回成功 -> dirsearch排除状态码
func ArrayContains[T comparable](t T, array []T) bool {
	if len(array) == 0 {
		return true
	}
	for _, ele := range array {
		if t == ele {
			return true
		}
	}
	return false
}

type Pair struct {
	Key   string
	Value int
}

// map按照value排序
func SortMap(temp map[string]int) []Pair {
	// 将map转换为切片
	var pairs []Pair
	for key, value := range temp {
		pairs = append(pairs, Pair{key, value})
	}
	func(pairs []Pair) {
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Value > pairs[j].Value
		})
	}(pairs)
	return pairs
}

func SplitInt(n, slice int) []int {
	var res []int
	for n > slice {
		res = append(res, slice)
		n = n - slice
	}
	res = append(res, n)
	return res
}

// 两个数组进行排列组合
func Combination(s1, s2 []string, split string) []string {
	var temp []string
	for _, v1 := range s1 {
		for _, v2 := range s2 {
			temp = append(temp, v1+split+v2)
		}
	}
	return temp
}
