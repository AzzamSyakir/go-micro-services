package model

import "database/sql"

type Transaction struct {
	Tx    *sql.Tx
	TxErr error
}
