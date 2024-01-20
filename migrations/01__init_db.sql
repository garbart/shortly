BEGIN;

CREATE SCHEMA shortly;

-- users
CREATE TABLE IF NOT EXISTS shortly.users
(
    id           BIGSERIAL PRIMARY KEY,
    email        TEXT NOT NULL,
    passwordHash TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx__users__email on shortly.users (email);

-- tokens
CREATE TABLE IF NOT EXISTS shortly.tokens
(
    id        BIGSERIAL PRIMARY KEY,
    userId    BIGINT      NOT NULL,
    value     TEXT        NOT NULL,
    expiredAt TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_userId FOREIGN KEY (userId) REFERENCES shortly.users (id)
);

CREATE INDEX IF NOT EXISTS idx__tokens__user_id on shortly.tokens (userId);
CREATE UNIQUE INDEX IF NOT EXISTS idx__tokens__value on shortly.tokens (value);

-- urls
CREATE TABLE IF NOT EXISTS shortly.urls
(
    id           BIGSERIAL PRIMARY KEY,
    userId       BIGINT NOT NULL,
    originalLink TEXT   NOT NULL,
    shortLink    TEXT   NOT NULL,
    views        BIGINT NOT NULL DEFAULT 0,

    CONSTRAINT fk_userId FOREIGN KEY (userId) REFERENCES shortly.users (id)
);

CREATE INDEX IF NOT EXISTS idx__urls__user_id on shortly.urls (userid);
CREATE UNIQUE INDEX IF NOT EXISTS idx__urls__shortLink on shortly.urls (shortLink);

COMMIT;