package util

import (
	"math"
	"net"
	"os"
)

const TimeFormat string = "2006-01-02-T15:04:05-0700"

/*
 * Name:     IsIP
 * Purpose:  Returns true if string is a valid IP address, false otherwise
 * comments:
 */
func IsIP(ip string) bool {
	if net.ParseIP(ip) != nil {
		return true
	}
	return false
}

/*
 * Name:     Exists
 * Purpose:  Returns true if file or directory exists, false otherwise
 * comments:
 */
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// ByStringLength Functions that, in combination with golang sort,
// allow users to sort a slice/list of strings by string length
// (shortest -> longest)
type ByStringLength []string

func (s ByStringLength) Len() int           { return len(s) }
func (s ByStringLength) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByStringLength) Less(i, j int) bool { return len(s[i]) < len(s[j]) }

// SortableInt64 functions that allow a golang sort of int64s
type SortableInt64 []int64

func (s SortableInt64) Len() int           { return len(s) }
func (s SortableInt64) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortableInt64) Less(i, j int) bool { return s[i] < s[j] }

//given a sorted slice, remove the duplicates
func RemoveSortedDuplicates(sortedIn []int64) []int64 {
	//Avoid some reallocations
	result := make([]int64, 0, len(sortedIn)/2)
	last := sortedIn[0]
	result = append(result, last)

	for idx := 1; idx < len(sortedIn); idx++ {
		if last != sortedIn[idx] {
			result = append(result, sortedIn[idx])
		}
		last = sortedIn[idx]
	}
	return result
}

func CountAndRemoveSortedDuplicates(sortedIn []int64) ([]int64, map[int64]int64) {
	//Avoid some reallocations
	result := make([]int64, 0, len(sortedIn)/2)
	counts := make(map[int64]int64)

	last := sortedIn[0]
	result = append(result, last)
	counts[last]++

	for idx := 1; idx < len(sortedIn); idx++ {
		if last != sortedIn[idx] {
			result = append(result, sortedIn[idx])
		}
		last = sortedIn[idx]
		counts[last]++
	}
	return result, counts
}

//two's complement 64 bit abs value
func Abs(a int64) int64 {
	mask := a >> 63
	a = a ^ mask
	return a - mask
}

//rounding function since go doesn't have it
func Round(f float64) int64 {
	return int64(math.Floor(f + .5))
}

//retun the smaller of two integers
func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
