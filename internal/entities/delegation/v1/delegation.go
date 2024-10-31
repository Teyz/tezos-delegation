package entities_delegation_v1

import "time"

type Delegation struct {
	Timestamp time.Time `json:"timestamp"`
	Amount    int64     `json:"amount"`
	Delegator string    `json:"delegator"`
	Level     int64     `json:"level"`
}

type Delegation_Create struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Amount    int64     `json:"amount"`
	Delegator string    `json:"delegator"`
	Level     int64     `json:"level"`
}

type TZKTDelegation struct {
	ID        int64       `json:"id"`
	Timestamp time.Time   `json:"timestamp"`
	Amount    int64       `json:"amount"`
	Sender    *TZKTSender `json:"sender"`
	Level     int64       `json:"level"`
	Hash      string      `json:"hash"`
}

type TZKTSender struct {
	Address string `json:"address"`
}
