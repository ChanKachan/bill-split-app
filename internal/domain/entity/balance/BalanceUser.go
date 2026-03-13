package balance

type BalanceUser struct {
	UserId  int     `json:"user_id"`  // Идентификатор пользователя
	Amount  float32 `json:"amount"`   // Сумма в данной событии (положительое - ему должны, отрицательное он должен)
	GroupId int     `json:"group_id"` // События в котором происходит рассчет
}
