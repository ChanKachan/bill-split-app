-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.operation (
                                  id serial4 NOT NULL,
                                  user_debtor_id int4 NOT NULL,
                                  user_recipient_id int4 NOT NULL,
                                  screen varchar NOT NULL,
                                  create_at timestamp NOT NULL,
                                  group_id int4 NOT NULL,
                                  "money" float4 NOT NULL,
                                  CONSTRAINT operation_id_pk PRIMARY KEY (id),
                                  CONSTRAINT operation_user_fk FOREIGN KEY (user_debtor_id) REFERENCES public."user"(id),
                                  CONSTRAINT operation_user_fk_1 FOREIGN KEY (user_recipient_id) REFERENCES public."user"(id),
                                  CONSTRAINT operation_group_fk FOREIGN KEY (group_id) REFERENCES public."group"(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.operation;
-- +goose StatementEnd
