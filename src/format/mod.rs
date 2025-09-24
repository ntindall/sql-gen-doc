use anyhow::Result;
use sqlx::{MySql, MySqlPool};

mod column_description;
mod database_analyzer;
mod foreign_description;
mod index_description;
mod markdown_generator;

pub use column_description::{ColumnDescription, TableInfo};
pub use database_analyzer::DatabaseAnalyzer;
pub use foreign_description::{ForeignDescription, ForeignDescriptions};
pub use index_description::{IndexDescription, IndexDescriptions, LogicalIndex};
pub use markdown_generator::{create_table_markdown, write_to_file};

#[derive(Debug, Clone)]
pub struct DatabaseConnection {
    pool: MySqlPool,
}

impl DatabaseConnection {
    pub async fn new(dsn: &str) -> Result<Self> {
        let pool = MySqlPool::connect(dsn).await?;
        Ok(Self { pool })
    }
    
    pub fn pool(&self) -> &MySqlPool {
        &self.pool
    }
}
