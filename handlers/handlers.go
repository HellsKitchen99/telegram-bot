package handlers

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	dataBase "tgBotWithDataBase/db"
	"tgBotWithDataBase/models"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

type BotSender interface {
	Send(telegram.Chattable) (telegram.Message, error)
}

func HandlerExpensesByCategories(bot BotSender, update *telegram.Update, db *sql.DB) error {
	expByC, err := dataBase.ExpensesByCategories(db)
	if err != nil {
		str := fmt.Sprintf("ошибка получения трат по категориям - %v", err)
		msg := telegram.NewMessage(update.CallbackQuery.Message.Chat.ID, str)
		bot.Send(msg)
		return fmt.Errorf("%w", err)
	}
	str := "Траты по категориям\n"
	for category, amount := range expByC {
		out := fmt.Sprintf("%v - %v", category, amount)
		str += out + "\n"
	}
	msg := telegram.NewMessage(update.CallbackQuery.Message.Chat.ID, str)
	_, errMsg := bot.Send(msg)
	if errMsg != nil {
		fmt.Println("ошибка отправки сообщения -", errMsg)
	}
	fmt.Printf("Длина сообщения: %d символов\n", len(str))
	return nil
}

func HandlerExpensesAllTime(bot BotSender, update *telegram.Update, db *sql.DB) error {
	expensesByCategoriesAllTime, err := dataBase.GetExpensesAllTime(db)
	if err != nil {
		str := fmt.Sprintf("ошибка получения трат за все время - %v", err)
		msg := telegram.NewMessage(update.CallbackQuery.Message.Chat.ID, str)
		bot.Send(msg)
		return fmt.Errorf("%w", err)
	}
	str := "Траты за все время\n\n"
	for index, expense := range expensesByCategoriesAllTime {
		str += fmt.Sprintf("Трата №%v\n", index+1)
		str += "ID - " + strconv.Itoa(expense.ID) + "\n"
		str += "Сумма - " + strconv.FormatFloat(expense.Amount, 'f', 2, 64) + "\n"
		str += "Категория - " + expense.Category + "\n"
		str += "Коммент - " + expense.Comment + "\n"
		str += "Время - " + fmt.Sprintf("%v\n\n", expense.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	msg := telegram.NewMessage(update.CallbackQuery.Message.Chat.ID, str)
	_, errMsg := bot.Send(msg)
	if errMsg != nil {
		fmt.Println("ошибка отправки сообщения -", errMsg)
	}
	fmt.Printf("Длина сообщения: %d символов\n", len(str))
	return nil
}

func HandlerAdd(bot BotSender, update *telegram.Update, status map[int64]string, db *sql.DB) (int, time.Time, error) {
	parts := strings.SplitN(update.Message.Text, " ", 3)
	if len(parts) != 3 {
		msg := telegram.NewMessage(update.Message.Chat.ID, "Неверный формат")
		bot.Send(msg)
		delete(status, update.Message.From.ID)
		return 0, time.Time{}, fmt.Errorf("неверный формат")
	}
	amount, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		msg := telegram.NewMessage(update.Message.Chat.ID, "Сумма должна быть числом")
		bot.Send(msg)
		delete(status, update.Message.From.ID)
		return 0, time.Time{}, fmt.Errorf("%w", err)
	}
	category := parts[1]
	if category == "" {
		msg := telegram.NewMessage(update.Message.Chat.ID, "Категория не должна быть пустой")
		bot.Send(msg)
		delete(status, update.Message.From.ID)
		return 0, time.Time{}, fmt.Errorf("пустая категория")
	}
	comment := parts[2]
	newExpense := models.GetExpense{
		Amount:   float64(amount),
		Category: category,
		Comment:  comment,
	}
	out, err := dataBase.AddExpense(db, newExpense)
	if err != nil {
		str := fmt.Sprintf("ошибка запроса к БД - %v", err)
		msg := telegram.NewMessage(update.Message.Chat.ID, str)
		bot.Send(msg)
		delete(status, update.Message.From.ID)
		return 0, time.Time{}, fmt.Errorf("%w", err)
	}
	id := out.ID
	amountOut := out.Amount
	categoryOut := out.Category
	commentOut := out.Comment
	createdAt := out.CreatedAt
	outputStr := fmt.Sprintf("Трата добавлена\nID траты - %v\nСумма - %v\nКатегория - %v\nКомментарий - %v\nВремя траты - %v", id, amountOut, categoryOut, commentOut, createdAt.Format("2006-01-02 15:04:05"))
	msg := telegram.NewMessage(update.Message.Chat.ID, outputStr)
	bot.Send(msg)
	delete(status, update.Message.From.ID)
	fmt.Println(outputStr)
	return id, createdAt, nil
}

func HandlerDelete(bot BotSender, update *telegram.Update, status map[int64]string, db *sql.DB) error {
	id, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		msg := telegram.NewMessage(update.Message.Chat.ID, "ID должен быть целочисленным числом")
		bot.Send(msg)
		delete(status, update.Message.From.ID)
		return fmt.Errorf("%w", err)
	}
	errDel := dataBase.DeleteExpense(db, id)
	if errDel != nil {
		msg := telegram.NewMessage(update.Message.Chat.ID, "ошибка удаления траты")
		bot.Send(msg)
		delete(status, update.Message.From.ID)
		return fmt.Errorf("%w", err)
	}
	msg := telegram.NewMessage(update.Message.Chat.ID, "трата удалена")
	bot.Send(msg)
	delete(status, update.Message.From.ID)
	return nil
}
