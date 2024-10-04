package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

// Функция для определения приоритета операций
func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1 // Сложение и вычитание имеют низкий приоритет
	case '*', '/':
		return 2 // Умножение и деление имеют высокий приоритет
	}
	return 0 // Если операция не распознана, возвращаем 0
}

// Функция для применения операции между двумя числами
func applyOperation(a, b int, op rune) int {
	switch op {
	case '+':
		return a + b // Выполняем сложение
	case '-':
		return a - b // Выполняем вычитание
	case '*':
		return a * b // Выполняем умножение
	case '/':
		if b == 0 {
			panic("division by zero") // Обрабатываем деление на ноль
		}
		return a / b // Выполняем деление
	}
	return 0 // Если операция не распознана, возвращаем 0
}

// Основная функция для вычисления выражения
func evaluateExpression(expr string) int {
	var values []int // Стек для хранения чисел
	var ops []rune   // Стек для хранения операторов

	// Проходим по каждому символу в выражении
	for i := 0; i < len(expr); i++ {
		if expr[i] == ' ' {
			continue // Игнорируем пробелы в выражении
		}

		if expr[i] >= '0' && expr[i] <= '9' { // Проверяем, является ли символ цифрой
			val := 0 // Переменная для хранения текущего числа
			// Пока символы являются цифрами, собираем число
			for i < len(expr) && expr[i] >= '0' && expr[i] <= '9' {
				val = val*10 + int(expr[i]-'0') // Формируем число из цифр
				i++                             // Переходим к следующему символу
			}
			values = append(values, val) // Добавляем число в стек значений
			i--                          // Корректируем индекс, чтобы не пропустить следующий символ
		} else if expr[i] == '(' { // Если встречаем открывающую скобку
			ops = append(ops, '(') // Добавляем её в стек операторов
		} else if expr[i] == ')' { // Если встречаем закрывающую скобку
			// Обрабатываем все операции до открывающей скобки
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				val2 := values[len(values)-1]   // Берём верхнее значение из стека значений
				values = values[:len(values)-1] // Удаляем его из стека

				val1 := values[len(values)-1]   // Берём следующее значение
				values = values[:len(values)-1] // Удаляем его из стека

				op := ops[len(ops)-1]  // Берём верхний оператор из стека операторов
				ops = ops[:len(ops)-1] // Удаляем его из стека

				// Применяем операцию и добавляем результат обратно в стек значений
				values = append(values, applyOperation(val1, val2, op))
			}
			ops = ops[:len(ops)-1] // Удаляем открывающую скобку из стека операторов
		} else { // Если это оператор (например, +, -, *, /)
			// Обрабатываем все операции с более высоким или равным приоритетом
			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(rune(expr[i])) {
				val2 := values[len(values)-1]   // Берём верхнее значение из стека значений
				values = values[:len(values)-1] // Удаляем его из стека

				val1 := values[len(values)-1]   // Берём следующее значение
				values = values[:len(values)-1] // Удаляем его из стека

				op := ops[len(ops)-1]  // Берём верхний оператор из стека операторов
				ops = ops[:len(ops)-1] // Удаляем его из стека

				// Применяем операцию и добавляем результат обратно в стек значений
				values = append(values, applyOperation(val1, val2, op))
			}
			ops = append(ops, rune(expr[i])) // Добавляем текущий оператор в стек операторов
		}
	}

	// После завершения прохода по выражению обрабатываем оставшиеся операции
	for len(ops) > 0 {
		val2 := values[len(values)-1]   // Берём верхнее значение из стека значений
		values = values[:len(values)-1] // Удаляем его из стека

		val1 := values[len(values)-1]   // Берём следующее значение
		values = values[:len(values)-1] // Удаляем его из стека

		op := ops[len(ops)-1]  // Берём верхний оператор из стека операторов
		ops = ops[:len(ops)-1] // Удаляем его из стека
		// Применяем операцию и добавляем результат обратно в стек значений
		values = append(values, applyOperation(val1, val2, op))
	}

	return values[0] // Возвращаем окончательный результат (единственное значение в стеке)
}

func evaluate(w http.ResponseWriter, r *http.Request) {

	var request struct {
		Expression string `json:"expression"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Здесь вы можете обработать выражение и вернуть результат
	// Например:
	result := evaluateExpression(request.Expression)
	json.NewEncoder(w).Encode(struct {
		Result int `json:"result"`
	}{result})
}

func main() {

	http.HandleFunc("/evaluate", evaluate)
	http.Handle("/", http.FileServer(http.Dir("."))) //
	fmt.Println("Сервер запущен на порту 8080")
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, address := range addrs {
		// Проверяем, является ли адрес IP-адресом
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("IP-адрес сервера:", ipnet.IP.String())
			}
		}
	}
	http.ListenAndServe(":8080", nil)

	// var input string                                 // Переменная для хранения ввода пользователя
	// fmt.Println("Введите математическое выражение:") // Запрашиваем ввод у пользователя
	// fmt.Scanln(&input)                               // Читаем ввод

	// result := evaluateExpression(input)   // Вычисляем результат выражения
	// fmt.Printf("Результат: %d\n", result) // Выводим результат на экран
}
