package main

import (
	"flag"
	"fmt"
	"os"
)

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

		res, err := remove(&expensesData, *delId)
		if err != nil {
			fmt.Println(res, err)
			os.Exit(1)
		}
		fmt.Println(res)

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

		res, err := update(&expensesData, *upId, *desc, *amount)
		if err != nil {
			fmt.Println(res, err)
			os.Exit(1)
		}
		fmt.Println(res)
	}

}
