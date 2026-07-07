package service

import (
	"github.com/ChanKachan/bill-split-app/internal/domain/entity/transfer"
	"math"
)

type OptimizationService interface {
	OptimizeDebts(participants []transfer.Participant) []transfer.Transaction
}

type optimizationService struct{}

func NewOptimizationService() OptimizationService {
	return &optimizationService{}
}

// OptimizeDebts - жадный алгоритм для минимизации количества транзакций.
// Принимает список участников с их балансами.
// Возвращает список рекомендуемых переводов.
func (o *optimizationService) OptimizeDebts(participants []transfer.Participant) []transfer.Transaction {
	// Разделяем на должников (balance < 0) и получателей (balance > 0)
	debtors := make([]*transfer.Participant, 0)
	receivers := make([]*transfer.Participant, 0)

	// Игнорируем участников с нулевым балансом (они никому ничего не должны)
	for i := range participants {
		p := &participants[i]
		if math.Abs(p.Balance) < 0.01 { // Небольшой допуск для ошибок округления
			continue
		}
		if p.IsOwes {
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
				FromUserId: debtor.UserId,
				ToUserId:   receiver.UserId,
				Amount:     transferAmount,
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
