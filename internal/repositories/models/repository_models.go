package models

type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Coins    int64  `db:"coins"`
}

type Transaction struct {
	ID         int64 `db:"id"`
	SenderID   int64 `db:"sender_id"`
	ReceiverID int64 `db:"receiver_id"`
	Amount     int64 `db:"amount"`
}

type Merch struct {
	ID   int64  `db:"id"`
	Type string `db:"type"`
	Cost int64  `db:"cost"`
}

type Inventory struct {
	ID      int64 `db:"id"`
	UserID  int64 `db:"user_id"`
	MerchID int64 `db:"merch_id"`
	Amount  int64 `db:"amount"`
}
