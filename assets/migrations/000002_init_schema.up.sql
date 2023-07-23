-- https://stackoverflow.com/questions/41970461/how-to-generate-a-random-unique-alphanumeric-id-of-length-n-in-postgres-9-6#41988979
CREATE OR REPLACE FUNCTION generate_uid(size INT) RETURNS TEXT AS
$$
DECLARE
    characters TEXT  := 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    bytes      BYTEA := gen_random_bytes(size);
    l          INT   := length(characters);
    i          INT   := 0;
    output     TEXT  := '';
BEGIN
    WHILE i < size
        LOOP
            output := output || substr(characters, get_byte(bytes, i) % l + 1, 1);
            i := i + 1;
        END LOOP;
    RETURN output;
END;
$$ LANGUAGE plpgsql VOLATILE;

CREATE TABLE ORGANIZATIONS
(
    ID         TEXT PRIMARY KEY   DEFAULT ('OG' || generate_uid(14)),
    NAME       TEXT      NOT NULL,
    TEAM_ID    TEXT      NOT NULL,
    CREATED_AT TIMESTAMP NOT NULL DEFAULT NOW(),
    UPDATED_AT TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE USERS
(

    ID              TEXT PRIMARY KEY   DEFAULT ('US' || generate_uid(14)),
    ORGANIZATION_ID TEXT      NOT NULL,
    USER_ID         TEXT      NOT NULL,
    CREATED_AT      TIMESTAMP NOT NULL DEFAULT NOW(),
    UPDATED_AT      TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_org_id_on_user
        FOREIGN KEY (ORGANIZATION_ID)
            REFERENCES ORGANIZATIONS (ID)
            ON DELETE CASCADE
);

CREATE TABLE ROTAS
(
    ID              TEXT PRIMARY KEY   DEFAULT ('RT' || generate_uid(14)),
    ORGANIZATION_ID TEXT      NOT NULL,
    CHANNEL_ID      TEXT      NOT NULL,
    NAME            TEXT      NOT NULL,
    TYPE            TEXT      NOT NULL,
    CREATED_AT      TIMESTAMP NOT NULL DEFAULT NOW(),
    UPDATED_AT      TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_org_id_on_rota
        FOREIGN KEY (ORGANIZATION_ID)
            REFERENCES ORGANIZATIONS (ID)
            ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_unique_rota_within_team_and_channel ON ROTAS (NAME, CHANNEL_ID, ORGANIZATION_ID);

CREATE TABLE MEMBERS
(
    ID         TEXT PRIMARY KEY   DEFAULT ('RM' || generate_uid(14)),
    ROTA_ID    TEXT      NOT NULL,
    CREATED_AT TIMESTAMP NOT NULL DEFAULT NOW(),
    UPDATED_AT TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_rota_id_on_member
        FOREIGN KEY (ROTA_ID)
            REFERENCES ROTAS (ID)
            ON DELETE CASCADE
);

CREATE TABLE SHIFTS
(
    ID         TEXT PRIMARY KEY   DEFAULT ('RS' || generate_uid(14)),
    ROTA_ID    TEXT      NOT NULL,
    MEMBER_ID  TEXT      NOT NULL,
    START_AT   TIMESTAMP NOT NULL,
    END_AT     TIMESTAMP NOT NULL,
    CREATED_AT TIMESTAMP NOT NULL DEFAULT NOW(),
    UPDATED_AT TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_rota_id_on_schedule
        FOREIGN KEY (ROTA_ID)
            REFERENCES ROTAS (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_member_id_on_schedule
        FOREIGN KEY (MEMBER_ID)
            REFERENCES MEMBERS (id)
            ON DELETE CASCADE
)
