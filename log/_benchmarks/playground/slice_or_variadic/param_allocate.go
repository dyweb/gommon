package main

type st struct {
	x int
	s string
	i interface{}
}

func fVariadic(nums ...int) int {
	s := 0
	for i := 0; i < len(nums); i++ {
		s += nums[i]
	}
	return s
}

func fVariadicSt(nums ...st) int {
	s := 0
	for i := 0; i < len(nums); i++ {
		s += nums[i].x
	}
	return s
}

func fSlice(nums []int) int {
	s := 0
	for i := 0; i < len(nums); i++ {
		s += nums[i]
	}
	return s
}

func fSliceSt(nums []st) int {
	s := 0
	for i := 0; i < len(nums); i++ {
		s += nums[i].x
	}
	return s
}

// go build -gcflags "-m -m" .
func main() {
	fVariadic(1, 2, 3)
	fSlice([]int{1, 2, 3})
}
