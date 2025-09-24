use anyhow::Result;
use sqlx::Row;

use super::{ColumnDescription, DatabaseConnection, ForeignDescriptions, IndexDescriptions, TableInfo};

pub struct DatabaseAnalyzer {
    conn: DatabaseConnection,
}

impl DatabaseAnalyzer {
    pub fn new(conn: DatabaseConnection) -> Self {
        Self { conn }
    }

    pub async fn get_tables(&self, db_name: &str) -> Result<Vec<TableInfo>> {
        let rows = sqlx::query(
            "SELECT table_name AS name, table_comment AS comment 
            FROM information_schema.tables 
            WHERE table_schema = ?"
        )
        .bind(db_name)
        .fetch_all(self.conn.pool())
        .await?;

        let mut tables = Vec::new();
        for row in rows {
            let name: String = row.try_get("name")?;
            let comment: String = row.try_get("comment")?;
            tables.push(TableInfo { name, comment });
        }

        Ok(tables)
    }

    pub async fn describe_table(&self, table_name: &str) -> Result<Vec<ColumnDescription>> {
        let rows = sqlx::query(&format!("SHOW FULL COLUMNS FROM {}", table_name))
            .fetch_all(self.conn.pool())
            .await?;

        let mut columns = Vec::new();
        for row in rows {
            let field: String = row.try_get("Field")?;
            let type_name: String = row.try_get("Type")?;
            let null: String = row.try_get("Null")?;
            let key: String = row.try_get("Key")?;
            let default: Option<String> = row.try_get("Default").ok();
            let extra: String = row.try_get("Extra")?;
            let comment: String = row.try_get("Comment")?;
            let collation: Option<String> = row.try_get("Collation").ok();
            let privileges: String = row.try_get("Privileges")?;

            columns.push(ColumnDescription {
                field,
                type_name,
                null,
                key,
                default_value: default,
                extra,
                comment,
                collation,
                privileges,
            });
        }

        Ok(columns)
    }

    pub async fn get_index_descriptions(&self, table_name: &str) -> Result<IndexDescriptions> {
        let rows = sqlx::query(&format!("SHOW INDEXES FROM {}", table_name))
            .fetch_all(self.conn.pool())
            .await?;

        let mut indexes = Vec::new();
        for row in rows {
            let table: String = row.try_get("Table")?;
            let non_unique: bool = row.try_get("Non_unique")?;
            let key_name: String = row.try_get("Key_name")?;
            let seq_in_index: i32 = row.try_get("Seq_in_index")?;
            let column_name: Option<String> = row.try_get("Column_name").ok();
            let comment: Option<String> = row.try_get("Comment").ok();
            let expression: Option<String> = row.try_get("Expression").ok();

            indexes.push(super::index_description::IndexDescription {
                table,
                non_unique,
                key_name,
                seq_in_index,
                column_name,
                comment,
                expression,
            });
        }

        Ok(IndexDescriptions(indexes))
    }

    pub async fn get_foreign_key_descriptions(&self, table_name: &str) -> Result<ForeignDescriptions> {
        let rows = sqlx::query(
            "SELECT table_name AS table_name, column_name AS column_name, 
             constraint_name AS constraint_name, referenced_table_name AS referenced_table_name, 
             referenced_column_name AS referenced_column_name
             FROM information_schema.key_column_usage
             WHERE constraint_schema = DATABASE() AND table_name = ? AND referenced_table_name IS NOT NULL
             ORDER BY 1,2"
        )
        .bind(table_name)
        .fetch_all(self.conn.pool())
        .await?;

        let mut foreign_keys = Vec::new();
        for row in rows {
            let table_name: String = row.try_get("table_name")?;
            let column_name: String = row.try_get("column_name")?;
            let constraint_name: String = row.try_get("constraint_name")?;
            let referenced_table_name: String = row.try_get("referenced_table_name")?;
            let referenced_column_name: String = row.try_get("referenced_column_name")?;

            foreign_keys.push(super::foreign_description::ForeignDescription {
                table_name,
                column_name,
                constraint_name,
                referenced_table_name,
                referenced_column_name,
            });
        }

        Ok(ForeignDescriptions(foreign_keys))
    }
}


