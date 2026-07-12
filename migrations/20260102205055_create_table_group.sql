-- +goose Up
-- +goose StatementBegin
CREATE TABLE "group"
(
    id         serial4   NOT NULL,
    "name"     varchar   NOT NULL,
    create_at  timestamp NOT NULL,
    date_start timestamp NOT NULL,
    date_end   timestamp NULL,
    CONSTRAINT group_id_pk PRIMARY KEY (id)
);

CREATE TABLE group_members
(
    id          serial4        NOT NULL,
    user_id     int4           NOT NULL,
    group_id    int4           NOT NULL,
    money_spent float4         NOT NULL,
    status      varchar        NOT NULL,
    del         int4 DEFAULT 0 NOT NULL,
    CONSTRAINT group_members_id_pk PRIMARY KEY (id),
    CONSTRAINT group_members_user_fk FOREIGN KEY (user_id) REFERENCES public."user" (id),
    CONSTRAINT group_members_group_fk_1 FOREIGN KEY (group_id) REFERENCES public."group" (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "group";
DROP TABLE IF EXISTS group_members;
-- +goose StatementEnd
