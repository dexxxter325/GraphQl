CREATE table users(
                      id SERIAL primary key,
                      name varchar(255)not null unique,
                      username varchar(255)not null unique ,
                      password varchar(255) not null
);
CREATE TABLE cars (
                      id SERIAL PRIMARY KEY,
                      userId INT REFERENCES users(id),
                      brand VARCHAR(255) NOT NULL,
                      model VARCHAR(255) NOT NULL,
                      year INT NOT NULL,
                      price FLOAT NOT NULL,
                      mileage INT NOT NULL,
                      description TEXT
);