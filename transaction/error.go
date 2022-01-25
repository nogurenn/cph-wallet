package transaction

type TransactionEntryMismatch struct {
	error
}

func (e *TransactionEntryMismatch) Error() string {
	return "some or all entries do not match given transaction id"
}

var ErrTransactionEntryMismatch = &TransactionEntryMismatch{}
