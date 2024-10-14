package util

import (
	"slack-wails/lib/structs"
	"sort"
	"strings"
)

// 去除某个元素
func RemoveElement[T comparable](arr []T, ele T) (newArr []T) {
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
	if len(arr) == 0 {
		return result
	}
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
func ArrayContains[T comparable](t T, array []T) bool {
	if len(array) == 0 {
		return false
	}
	for _, ele := range array {
		if t == ele {
			return true
		}
	}
	return false
}

// IntArrayToUint16Array 将 int 数组转换为 uint16 数组
func IntArrayToUint16Array(intArray []int) []uint16 {
	uint16Array := make([]uint16, len(intArray))
	for i, v := range intArray {
		uint16Array[i] = uint16(v)
	}
	return uint16Array
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

// 拼接非空值项数组
func MergeNonEmpty(arrayList []string, connector string) string {
	var newArr []string
	for _, v := range arrayList {
		if v != "" {
			newArr = append(newArr, v)
		}
	}
	return strings.Join(newArr, connector)
}

// 处理空间引擎的位置数据，删除非空空项，并处理直辖市只显示一次，按照连接符返回字符串
func MergePosition(position structs.Position) string {
	fields := []string{}
	if position.Country != "" {
		fields = append(fields, position.Country)
	}
	if position.Province != "" && position.Province != position.City {
		fields = append(fields, position.Province)
	}
	if position.City != "" {
		fields = append(fields, position.City)
	}
	if position.District != "" {
		fields = append(fields, position.District)
	}
	return strings.Join(fields, position.Connector)
}
