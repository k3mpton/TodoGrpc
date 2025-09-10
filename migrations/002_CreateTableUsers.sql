-- +goose Up
CREATE TABLE Users(
    id serial primary key
); 

-- +goose Down 
drop table Users; 