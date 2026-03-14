package service

import (
	"bill-split/internal/domain/entity/transfer"
	"github.com/gin-gonic/gin"
	"log"
	"math"
)

type OptimizationService interface{ Optimize(c *gin.Context) }

type optimizationService struct {
}

type transaction struct {
	fromId int // от кого
	toId   int // к кому
	amount int // сумма перевода
}

// participant представляет участника группы
type participant struct {
	ID      int
	Name    string
	Balance float64 // >0 — должен получить, <0 — должен заплатить
}

// Метод реализует жадный алгорим.
// Самый большой должник погашает долг самого большого кредитора
func (os *optimizationService) Optimize(c *gin.Context) {
	var participants []transfer.Participant

	if err := c.ShouldBind(&participants); err != nil {
		log.Println("OptimizationService.Optimize err:", err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	transaction := optimizeDebts(participants)
	c.JSON(200, gin.H{
		"code": 200,
		"data": transaction,
	})
	return
}

// OptimizeDebts - жадный алгоритм для минимизации количества транзакций.
// Принимает список участников с их балансами.
// Возвращает список рекомендуемых переводов.
func optimizeDebts(participants []transfer.Participant) []transfer.Transaction {
	// Разделяем на должников (balance < 0) и получателей (balance > 0)
	debtors := make([]*transfer.Participant, 0)
	receivers := make([]*transfer.Participant, 0)

	// Игнорируем участников с нулевым балансом (они никому ничего не должны)
	for i := range participants {
		p := &participants[i]
		if math.Abs(p.Balance) < 0.01 { // Небольшой допуск для ошибок округления
			continue
		}
		if p.Balance < 0 {
			// Для удобства делаем баланс должника положительным числом (сумма долга)
			p.Balance = math.Abs(p.Balance)
			debtors = append(debtors, p)
		} else {
			receivers = append(receivers, p)
		}
	}

	transactions := make([]transfer.Transaction, 0)

	i, j := 0, 0
	for i < len(debtors) && j < len(receivers) {
		debtor := debtors[i]
		receiver := receivers[j]

		// Сумма перевода - минимальная из того, сколько должен должник и сколько должен получить получатель
		transferAmount := math.Min(debtor.Balance, receiver.Balance)

		// Создаем транзакцию
		if transferAmount > 0.01 {
			transactions = append(transactions, transfer.Transaction{
				From:   debtor.Name,
				To:     receiver.Name,
				Amount: transferAmount,
			})
		}

		// Уменьшаем балансы
		debtor.Balance -= transferAmount
		receiver.Balance -= transferAmount

		// Если баланс должника стал почти нулевым, переходим к следующему должнику
		if math.Abs(debtor.Balance) < 0.01 {
			i++
		}
		// Если баланс получателя стал почти нулевым, переходим к следующему получателю
		if math.Abs(receiver.Balance) < 0.01 {
			j++
		}
	}

	return transactions
}
