package main

import (
	"fmt"
	"os"
	database "tgBotWithDataBase/db"
	"tgBotWithDataBase/handlers"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var status map[int64]string = make(map[int64]string)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("ошибка загрузки файла окружения - %v\n", err)
		return
	}

	token := os.Getenv("Token")
	if token == "" {
		fmt.Printf("пустой токен - %v\n", err)
		return
	}

	bot, err := telegram.NewBotAPI(token)
	if err != nil {
		fmt.Printf("ошибка создания бота - %v\n", err)
		return
	}
	fmt.Printf("Бот успешно запущен\n")
	u := telegram.NewUpdate(0)
	u.Timeout = 10

	for {
		updates, err := bot.GetUpdates(u)
		if err != nil {
			fmt.Printf("ошибка при получении апдейта - %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, update := range updates {
			if update.Message != nil {
				fmt.Printf("АПДЕЙТ %d: Message = true, Text = %q\n", update.UpdateID, update.Message.Text)
			} else {
				fmt.Printf("АПДЕЙТ %d: Message = false\n", update.UpdateID)
			}

			if update.Message != nil {
				if status[update.Message.From.ID] == "add" {
					db, err := database.Connect()
					if err != nil {
						str := fmt.Sprintf("ошибка подключения к БД - %v", err)
						msg := telegram.NewMessage(update.Message.Chat.ID, str)
						bot.Send(msg)
						delete(status, update.Message.From.ID)
						continue
					}
					if _, _, err := handlers.HandlerAdd(bot, &update, status, db); err != nil {
						continue
					}
				} else if status[update.Message.From.ID] == "delete" {
					db, err := database.Connect()
					if err != nil {
						str := fmt.Sprintf("ошибка подключения к БД - %v", err)
						msg := telegram.NewMessage(update.Message.Chat.ID, str)
						bot.Send(msg)
						delete(status, update.Message.From.ID)
						continue
					}
					if err := handlers.HandlerDelete(bot, &update, status, db); err != nil {
						continue
					}
				}
				buttons1 := telegram.NewInlineKeyboardRow(
					telegram.NewInlineKeyboardButtonData("Start", "start"),
					telegram.NewInlineKeyboardButtonData("Help", "help"),
				)
				buttons2 := telegram.NewInlineKeyboardRow(
					telegram.NewInlineKeyboardButtonData("Expenses (All time)", "expensesAllTime"),
					telegram.NewInlineKeyboardButtonData("Expenses (By categories)", "expensesByCategories"),
				)
				buttons3 := telegram.NewInlineKeyboardRow(
					telegram.NewInlineKeyboardButtonData("Add expense", "add"),
					telegram.NewInlineKeyboardButtonData("Delete expense", "delete"),
				)
				keyboard := telegram.NewInlineKeyboardMarkup(buttons1, buttons2, buttons3)
				msg := telegram.NewMessage(update.Message.Chat.ID, "Кнопки")
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			}

			if update.CallbackQuery != nil {
				bot.Request(telegram.NewCallback(update.CallbackQuery.ID, ""))

				button := update.CallbackQuery.Data
				switch button {
				case "start":
					msg := telegram.NewMessage(update.CallbackQuery.Message.Chat.ID, "Привет! Это бот для учета расходов. Нажми Help для помощи.")
					bot.Send(msg)
				case "help":
					msg := telegram.NewMessage(update.CallbackQuery.Message.Chat.ID, "Start - начало и приветствие\nHelp - помощь\nExpense (All time)- все траты\nExpenses (By categories) - траты по категориям\nAdd expense - добавить трату\nDelete expense - удалить трату")
					bot.Send(msg)
				case "expensesAllTime":
					db, err := database.Connect()
					if err != nil {
						str := fmt.Sprintf("ошибка подключения к БД - %v", err)
						msg := telegram.NewMessage(update.Message.Chat.ID, str)
						bot.Send(msg)
						delete(status, update.Message.From.ID)
						continue
					}
					if err := handlers.HandlerExpensesAllTime(bot, &update, db); err != nil {
						continue
					}
				case "expensesByCategories":
					db, err := database.Connect()
					if err != nil {
						str := fmt.Sprintf("ошибка подключения к БД - %v", err)
						msg := telegram.NewMessage(update.Message.Chat.ID, str)
						bot.Send(msg)
						delete(status, update.Message.From.ID)
						continue
					}
					if err := handlers.HandlerExpensesByCategories(bot, &update, db); err != nil {
						continue
					}
				case "add":
					msg := telegram.NewMessage(update.CallbackQuery.Message.Chat.ID, "Введите трату (пример 155.55 Категория Коммент)")
					bot.Send(msg)
					status[update.CallbackQuery.From.ID] = "add"
				case "delete":
					msg := telegram.NewMessage(update.CallbackQuery.Message.Chat.ID, "Введите ID траты")
					bot.Send(msg)
					status[update.CallbackQuery.From.ID] = "delete"
				}
			}
			u.Offset = update.UpdateID + 1
		}
	}
}
