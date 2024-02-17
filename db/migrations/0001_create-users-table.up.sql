CREATE TABLE IF NOT EXISTS users(
   user_id serial PRIMARY KEY,
   name VARCHAR (50) NOT NULL,
   password VARCHAR (64) NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL
   admin BOOL,
);
