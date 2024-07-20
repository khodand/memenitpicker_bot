create table memes
(
    hash        numeric not null,
    hash_kind   text    not null
        constraint memes_hash_kind_check
            check (hash_kind = ANY (ARRAY ['AVERAGE'::text, 'PERCEPTION'::text, 'DIFFERENCE'::text])),
    chat_id     numeric not null,
    message_id  bigint  not null,
    inserted_at timestamp with time zone default now()
);

alter table memes
    owner to meme_bot;

create unique index memes_hash_kind_chat_id_hash_key
    ON memes (hash_kind, chat_id, hash);

