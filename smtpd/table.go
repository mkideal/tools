package main

const createTableEmail = `CREATE TABLE IF NOT EXISTS 'email' (
	id INT NOT NULL AUTO_INCREMENT,
	from varchar(64) NOT NULL,
	tos text NOT NULL,
	data blob
)`
