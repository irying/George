package main

func findAnagrams(s string, p string) {
	var res = []int{}
	var target, window [26]int
	for i:=0; i < len(p); i++ {
		target[p[i] - 'a']++
	}

	i := 0
	r := -1
	len := len(s)
	return res

}

func main() {
}
