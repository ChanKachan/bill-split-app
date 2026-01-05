-- +goose Up
-- +goose StatementBegin
CREATE TABLE public."group" (
                                id serial4 NOT NULL,
                                "name" varchar NOT NULL,
                                create_at timestamp NOT NULL,
                                date_start timestamp NOT NULL,
                                date_end timestamp NULL,
                                CONSTRAINT group_id_pk PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public."group";
-- +goose StatementEnd
