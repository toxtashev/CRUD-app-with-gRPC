create table books(
	book_id serial primary key,
	author varchar(48),
	name varchar(64),
	price decimal(16,2),
	genre varchar(24),
	cover varchar(16),
	page int
);