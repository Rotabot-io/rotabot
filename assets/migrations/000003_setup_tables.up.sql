CREATE TABLE ROTAS
(
    ID         TEXT PRIMARY KEY   DEFAULT ('RT' || generate_uid(14)),
    TEAM_ID    TEXT      NOT NULL,
    CHANNEL_ID TEXT      NOT NULL,
    NAME       TEXT      NOT NULL,
    METADATA   JSONB     NOT NULL,
    CREATED_AT TIMESTAMP NOT NULL DEFAULT NOW(),
    UPDATED_AT TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_unique_rota_within_team_and_channel ON ROTAS (NAME, CHANNEL_ID, TEAM_ID);

CREATE TRIGGER rotas_updated_at_trigger
    BEFORE UPDATE
    ON ROTAS
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE MEMBERS
(
    ID         TEXT PRIMARY KEY   DEFAULT ('RM' || generate_uid(14)),
    ROTA_ID    TEXT      NOT NULL,
    USER_ID    TEXT      NOT NULL,
    METADATA   JSONB     NOT NULL,
    CREATED_AT TIMESTAMP NOT NULL DEFAULT NOW(),
    UPDATED_AT TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_rota_id_on_member
        FOREIGN KEY (ROTA_ID)
            REFERENCES ROTAS (ID)
            ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_unique_user_within_rota ON MEMBERS (ROTA_ID, USER_ID);
CREATE INDEX idx_user_id_on_members ON MEMBERS (USER_ID);

CREATE TRIGGER members_updated_at_trigger
    BEFORE UPDATE
    ON MEMBERS
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();