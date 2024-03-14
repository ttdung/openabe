// package main

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/thoas/go-funk"
// )

// func replaceASCIIChars(input string) string {
// 	output := ""
// 	for index := 0 ; index < len(input) ; index++ {
// 		cur := index
// 		for (input[cur] >= 'a' && input[cur] <= 'z') || (input[cur] >= 'A' && input[cur] <= 'Z') {
// 			cur++
// 		}
// 		if cur == index {
// 			output += string(input[index])
// 		} else {
// 			output += "0"
// 			index = cur - 1
// 		}
	
// 	}
// 	return output
// }

// func reformatString(A string, B string) string {
// 	A = strings.ReplaceAll(A, "(", "4")
// 	A = strings.ReplaceAll(A, ")", "5")
// 	A = strings.ReplaceAll(A, "and", "2")
// 	A = strings.ReplaceAll(A, "or", "3")
// 	A = strings.ReplaceAll(A, " ", "")

// 	splittedB := strings.Split(B, "|") 
// 	funk.ForEach(splittedB, func(value string) {
// 		A = strings.ReplaceAll(A, value, "1")
// 	})

	
// 	return replaceASCIIChars(A)
// }

// func main() {
// 	A := "(student) and (CS)"
// 	B := "student|EE"
	
// 	formattedA := reformatString(A, B)

// 	stack := make([]string, 0)

// 	fmt.Println(formattedA)
// 	for i := 0 ; i < len(formattedA) ; i++ {
// 		if formattedA[i] == '5' {
// 			currentChain := ""
// 			for len(stack) > 0 && stack[len(stack) - 1] != "4" {
// 				lenS := len(stack)
// 				currentChain = stack[lenS - 1] + currentChain 
// 				stack = stack[:lenS - 1]
// 			}
// 			if len(currentChain) > 0 {
// 				stack = stack[:len(stack) - 1] // remove ")"
// 				cur := 1
// 				var op int
// 				for _, c := range currentChain {
// 					if c - '0' >= 1 {
// 						op = int(c - '0')
// 					}else {
// 						switch op {
// 						case 3: 
// 							cur = cur|int(c - '0')
// 							break 
// 						case 2:
// 							cur = cur&int(c - '0')
// 							break
// 						default:
// 							cur = int(c - '0')
// 						}
// 					}
// 				}
// 				stack = append(stack, string(cur + '0'))
// 				fmt.Println(cur)
// 			}
// 		}else {
// 			stack = append(stack, string(formattedA[i]))
// 		}
// 	}

// 	fmt.Println(stack)
// 	if len(stack) > 0 {
// 		cur := 1
// 		var op int
// 		for _, c := range stack {
// 			if c[0] - '0' >= 1 {
// 				op = int(c[0] - '0')
// 			}else {
// 				switch op {
// 				case 3: 
// 					cur = cur|int(c[0] - '0')
// 					break 
// 				case 2:
// 					cur = cur&int(c[0] - '0')
// 					break
// 				}
// 			}
// 		}
// 		fmt.Println(cur)
// 	}
// }
