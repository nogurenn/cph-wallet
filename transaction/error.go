package transaction

type TransactionEntryMismatch struct {
	error
}

func (e *TransactionEntryMismatch) Error() string {
	return "some or all entries do not match given transaction id"
}

var ErrTransactionEntryMismatch = &TransactionEntryMismatch{}

type CreditAmountInvalid struct {
	error
}

func (e *CreditAmountInvalid) Error() string {
	return "amount to credit is either zero or negative"
}

var ErrCreditAmountInvalid = &CreditAmountInvalid{}

type BalanceInsufficient struct {
	error
}

func (e *BalanceInsufficient) Error() string {
	return "balance of sender is insufficient"
}

var ErrBalanceInsufficient = &BalanceInsufficient{}

type PaymentSenderReceiverIdentical struct {
	error
}

func (e *PaymentSenderReceiverIdentical) Error() string {
	return "sender and receiver usernames are identical"
}

var ErrPaymentSenderReceiverIdentical = &PaymentSenderReceiverIdentical{}
