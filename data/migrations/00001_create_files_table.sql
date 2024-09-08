-- +goose Up
CREATE TABLE files (
  `path` TEXT NOT NULL,
  `hash` BLOB NOT NULL,
  `media_type` TEXT NOT NULL,
  `size` INTEGER NOT NULL,
  `mod` TEXT NOT NULL,
  --
  PRIMARY KEY (`path`)
);
-- +goose Down
DROP TABLE files;