package entity

import "time"

type HashKind string

const (
	HashKindAverage    = "AVERAGE"
	HashKindPerception = "PERCEPTION"
	HashKindDifference = "DIFFERENCE"
)

type Meme struct {
	Hash       uint64    `db:"hash"`
	HashKind   HashKind  `db:"hash_kind"`
	ChatID     int64     `db:"chat_id"`
	MessageID  int       `db:"message_id"`
	InsertedAt time.Time `db:"inserted_at"`
}
