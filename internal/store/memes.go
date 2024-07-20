package store

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/khodand/memenitpicker_bot/internal/entity"
)

type Memes interface {
	Insert(ctx context.Context, meme *entity.Meme) (*entity.Meme, error)
}

type MemesDB struct {
	conn *sqlx.DB
}

func NewMemes(conn *sqlx.DB) *MemesDB {
	return &MemesDB{conn: conn}
}

func (s *MemesDB) Insert(ctx context.Context, meme *entity.Meme) (*entity.Meme, error) {
	const query = `
INSERT INTO memes (hash,
                   hash_kind,
                   chat_id,
                   message_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT(hash, hash_kind, chat_id) DO
UPDATE SET
    hash = EXCLUDED.hash
RETURNING *;
`
	row := s.conn.QueryRowxContext(
		ctx,
		query,
		meme.Hash,
		meme.HashKind,
		meme.ChatID,
		meme.MessageID,
	)

	var insertedMeme entity.Meme
	err := row.StructScan(&insertedMeme)
	if err != nil {
		return nil, err
	}
	return &insertedMeme, nil
}
