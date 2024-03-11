
CREATE TABLE IF NOT EXISTS discipline
(
    id              bigserial PRIMARY KEY,
    created_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name            text                        NOT NULL,
    description     text,
    credits         text
);

CREATE TABLE IF NOT EXISTS schedule
(
    id          bigserial PRIMARY KEY,
    created_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    discipline  bigserial                   NOT NULL,
    cabinet     text                        NOT NULL,
    time_period text                        NOT NULL
);

CREATE TABLE IF NOT EXISTS discipline_schedule
(
    id          bigserial PRIMARY KEY,
    created_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    discipline  bigserial,
    schedule    bigserial,
    FOREIGN KEY (discipline)
    REFERENCES discipline(id),
    FOREIGN KEY (schedule)
    REFERENCES schedule(id)
);
