-- +goose Up
-- +goose StatementBegin
CREATE TABLE public."cost" (
                               id serial4 NOT NULL,
                               user_id int4 NOT NULL,
                               group_id int4 NOT NULL,
                               description text NOT NULL,
                               sum float4 NOT NULL,
                               CONSTRAINT cost_id_pk PRIMARY KEY (id),
                               CONSTRAINT cost_user_fk FOREIGN KEY (user_id) REFERENCES public."user"(id),
                               CONSTRAINT cost_group_fk FOREIGN KEY (group_id) REFERENCES public."group"(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public."cost";
-- +goose StatementEnd
