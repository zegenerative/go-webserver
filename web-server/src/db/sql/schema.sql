CREATE TABLE books
(
  "isbn" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "author" varchar NOT NULL,
  "price" INT NOT NULL
);