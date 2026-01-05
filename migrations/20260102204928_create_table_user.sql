-- +goose Up
-- +goose StatementBegin
CREATE TABLE public."user" (
                               id serial4 NOT NULL,
                               "name" varchar NOT NULL,
                               email varchar NOT NULL,
                               phone varchar NOT NULL,
                               login varchar NOT NULL,
                               "password" varchar NOT NULL,
                               CONSTRAINT user_id_pk PRIMARY KEY (id)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public."user";
-- +goose StatementEnd
