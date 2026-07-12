-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat
(
    id       serial4            NOT NULL,
    del      bool DEFAULT false NOT NULL,
    group_id int4               NOT NULL,
    CONSTRAINT chat_id_pk PRIMARY KEY (id),
    CONSTRAINT chat_id_group_fk FOREIGN KEY (group_id) REFERENCES "group" (id)
);

CREATE TABLE chat_members
(
    id      serial4 NOT NULL,
    chat_id int4    NOT NULL,
    user_id int4    NULL,
    del     bool    NOT NULL,
    CONSTRAINT chat_members_pk PRIMARY KEY (id),
    CONSTRAINT chat_members_chat_fk FOREIGN KEY (chat_id) REFERENCES chat (id),
    CONSTRAINT chat_members_user_fk FOREIGN KEY (user_id) REFERENCES "user" (id)
);


CREATE TABLE chat_message
(
    id          serial4   NOT NULL,
    user_id     int4      NOT NULL,
    message     text      NOT NULL,
    date_create timestamp NOT NULL,
    date_update timestamp NOT NULL,
    chat_id     int4      NOT NULL,
    CONSTRAINT chat_message_pk PRIMARY KEY (id),
    CONSTRAINT chat_message_chat_fk FOREIGN KEY (chat_id) REFERENCES chat (id),
    CONSTRAINT chat_message_user_fk FOREIGN KEY (user_id) REFERENCES "user" (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chat_message;
DROP TABLE IF EXISTS chat_members;
DROP TABLE IF EXISTS chat;
-- +goose StatementEnd
