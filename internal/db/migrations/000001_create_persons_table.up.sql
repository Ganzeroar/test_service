CREATE TABLE IF NOT EXISTS persons(
   id serial PRIMARY KEY,
   created_at timestamptz,
   updated_at timestamptz,
   deleted_at timestamptz,
   name VARCHAR (50) NOT NULL,
   surname VARCHAR (50) NOT NULL,
   patronymic VARCHAR (50)
);