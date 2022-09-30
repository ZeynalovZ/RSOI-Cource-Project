CREATE TABLE chats (
    id                   uuid PRIMARY KEY,
    user_id1             uuid NOT NULL,
    user_id2             uuid NOT NULL
);

CREATE TABLE messages (
    id                   uuid PRIMARY KEY,
    creator_user_id uuid,
    chat_id uuid NOT NULL,
    content VARCHAR(65535) NOT NULL,
    created_at DATE NOT NULL,
    parent_message uuid,
    CONSTRAINT "CHAT_ID_FK" FOREIGN KEY (chat_id)
    REFERENCES chats (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
);

CREATE TABLE messagesStatuses ( 
    id                   uuid PRIMARY KEY,
    message_id uuid NOT NULL,
    user_id uuid NOT NULL,
    status boolean NOT NULL,
    CONSTRAINT "MESSAGE_ID_FK" FOREIGN KEY (message_id)
    REFERENCES messages (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
);
