## all_data_types
#### SCHEMA
|    FIELD    |         TYPE          | NULL  |  KEY  |        DEFAULT         |              EXTRA               |    COMMENT     |
|-------------|-----------------------|-------|-------|------------------------|----------------------------------|----------------|
| `char_`     | `char(12)`            | `NO`  | `PRI` |                        |                                  |                |
| `varchar_`  | `varchar(120)`        | `YES` |       |                        |                                  |                |
| `boolean_`  | `tinyint(1)`          | `YES` |       |                        |                                  |                |
| `tinyint_`  | `tinyint(3) unsigned` | `NO`  |       |                        |                                  |                |
| `smallint_` | `smallint(6)`         | `NO`  |       | `31`                   |                                  | `test comment` |
| `int_`      | `int(11)`             | `YES` |       |                        |                                  |                |
| `decimal_`  | `decimal(3,0)`        | `YES` |       |                        |                                  |                |
| `numeric_`  | `decimal(3,2)`        | `NO`  |       |                        |                                  |                |
| `real_`     | `double`              | `NO`  |       | `0`                    |                                  |                |
| `float_`    | `float`               | `YES` |       |                        |                                  |                |
| `double_p`  | `double`              | `YES` |       |                        |                                  |                |
| `date_`     | `date`                | `YES` |       |                        |                                  |                |
| `time_`     | `time`                | `YES` |       |                        |                                  |                |
| `ts`        | `timestamp(6)`        | `NO`  |       | `CURRENT_TIMESTAMP(6)` | `on update CURRENT_TIMESTAMP(6)` |                |
| `bin_`      | `binary(1)`           | `YES` |       |                        |                                  |                |
| `varbin_`   | `varbinary(255)`      | `YES` |       |                        |                                  |                |
| `blob_`     | `blob`                | `YES` |       |                        |                                  |                |
| `longblog_` | `longblob`            | `YES` |       |                        |                                  |                |
#### INDEXES
| KEY NAME  | UNIQUE |  COLUMNS  | COMMENT |
|-----------|--------|-----------|---------|
| `PRIMARY` | `true` | `(char_)` |         |

## companies
#### SCHEMA
|    FIELD     |   TYPE    | NULL |  KEY  | DEFAULT | EXTRA | COMMENT |
|--------------|-----------|------|-------|---------|-------|---------|
| `company_id` | `int(11)` | `NO` | `PRI` | `0`     |       |         |
#### INDEXES
| KEY NAME  | UNIQUE |    COLUMNS     | COMMENT |
|-----------|--------|----------------|---------|
| `PRIMARY` | `true` | `(company_id)` |         |

## employees

this is a table comment

#### SCHEMA
|    FIELD     |   TYPE    | NULL |  KEY  | DEFAULT | EXTRA | COMMENT |
|--------------|-----------|------|-------|---------|-------|---------|
| `company_id` | `int(11)` | `NO` | `PRI` | `0`     |       |         |
| `person_id`  | `int(11)` | `NO` | `PRI` | `0`     |       |         |
#### INDEXES
|        KEY NAME        | UNIQUE  |          COLUMNS          | COMMENT |
|------------------------|---------|---------------------------|---------|
| `PRIMARY`              | `true`  | `(company_id, person_id)` |         |
| `fk_persons_person_id` | `false` | `(person_id)`             |         |
#### Foreign Key
|      KEY NAME      | TABLE NAME  | COLUMN NAME  |       REFERENCES       |
|--------------------|-------------|--------------|------------------------|
| `employees_ibfk_1` | `employees` | `company_id` | `companies.company_id` |
| `employees_ibfk_2` | `employees` | `person_id`  | `persons.person_id`    |

## goose_db_version
#### SCHEMA
|    FIELD     |         TYPE          | NULL  |  KEY  |       DEFAULT       |      EXTRA       | COMMENT |
|--------------|-----------------------|-------|-------|---------------------|------------------|---------|
| `id`         | `bigint(20) unsigned` | `NO`  | `PRI` |                     | `auto_increment` |         |
| `version_id` | `bigint(20)`          | `NO`  |       |                     |                  |         |
| `is_applied` | `tinyint(1)`          | `NO`  |       |                     |                  |         |
| `tstamp`     | `timestamp`           | `YES` |       | `CURRENT_TIMESTAMP` |                  |         |
#### INDEXES
| KEY NAME  | UNIQUE | COLUMNS | COMMENT |
|-----------|--------|---------|---------|
| `PRIMARY` | `true` | `(id)`  |         |
| `id`      | `true` | `(id)`  |         |

## persons
#### SCHEMA
|    FIELD     |      TYPE      | NULL  |  KEY  | DEFAULT | EXTRA | COMMENT |
|--------------|----------------|-------|-------|---------|-------|---------|
| `person_id`  | `int(11)`      | `NO`  | `PRI` | `0`     |       |         |
| `last_name`  | `varchar(100)` | `YES` | `MUL` |         |       |         |
| `first_name` | `varchar(100)` | `YES` |       |         |       |         |
| `address`    | `varchar(255)` | `YES` |       |         |       |         |
| `city`       | `varchar(255)` | `YES` |       |         |       |         |
#### INDEXES
|           KEY NAME            | UNIQUE  |          COLUMNS          | COMMENT |
|-------------------------------|---------|---------------------------|---------|
| `PRIMARY`                     | `true`  | `(person_id)`             |         |
| `index__last_name`            | `false` | `(last_name)`             |         |
| `index__last_name_first_name` | `false` | `(last_name, first_name)` |         |

## random_times
#### SCHEMA
|    FIELD    |         TYPE          | NULL |  KEY  |        DEFAULT         |      EXTRA       | COMMENT |
|-------------|-----------------------|------|-------|------------------------|------------------|---------|
| `id`        | `bigint(20) unsigned` | `NO` | `PRI` |                        | `auto_increment` |         |
| `created`   | `timestamp(6)`        | `NO` |       | `CURRENT_TIMESTAMP(6)` |                  |         |
| `timestamp` | `datetime(6)`         | `NO` |       |                        |                  |         |
#### INDEXES
| KEY NAME  | UNIQUE | COLUMNS | COMMENT |
|-----------|--------|---------|---------|
| `PRIMARY` | `true` | `(id)`  |         |
