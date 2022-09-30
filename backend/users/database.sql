CREATE DATABASE users;

CREATE TABLE users
(
    id uuid NOT NULL,
    name text NOT NULL,
    surname text NOT NULL,
    email text NOT NULL,
    phone_number text NOT NULL,
    has_access boolean NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE musics
(
    id uuid NOT NULL,
    name text NOT NULL,
    author text NOT NULL,
    url text NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE images
(
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    image text NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT "USER_ID_FK" FOREIGN KEY (user_id)
        REFERENCES users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

CREATE TABLE users_to_musics
(
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    music_id uuid NOT NULL,
    favourite_level integer NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT "USER_ID_FK" FOREIGN KEY (user_id)
        REFERENCES users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT "MUSIC_ID_FK" FOREIGN KEY (music_id)
        REFERENCES musics (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);