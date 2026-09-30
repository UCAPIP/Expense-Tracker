package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

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

func update(expenses *[]Expense, id string, description string, amount float64) {
	targerId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("Ошибка преобразования строки ID: %v", err)
		return
	}
	for i, e := range *expenses {
		if e.Id == targerId {
			(*expenses)[i].Description = description
			(*expenses)[i].Amount = amount

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

			fmt.Println("# Expense updated successfully")

		}
	}
}

func delete(expenses *[]Expense, id string) {
	targerId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("Ошибка преобразования строки ID: %v", err)
		return
	}
	for i, e := range *expenses {
		if e.Id == targerId {
			*expenses = append((*expenses)[:i], (*expenses)[i+1:]...)

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

			fmt.Println("# Expense deleted successfully")

		}
	}
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

func main() {

	if len(os.Args) < 2 {
		os.Exit(1)
	}

	expensesData, err := loadExpenses()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка загрузки трат: %v\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		viewAll(&expensesData)

	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)

		desc := addCmd.String("description", "", "description")
		amount := addCmd.Float64("amount", 0.00, "amount")

		addCmd.Parse(os.Args[2:])

		if *desc == "" || *amount <= 0 {
			fmt.Println("Ошибка: необходимо указать --description и --amount (больше 0)")
			addCmd.Usage()
			os.Exit(1)
		}

		add(&expensesData, *desc, *amount)

	case "summary":
		addCmd := flag.NewFlagSet("summary", flag.ExitOnError)

		month := addCmd.String("month", "", "month")

		addCmd.Parse(os.Args[2:])

		if *month == "" {
			viewSumExp(&expensesData)
		} else {
			viewSumExpMnth(&expensesData, *month)
		}

	case "delete":
		addCmd := flag.NewFlagSet("delete", flag.ExitOnError)

		delId := addCmd.String("id", "", "id")

		addCmd.Parse(os.Args[2:])

		if *delId == "" {
			fmt.Println("Ошибка: необходимо указать --id")
			addCmd.Usage()
			os.Exit(1)
		}

		delete(&expensesData, *delId)

	case "update":
		addCmd := flag.NewFlagSet("update", flag.ExitOnError)

		upId := addCmd.String("id", "", "id")
		desc := addCmd.String("description", "", "description")
		amount := addCmd.Float64("amount", 0.00, "amount")

		addCmd.Parse(os.Args[2:])

		if *desc == "" || *amount <= 0 || *upId == "" {
			fmt.Println("Ошибка: необходимо указать --id, --description и --amount (больше 0)")
			addCmd.Usage()
			os.Exit(1)
		}

		update(&expensesData, *upId, *desc, *amount)
	}

}
