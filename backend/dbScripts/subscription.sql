CREATE TABLE payments (
    id                  uuid PRIMARY KEY,
    user_id             VARCHAR(255) NOT NULL,
    subscription_type   INT,
    active_till_to      DATE NOT NULL,
    status              VARCHAR(80) NOT NULL CHECK (status in ('active', 'cancelled'))
);

CREATE TABLE subscriptions (
    id                  INT PRIMARY KEY,
    name                VARCHAR(255) NOT NULL,
    description         VARCHAR(1023) NOT NULL,
    price               INT NOT NULL
);

INSERT INTO subscriptions VALUES 
(0, 'Light', 'This is light subscription', 0), 
(1, 'Prime', 'This is primary subscription. Allows to chat with anyone you want.', 1500);