package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Calc принимает строку выражения и возвращает вычисленное значение или ошибку
func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")

	var stack []float64
	var opStack []rune
	precedence := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}
	var numberBuffer []rune

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		// tсли цифра или точка, собираем число
		if unicode.IsDigit(char) || char == '.' {
			numberBuffer = append(numberBuffer, char)
			continue
		}

		// если число накопилось в буфере, добавляем его в стек
		if len(numberBuffer) > 0 {
			num, err := strconv.ParseFloat(string(numberBuffer), 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
			numberBuffer = nil // Очищаем буфер
		}

		// если откр скобка
		if char == '(' {
			opStack = append(opStack, char)
			continue
		}

		// если закр скобка выполн операции до открывающей скобки
		if char == ')' {
			for len(opStack) > 0 && opStack[len(opStack)-1] != '(' {
				if len(stack) < 2 {
					return 0, errors.New("недостаточно операндов")
				}
				res, err := applyOperator(&stack, &opStack)
				if err != nil {
					return 0, err
				}
				stack = append(stack, res)
			}
			if len(opStack) == 0 {
				return 0, errors.New("некорректное выражение: несоответствие скобок")
			}
			opStack = opStack[:len(opStack)-1] // Убираем '('
			continue
		}

		// если оператор, обрабатываем приоритеты
		if precedence[char] > 0 {
			for len(opStack) > 0 && precedence[char] <= precedence[opStack[len(opStack)-1]] {
				if len(stack) < 2 {
					return 0, errors.New("недостаточно операндов")
				}
				res, err := applyOperator(&stack, &opStack)
				if err != nil {
					return 0, err
				}
				stack = append(stack, res)
			}
			opStack = append(opStack, char)
		}
	}

	// если числа остались в буфере, добавляем их в стек
	if len(numberBuffer) > 0 {
		num, err := strconv.ParseFloat(string(numberBuffer), 64)
		if err != nil {
			return 0, err
		}
		stack = append(stack, num)
	}

	// применяем оставшиеся операторы
	for len(opStack) > 0 {
		if len(stack) < 2 {
			return 0, errors.New("недостаточно операндов")
		}
		res, err := applyOperator(&stack, &opStack)
		if err != nil {
			return 0, err
		}
		stack = append(stack, res)
	}

	if len(stack) != 1 {
		return 0, errors.New("лишние операнды")
	}

	return stack[0], nil
}

// applyOperator применяет оператор к двум верхним элементам стека
func applyOperator(stack *[]float64, opStack *[]rune) (float64, error) {
	if len(*stack) < 2 {
		return 0, errors.New("недостаточно операндов")
	}

	b := (*stack)[len(*stack)-1]
	a := (*stack)[len(*stack)-2]
	*stack = (*stack)[:len(*stack)-2] // Убираем два последних операнда

	op := (*opStack)[len(*opStack)-1]
	*opStack = (*opStack)[:len(*opStack)-1] // Убираем оператор

	var result float64
	switch op {
	case '+':
		result = a + b
	case '-':
		result = a - b
	case '*':
		result = a * b
	case '/':
		if b == 0 {
			return 0, errors.New("ошибка: деление на ноль")
		}
		result = a / b
	default:
		return 0, errors.New("некорректное выражение: неизвестный оператор")
	}

	return result, nil
}

func main() {
	expression := "2+2*2"
	result, err := Calc(expression)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println(result)
	}
}
