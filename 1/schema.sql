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