-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS random_times (
  id                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  created           TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  timestamp         DATETIME(6) NOT NULL
) ENGINE=INNODB, CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE IF NOT EXISTS persons (
  person_id         INT,
  last_name         VARCHAR(255),
  first_name        VARCHAR(255),
  address           VARCHAR(255),
  city              VARCHAR(255)
) ENGINE=INNODB, CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE IF NOT EXISTS all_data_types (
  char_             CHAR(12) NOT NULL PRIMARY KEY,
  varchar_          VARCHAR(120),
  boolean_          BOOLEAN,
  tinyint_          TINYINT UNSIGNED NOT NULL,
  smallint_         SMALLINT NOT NULL DEFAULT 31,
  int_              INTEGER,
  decimal_          DECIMAL(3),
  numeric_          NUMERIC(3,2) NOT NULL,
  real_             REAL NOT NULL DEFAULT 0,
  float_            FLOAT(6),
  double_p          DOUBLE PRECISION,
  date_             DATE,
  time_             TIME,
  ts                TIMESTAMP(6),
  bin_              BINARY,
  varbin_           VARBINARY(255),
  blob_             BLOB,
  longblog_         LONGBLOB
) ENGINE=INNODB, CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS all_data_types;
DROP TABLE IF EXISTS persons;
DROP TABLE IF EXISTS random_times;
