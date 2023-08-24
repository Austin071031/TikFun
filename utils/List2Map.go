package utils

func Tomap(list []int) map[int]int{
	data := make(map[int]int)
	for i, id := range list {
		data[id] = i
	}
	return data
}
