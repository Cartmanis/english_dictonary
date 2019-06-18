package crypto

import (
	"fmt"
	"testing"
)

func TestRandPassword(t *testing.T) {
	test := RandPassword(3)
	fmt.Println(test)

	test1 := RandPassword(6)
	fmt.Println(test1)

	test2 := RandPassword(7)
	fmt.Println(test2)

	test3 := RandPassword(8)
	fmt.Println(test3)

	test4 := RandPassword(9)
	fmt.Println(test4)

	test5 := RandPassword(10)
	fmt.Println(test5)

	test6 := RandPassword(11)
	fmt.Println(test6)

	test7 := RandPassword(12)
	fmt.Println(test7)

	test8 := RandPassword(13)
	fmt.Println(test8)

	test9 := RandPassword(14)
	fmt.Println(test9)

	test10 := RandPassword(15)
	fmt.Println(test10)

	test11 := RandPassword(16)
	fmt.Println(test11)

	test12 := RandPassword(18)
	fmt.Println(test12)
}
