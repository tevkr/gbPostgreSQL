/*
Видео-чат. Пользователи присоединяются к лобби, 
сервер начинает слушать свободный udp порт и отсылать 
кадры всем пользователям, находящимся в лобби.
*/
CREATE TABLE available_udp_ports (
	port INT NOT NULL,
	availability BOOLEAN NOT NULL,
	PRIMARY KEY (port)
);

CREATE TABLE lobbies (
	id TEXT NOT NULL,
	name TEXT NOT NULL,
	password TEXT NOT NULL,
	capacity INT NOT NULL DEFAULT '2',
	users_count INT NOT NULL DEFAULT '1',
	udp_port INT REFERENCES available_udp_ports (port),
	PRIMARY KEY (id)
);

CREATE TABLE users (
	id TEXT NOT NULL,
	name TEXT NOT NULL,
	lobby_id TEXT REFERENCES lobbies (id),
	PRIMARY KEY (id)
);

-- Задание 1
ALTER TABLE lobbies ADD CONSTRAINT lobbies_fk_udp_port FOREIGN KEY (udp_port) REFERENCES available_udp_ports (port);
ALTER TABLE users ADD CONSTRAINT users_fk_lobby_id FOREIGN KEY (lobby_id) REFERENCES lobbies (id);

-- Задание 2
ALTER TABLE lobbies ADD CONSTRAINT lobbies_name_longer_than_four_check CHECK (char_length(name) >= 4);
ALTER TABLE lobbies ADD CONSTRAINT lobbies_capacity_greater_than_two_check CHECK (capacity >= 2);
ALTER TABLE lobbies ADD CONSTRAINT lobbies_users_count_greater_than_one_check CHECK (users_count >= 0);

ALTER TABLE users ADD CONSTRAINT users_name_longer_than_four_check CHECK (char_length(name) >= 4);

-- Задание 3

INSERT INTO available_udp_ports VALUES
(9334, TRUE),
(9335, TRUE),
(9336, TRUE),
(9337, TRUE);

CREATE EXTENSION "uuid-ossp";

INSERT INTO users (id, name) VALUES 
((SELECT uuid_generate_v4())::TEXT, 'name1'),
((SELECT uuid_generate_v4())::TEXT, 'name2'),
((SELECT uuid_generate_v4())::TEXT, 'name3');

WITH available_udp_port AS (
    UPDATE available_udp_ports 
    SET availability = FALSE 
    WHERE 
    	port = (SELECT port FROM available_udp_ports WHERE availability = TRUE LIMIT 1) 
    RETURNING port
)
INSERT INTO lobbies (id, name, password, capacity, users_count, udp_port) VALUES ((SELECT uuid_generate_v4())::TEXT, 'lobby_name1', '', 2, 1, (SELECT * FROM available_udp_port));

-- Добавление пользователя в лобби
UPDATE users SET lobby_id = (SELECT id FROM lobbies LIMIT 1) WHERE id = (SELECT id FROM users LIMIT 1);