package transfer

// Участник представляет пользователя с его итоговым балансом.
// Отрицательный баланс означает, что участник должен деньги (должник).
// Положительный баланс означает, что участнику должны деньги (получатель).
type Participant struct {
	UserId  int
	Balance float64
	IsOwes  bool
}

// Transaction представляет собой рекомендуемую операцию перевода.
type Transaction struct {
	FromUserId int
	ToUserId   int
	Amount     float64
}
