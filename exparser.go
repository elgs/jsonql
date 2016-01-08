package jsonql

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type Operator struct {
	Precedence int
	Eval       func(symbolTable interface{}, left string, right string) (string, error)
}

type Parser struct {
	Operators   map[string]*Operator
	SymbolTable interface{}
	maxOpLen    int
	initialized bool
}

func (this *Parser) Init() {
	for k, _ := range this.Operators {
		if len(k) > this.maxOpLen {
			this.maxOpLen = len(k)
		}
	}
}

func (this *Parser) Calculate(expression string) (string, error) {
	tokens := this.Tokenize(expression)
	//fmt.Println(expression, tokens)
	rpn, err := this.ParseRPN(tokens)
	if err != nil {
		return "", err
	}
	return this.Evaluate(rpn, true)
}

func (this *Parser) Evaluate(ts *Lifo, postfix bool) (string, error) {
	newTs := &Lifo{}
	usefulWork := false
	for ti := ts.Pop(); ti != nil; ti = ts.Pop() {
		t := ti.(string)
		//		fmt.Println("t:", t)
		switch {
		case this.Operators[t] != nil:
			// operators
			usefulWork = true
			if postfix {
				right := newTs.Pop()
				left := newTs.Pop()
				l := "0"
				r := "0"
				if left != nil {
					l = left.(string)
				}
				if right != nil {
					r = right.(string)
				}
				result, err := this.Operators[t].Eval(this.SymbolTable, l, r)
				newTs.Push(result)
				if err != nil {
					return "", errors.New(fmt.Sprint("Failed to evaluate:", l, t, r))
				}
			} else {
				right := ts.Pop()
				left := ts.Pop()
				l := ""
				r := ""
				if left != nil {
					l = left.(string)
				}
				if right != nil {
					r = right.(string)
				}
				result, err := this.Operators[t].Eval(this.SymbolTable, l, r)
				newTs.Push(result)
				if err != nil {
					return "", errors.New(fmt.Sprint("Failed to evaluate:", l, t, r))
				}
			}
		default:
			// operands
			newTs.Push(t)
		}
		//newTs.Print()
	}
	if !usefulWork {
		return "", errors.New("Failed to evaluate: no valid operator found.")
	}
	if newTs.Len() == 1 {
		return newTs.Pop().(string), nil
	} else {
		return this.Evaluate(newTs, !postfix)
	}
}

// false o1 in first, true o2 out first
func (this *Parser) shunt(o1, o2 string) (bool, error) {
	op1 := this.Operators[o1]
	op2 := this.Operators[o2]
	if op1 == nil || op2 == nil {
		return false, errors.New(fmt.Sprint("Invalid operators:", o1, o2))
	}
	if op1.Precedence < op2.Precedence || (op1.Precedence <= op2.Precedence && op1.Precedence%2 == 1) {
		return true, nil
	}
	return false, nil
}

func (this *Parser) ParseRPN(tokens []string) (output *Lifo, err error) {
	opStack := &Lifo{}
	outputQueue := []string{}
	for _, token := range tokens {
		switch {
		case this.Operators[token] != nil:
			// operator
			for o2 := opStack.Peep(); o2 != nil; o2 = opStack.Peep() {
				stackToken := o2.(string)
				if this.Operators[stackToken] == nil {
					break
				}
				o2First, err := this.shunt(token, stackToken)
				if err != nil {
					return output, err
				}
				if o2First {
					outputQueue = append(outputQueue, opStack.Pop().(string))
				} else {
					break
				}
			}
			opStack.Push(token)
		case token == "(":
			opStack.Push(token)
		case token == ")":
			for o2 := opStack.Pop(); o2 != nil && o2.(string) != "("; o2 = opStack.Pop() {
				outputQueue = append(outputQueue, o2.(string))
			}
		default:
			outputQueue = append(outputQueue, token)
		}
	}
	for o2 := opStack.Pop(); o2 != nil; o2 = opStack.Pop() {
		outputQueue = append(outputQueue, o2.(string))
	}
	//fmt.Println(outputQueue)
	output = &Lifo{}
	for i := 0; i < len(outputQueue); i++ {
		(*output).Push(outputQueue[len(outputQueue)-i-1])
	}
	return
}

func normalize(exp string) string {
	re := regexp.MustCompile("\\s+OR")
	fmt.Println(re.ReplaceAllLiteralString(exp, "T"))
	return exp
}

func (this *Parser) Tokenize(exp string) (tokens []string) {
	if !this.initialized {
		this.Init()
	}
	exp = normalize(exp)
	sq, dq := false, false
	var tmp string
	expRunes := []rune(exp)
	for i := 0; i < len(expRunes); i++ {
		v := expRunes[i]
		s := string(v)
		switch {
		case unicode.IsSpace(v):
			//			if sq || dq {
			//				tmp += s
			//			} else {
			//				if len(tmp) > 0 && !sq && !dq {
			//					tokens = append(tokens, tmp)
			//					tmp = ""
			//				}
			//			}

			if sq || dq {
				tmp += s
			} else if len(tmp) > 0 {
				tokens = append(tokens, tmp)
				tmp = ""
			}
		case s == "'":
			tmp += s
			if !dq {
				sq = !sq
				if !sq {
					tokens = append(tokens, tmp)
					tmp = ""
				}
			}
		case s == "\"":
			tmp += s
			if !sq {
				dq = !dq
				if !dq {
					tokens = append(tokens, tmp)
					tmp = ""
				}
			}
		case s == "+" || s == "-" || s == "(" || s == ")":
			if sq || dq {
				tmp += s
			} else {
				if len(tmp) > 0 {
					tokens = append(tokens, tmp)
					tmp = ""
				}
				lastToken := ""
				if len(tokens) > 0 {
					lastToken = tokens[len(tokens)-1]
				}
				if (s == "+" || s == "-") && (len(tokens) == 0 || lastToken == "(" || this.Operators[lastToken] != nil) {
					// sign
					tmp += s
				} else {
					// operator
					tokens = append(tokens, s)
				}
			}
		default:
			if sq || dq {
				tmp += s
			} else {
				// until the max length of operators(n), check if next 1..n runes are operator, greedily
				opCandidateTmp := ""
				opCandidate := ""
				for j := 0; j < this.maxOpLen && i < len(expRunes)-j-1; j++ {
					next := string(expRunes[i+j])
					opCandidateTmp += strings.ToUpper(next)
					if this.Operators[opCandidateTmp] != nil {
						opCandidate = opCandidateTmp
					}
				}
				if len(opCandidate) > 0 {
					if len(tmp) > 0 {
						tokens = append(tokens, tmp)
						tmp = ""
					}
					tokens = append(tokens, opCandidate)
					i += len(opCandidate) - 1
				} else {
					tmp += s
				}
			}
		}
	}
	if len(tmp) > 0 {
		tokens = append(tokens, tmp)
		tmp = ""
	}
	return
}
