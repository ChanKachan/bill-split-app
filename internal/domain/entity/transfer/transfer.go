package transfer

// Участник представляет пользователя с его итоговым балансом.
// Отрицательный баланс означает, что участник должен деньги (должник).
// Положительный баланс означает, что участнику должны деньги (получатель).
type Participant struct {
	Name    string
	Balance float64
}

// Transaction представляет собой рекомендуемую операцию перевода.
type Transaction struct {
	From   string
	To     string
	Amount float64
}
