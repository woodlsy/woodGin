package number

import "strconv"

//
// Int64ToInt
// @Description: int64转int
// @param num
// @return numInt
//
func Int64ToInt(num int64) (numInt int) {
	numString := strconv.FormatInt(num, 10)
	numInt, _ = strconv.Atoi(numString)
	return
}
