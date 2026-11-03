package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

var errNotFound = errors.New("Not found")

type Expense struct {
	Id          int     `json:"id"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount,string"`
}

const jsonPath string = "Expenses.json"

// Загрузка трат из json
func loadExpenses() ([]Expense, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Файл не существует")
			return []Expense{}, nil
		}
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	if len(data) == 0 {
		fmt.Println(`Трат еще нет. Для добавления новой записи, введите команду: expense-tracker add --description <description> --amount <amount>`)
		return []Expense{}, nil
	}

	var expenses []Expense
	if err := json.Unmarshal(data, &expenses); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return expenses, nil

}

func add(expenses *[]Expense, description string, amount float64) {
	maxId := func(expense []Expense) int {
		maxId := 0
		for _, e := range *expenses {
			if e.Id > maxId {
				maxId = e.Id
			}
		}
		return maxId
	}(*expenses)

	newExpense := Expense{
		Id:          maxId + 1,
		Date:        time.Now().Format("2006-01-02"),
		Description: description,
		Amount:      amount,
	}

	*expenses = append(*expenses, newExpense)

	newExpenses, err := json.MarshalIndent(*expenses, "", "  ")
	if err != nil {
		fmt.Println("ошибка сериализации:", err)
		return
	}

	err = os.WriteFile(jsonPath, newExpenses, 0644)
	if err != nil {
		fmt.Println("ошибка записи в JSON:", err)
		return
	}

	fmt.Println("# Expense added successfully (ID:" + strconv.Itoa(newExpense.Id) + ")")
}

func update(expenses *[]Expense, id string, description string, amount float64) (string, error) {
	targerId, err := strconv.Atoi(id)
	if err != nil {
		return "Ошибка преобразования строки", err
	}
	for i, e := range *expenses {
		if e.Id == targerId {
			(*expenses)[i].Description = description
			(*expenses)[i].Amount = amount

			newExpenses, err := json.MarshalIndent(*expenses, "", "  ")
			if err != nil {
				return "ошибка сериализации", err
			}

			err = os.WriteFile(jsonPath, newExpenses, 0644)
			if err != nil {
				return "ошибка записи в JSON", err
			}

			return "# Expense updated successfully", nil
		}
	}
	return "", errNotFound
}

func remove(expenses *[]Expense, id string) (string, error) {
	targerId, err := strconv.Atoi(id)
	if err != nil {
		return "Ошибка преобразования строки", err
	}
	for i, e := range *expenses {
		if e.Id == targerId {
			*expenses = append((*expenses)[:i], (*expenses)[i+1:]...)

			newExpenses, err := json.MarshalIndent(*expenses, "", "  ")
			if err != nil {
				return "ошибка сериализации", err
			}

			err = os.WriteFile(jsonPath, newExpenses, 0644)
			if err != nil {
				return "ошибка записи в JSON", err
			}

			return "# Expense deleted successfully", nil
		}
	}
	return "", errNotFound
}

func viewAll(expenses *[]Expense) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	fmt.Fprintln(w, "# ID\tDate\tDescription\tAmount")
	for _, e := range *expenses {
		fmt.Fprintf(w, "# %d\t%s\t%s\t$%.2f\n", e.Id, e.Date, e.Description, e.Amount)
	}
	w.Flush()
}

func viewSumExp(expenses *[]Expense) {
	var sumInCents int
	for _, e := range *expenses {
		sumInCents = sumInCents + int(math.Round(e.Amount*100))
	}
	sumPrintable := float64(sumInCents) / 100
	fmt.Printf("# Total expenses: $%.2f", sumPrintable)
}

func viewSumExpMnth(expenses *[]Expense, mnth string) {

	monthNum, err := strconv.Atoi(mnth)
	if err != nil || monthNum < 1 || monthNum > 12 {
		fmt.Fprintf(os.Stderr, "Ошибка: некорректный номер месяца: %s\n", mnth)
		return
	}

	var sumInCents int
	for _, e := range *expenses {
		t, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			continue
		}
		if int(t.Month()) == monthNum {
			sumInCents += int(math.Round(e.Amount * 100))
		}
	}

	sumPrintable := float64(sumInCents) / 100
	fmt.Printf("# Total expenses: $%.2f", sumPrintable)
}
