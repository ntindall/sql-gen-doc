# a title

<!-- sql-gen-doc BEGIN -->
### all_data_types
| Field     | Type                | Null | Key | Default              | Extra                          |
|-----------|---------------------|------|-----|----------------------|--------------------------------|
| char_     | char(12)            | NO   | PRI | NULL                 |                                |
| varchar_  | varchar(120)        | YES  |     | NULL                 |                                |
| boolean_  | tinyint(1)          | YES  |     | NULL                 |                                |
| tinyint_  | tinyint(3) unsigned | NO   |     | NULL                 |                                |
| smallint_ | smallint(6)         | NO   |     | 31                   |                                |
| int_      | int(11)             | YES  |     | NULL                 |                                |
| decimal_  | decimal(3,0)        | YES  |     | NULL                 |                                |
| numeric_  | decimal(3,2)        | NO   |     | NULL                 |                                |
| real_     | double              | NO   |     | 0                    |                                |
| float_    | float               | YES  |     | NULL                 |                                |
| double_p  | double              | YES  |     | NULL                 |                                |
| date_     | date                | YES  |     | NULL                 |                                |
| time_     | time                | YES  |     | NULL                 |                                |
| ts        | timestamp(6)        | NO   |     | CURRENT_TIMESTAMP(6) | on update CURRENT_TIMESTAMP(6) |
| bin_      | binary(1)           | YES  |     | NULL                 |                                |
| varbin_   | varbinary(255)      | YES  |     | NULL                 |                                |
| blob_     | blob                | YES  |     | NULL                 |                                |
| longblog_ | longblob            | YES  |     | NULL                 |                                |

### goose_db_version
| Field      | Type                | Null | Key | Default           | Extra          |
|------------|---------------------|------|-----|-------------------|----------------|
| id         | bigint(20) unsigned | NO   | PRI | NULL              | auto_increment |
| version_id | bigint(20)          | NO   |     | NULL              |                |
| is_applied | tinyint(1)          | NO   |     | NULL              |                |
| tstamp     | timestamp           | YES  |     | CURRENT_TIMESTAMP |                |

### persons
| Field      | Type         | Null | Key | Default | Extra |
|------------|--------------|------|-----|---------|-------|
| person_id  | int(11)      | YES  |     | NULL    |       |
| last_name  | varchar(255) | YES  |     | NULL    |       |
| first_name | varchar(255) | YES  |     | NULL    |       |
| address    | varchar(255) | YES  |     | NULL    |       |
| city       | varchar(255) | YES  |     | NULL    |       |

### random_times
| Field     | Type                | Null | Key | Default              | Extra          |
|-----------|---------------------|------|-----|----------------------|----------------|
| id        | bigint(20) unsigned | NO   | PRI | NULL                 | auto_increment |
| created   | timestamp(6)        | NO   |     | CURRENT_TIMESTAMP(6) |                |
| timestamp | datetime(6)         | NO   |     | NULL                 |                |

<!-- sql-gen-doc END -->

more stuff to follow
