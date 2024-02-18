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
create table sessions(
    session_id varchar(255) primary key,
    userId int references users(id),
    created_at TIMESTAMP,
    expires_at TIMESTAMP
);