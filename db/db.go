package db

import (
	"database/sql"
	"fmt"
	"tgBotWithDataBase/models"
)

func Connect() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password='C529@4H0OdEO%}Y' dbname=puvelka sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД - %w", err)
	}
	errPing := db.Ping()
	if errPing != nil {
		return nil, fmt.Errorf("ошибка пинга БД - %w", errPing)
	}
	return db, nil
}

func GetExpensesAllTime(db *sql.DB) ([]models.Expense, error) {
	defer db.Close()
	var ExpensesAllTime []models.Expense
	rows, err := db.Query("SELECT * FROM expensesTG")
	if err != nil {
		return ExpensesAllTime, fmt.Errorf("ошибка запроса к БД - %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(&expense.ID, &expense.Amount, &expense.Category, &expense.Comment, &expense.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения данных из строки - %w", err)
		}
		ExpensesAllTime = append(ExpensesAllTime, expense)
	}

	return ExpensesAllTime, nil
}

func ExpensesByCategories(db *sql.DB) (map[string]float64, error) {
	defer db.Close()
	var ExpensesByCategories map[string]float64 = make(map[string]float64)
	rows, err := db.Query(`SELECT category, SUM(amount)
	FROM expensesTG
	GROUP BY category`)
	if err != nil {
		return ExpensesByCategories, fmt.Errorf("ошибка запроса к БД - %w", err)
	}

	for rows.Next() {
		var category string
		var amount float64
		err := rows.Scan(&category, &amount)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения данных из строки - %w", err)
		}
		ExpensesByCategories[category] = amount
	}
	return ExpensesByCategories, nil
}

func AddExpense(db *sql.DB, exp models.GetExpense) (models.Expense, error) {
	defer db.Close()
	var outputExpense models.Expense
	row := db.QueryRow(`INSERT INTO expensesTG (amount, category, comment)
	VALUES ($1, $2, $3)
	RETURNING id, created_at`, exp.Amount, exp.Category, exp.Comment)
	err := row.Scan(&outputExpense.ID, &outputExpense.CreatedAt)
	if err != nil {
		return outputExpense, fmt.Errorf("ошибка при получении данных из строки - %w", err)
	}
	outputExpense.Amount = exp.Amount
	outputExpense.Category = exp.Category
	outputExpense.Comment = exp.Comment

	return outputExpense, nil
}

func DeleteExpense(db *sql.DB, id int) error {
	defer db.Close()
	res, err := db.Exec("DELETE FROM expensesTG WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("ошибка запроса к БД на удаление - %w", err)
	}
	rowsAffected, errRows := res.RowsAffected()
	if errRows != nil {
		return fmt.Errorf("ошибка получения количества затронутых строк - %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ничего не удалилось")
	}
	return nil
}
