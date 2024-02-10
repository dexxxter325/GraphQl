CREATE TABLE cars (
                      id SERIAL PRIMARY KEY,
                      brand VARCHAR(255) NOT NULL,
                      model VARCHAR(255) NOT NULL,
                      year INT NOT NULL,
                      price FLOAT NOT NULL,
                      mileage INT NOT NULL,
                      description TEXT NOT NULL
);
