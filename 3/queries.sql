----- Получение свободного UDP порта 
SELECT port FROM available_udp_ports WHERE availability = TRUE LIMIT 1;

----- Получение пользователей, находящихся в конкретном лобби 
-- Допустим id лобби f7677b26-c3f9-4028-9fa8-d4732e7e56f6
SELECT id, name FROM users WHERE lobby_id = 'f7677b26-c3f9-4028-9fa8-d4732e7e56f6';

----- Подключение пользователя к конкретному лобби 
-- Допустим 
-- id лобби			f7677b26-c3f9-4028-9fa8-d4732e7e56f6
-- id пользователя 	01d3e2b5-9834-4358-bdbc-5b136f44f166

UPDATE users SET lobby_id = 'f7677b26-c3f9-4028-9fa8-d4732e7e56f6' WHERE id = '01d3e2b5-9834-4358-bdbc-5b136f44f166';