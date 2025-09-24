use anyhow::Result;

#[derive(Debug, Clone, PartialEq)]
pub struct IndexDescription {
    pub table: String,
    pub non_unique: bool,
    pub key_name: String,
    pub seq_in_index: i32,
    pub column_name: Option<String>,
    pub comment: Option<String>,
    pub expression: Option<String>,
}

#[derive(Debug)]
pub struct IndexDescriptions(pub Vec<IndexDescription>);

#[derive(Debug, Clone, PartialEq)]
pub struct LogicalIndex {
    pub table: String,
    pub non_unique: bool,
    pub key_name: String,
    pub indexed_column_names_ordered: Vec<String>,
    pub comment: String,
    pub expression: String,
}

impl IndexDescriptions {
    /// Convert raw index descriptions into aggregated type
    pub fn convert_to_logical_indexes(&self) -> Result<Vec<LogicalIndex>> {
        let mut indices = std::collections::HashMap::new();

        for description in &self.0 {
            let li = indices.entry(&description.key_name).or_insert_with(|| LogicalIndex {
                table: description.table.clone(),
                non_unique: description.non_unique,
                key_name: description.key_name.clone(),
                indexed_column_names_ordered: Vec::new(),
                comment: description.comment.clone().unwrap_or_default(),
                expression: description.expression.clone().unwrap_or_default(),
            });

            if let Some(column_name) = &description.column_name {
                // Check if the number of indexed columns equals the sequence number
                if li.indexed_column_names_ordered.len() + 1 != description.seq_in_index as usize {
                    return Err(anyhow::anyhow!(
                        "internal logic error: expecting for indexed columns to always be returned in sequence"
                    ));
                }
                li.indexed_column_names_ordered.push(column_name.clone());
            }
        }

        // Return things in the order they were received for deterministic output
        let mut result = Vec::new();
        let mut added_keys = std::collections::HashSet::new();
        
        for description in &self.0 {
            if !added_keys.contains(&description.key_name) {
                if let Some(li) = indices.get(&description.key_name) {
                    result.push(li.clone());
                    added_keys.insert(description.key_name.clone());
                }
            }
        }

        Ok(result)
    }
}


