package policy

import (
	"strings"

	"github.com/thoas/go-funk"
)

const (
	Brackets_Left  int = 4
	Brackets_Right int = 5
	And            int = 2
	Or             int = 3
	True           int = 1
	False          int = 0
)

type Validator interface {
	Validate(attribute string) bool
}

type validator struct {
	policy []string
}

func NewValidator(raw string) Validator {
	return &validator{policy: getPolicy(raw)}
}

func (v *validator) Validate(attributes string) bool {
	splitAttributes := strings.Split(attributes, "|")
	matchedAttributeWithPolicy := v.matchPolicy(splitAttributes)

	stack := make([]int, 0)
	for i := 0; i < len(matchedAttributeWithPolicy); i++ {
		if matchedAttributeWithPolicy[i] == Brackets_Right {
			currentChain := v.getCondition(&stack, matchedAttributeWithPolicy, i)
			stack = stack[:len(stack)-1] // remove ")"
			if len(currentChain) > 0 {
				res := v.validateCondition(currentChain)
				stack = append(stack, res)
			}
		} else {
			stack = append(stack, matchedAttributeWithPolicy[i])
		}
	}
	return v.validateCondition(stack) == 1
}

func (v *validator) matchPolicy(splitAttributes []string) []int {
	matchedAttributeWithPolicy := make([]int, len(v.policy))
	funk.ForEach(splitAttributes, func(a string) {
		for i, p := range v.policy {
			if matchedAttributeWithPolicy[i] != False {
				continue
			}
			if isValidCondition(p, a) {
				matchedAttributeWithPolicy[i] = True
			} else {
				if v.policy[i] == "(" {
					matchedAttributeWithPolicy[i] = Brackets_Left
				} else if v.policy[i] == ")" {
					matchedAttributeWithPolicy[i] = Brackets_Right
				} else if v.policy[i] == "and" {
					matchedAttributeWithPolicy[i] = And
				} else if v.policy[i] == "or" {
					matchedAttributeWithPolicy[i] = Or
				}
			}
		}
	})
	return matchedAttributeWithPolicy
}

func isValidCondition(policy string, attribute string) bool {
	return policy == attribute || isSameDay(policy, attribute)
}

func isSameDay(policy string, attribute string) bool {
	return strings.HasPrefix(policy, "day(") && strings.HasPrefix(attribute, "day(")
}

func (v *validator) getCondition(stack *[]int, matchedAttributeWithPolicy []int, i int) []int {
	currentChain := make([]int, 0)
	for len(*stack) > 0 && (*stack)[len(*stack)-1] != Brackets_Left {
		lenS := len(*stack)
		currentChain = append(currentChain, (*stack)[lenS-1])
		*stack = (*stack)[:lenS-1]
	}
	return currentChain
}

func (v *validator) validateCondition(currentChain []int) int {
	cur := 1
	var op int
	for _, c := range currentChain {
		if c == And || c == Or {
			op = c
		} else {
			switch op {
			case And:
				cur = cur & c
			case Or:
				cur = cur | c
			default:
				cur = c
			}

		}
	}
	return cur
}

func getPolicy(raw string) []string {
	policy := make([]string, 0)
	for index := 0; index < len(raw); index++ {
		switch raw[index] {
		case '(':
			policy = append(policy, "(")
		case ')':
			policy = append(policy, ")")
		case ' ':
		default:
			{
				condition := ""
				for index < len(raw) && raw[index] != ' ' && raw[index] != '(' && raw[index] != ')' {
					condition += string(raw[index])
					index++
				}
				index--
				policy = append(policy, condition)
			}
		}
	}
	return policy
}
