use anyhow::Result;

use super::{ColumnDescription, ForeignDescriptions, LogicalIndex};

/// Wrap string in backticks for markdown code formatting
fn wrap_back_ticks(s: &str) -> String {
    if s.is_empty() {
        s.to_string()
    } else {
        format!("`{}`", s)
    }
}

/// Create table markdown with name, comment, columns, indexes, and foreign keys
pub fn create_table_markdown(
    table_name: &str,
    comment: &str,
    columns: &[ColumnDescription],
    indexes: &[LogicalIndex],
    foreign_keys: &ForeignDescriptions,
) -> String {
    let mut markdown = format!("## {}\n", table_name);

    if !comment.is_empty() {
        markdown.push_str(&format!("\n{}\n\n", comment));
    }

    // Schema section - format table as in the Go version
    markdown.push_str("#### SCHEMA\n");
    markdown.push_str("||    FIELD    |         TYPE          | NULL  |  KEY  |        DEFAULT         |              EXTRA               |    COMMENT     |\n");
    markdown.push_str("||-------------|-----------------------|-------|-------|------------------------|----------------------------------|----------------|\n");

    for col in columns {
        let default_str = col.default_value.clone().unwrap_or_default();
        let default_len = default_str.len();
        let extra_len = col.extra.len();
        
        markdown.push_str(&format!(
            "|| `{}`{}| `{}`{} | `{}`{} | `{}`{} | {}| {}| `{}`{} |\n",
            col.field,
            " ".repeat(11_usize.saturating_sub(col.field.len() + 2)),
            col.type_name,
            " ".repeat(24_usize.saturating_sub(col.type_name.len() + 2)), 
            col.null,
            " ".repeat(6_usize.saturating_sub(col.null.len() + 2)),
            col.key,
            " ".repeat(6_usize.saturating_sub(col.key.len() + 2)),
            if default_str.is_empty() {
                " ".repeat(24)
            } else {
                format!(" `{}`{}", default_str, " ".repeat(23_usize.saturating_sub(default_len + 2)))
            },
            if col.extra.is_empty() {
                " ".repeat(33)
            } else {
                format!(" `{}`{}", col.extra, " ".repeat(32_usize.saturating_sub(extra_len + 2)))
            },
            col.comment,
            " ".repeat(15_usize.saturating_sub(col.comment.len() + 2)),
        ));
    }

    markdown.push('\n');

    // Indexes section
    markdown.push_str("#### INDEXES\n");
    // Check if any index has an expression
    let has_expression = indexes.iter().any(|idx| !idx.expression.is_empty());

    if has_expression {
        markdown.push_str("||  KEY NAME  | UNIQUE |  COLUMNS  | COMMENT |  EXPRESSION  |\n");
        markdown.push_str("||------------|--------|-----------|---------|--------------|\n");
    } else {
        markdown.push_str("||  KEY NAME  | UNIQUE |  COLUMNS  | COMMENT |\n");
        markdown.push_str("||------------|--------|-----------|---------|\n");
    }

    for idx in indexes {
        let columns_text = format!("({})", idx.indexed_column_names_ordered.join(", "));
        let comment_str = idx.comment;
        let expression_str = &idx.expression;
        
        let padded_key = format!("`{}`{}", idx.key_name, " ".repeat(11_usize.saturating_sub(idx.key_name.len() + 2)));
        let padded_unique = format!("`{}`{}", (!idx.non_unique).to_string(), " ".repeat(6_usize.saturating_sub(5)));
        let padded_cols = format!("`{}`{}", columns_text, " ".repeat(10_usize.saturating_sub(columns_text.len() + 2)));
        let padded_comment = format!("`{}`{}", comment_str, " ".repeat(8_usize.saturating_sub(comment_str.len() + 2)));

        if has_expression {
            let padded_expr = format!("`{}`", expression_str);
            markdown.push_str(&format!("||{}|{}|{}|{}|{}|\n", 
                padded_key, padded_unique, padded_cols, padded_comment, padded_expr));
        } else {
            markdown.push_str(&format!("||{}|{}|{}|{}|\n", 
                padded_key, padded_unique, padded_cols, padded_comment));
        }
    }

    markdown.push('\n');

    // Foreign keys section
    if !foreign_keys.0.is_empty() {
        markdown.push_str("#### Foreign Key\n");
        markdown.push_str("||      KEY NAME      | TABLE NAME  | COLUMN NAME  |       REFERENCES       |\n");
        markdown.push_str("||--------------------|-------------|--------------|------------------------|\n");

        for fk in &foreign_keys.0 {
            let constraint_name = &fk.constraint_name;
            let table_name = &fk.table_name;
            let column_name = &fk.column_name;
            let references = format!("{}.{}", fk.referenced_table_name, fk.referenced_column_name);

            markdown.push_str(&format!(
                "|| `{}`{} | `{}`{} | `{}`{} | `{}`{} |\n",
                constraint_name,
                " ".repeat(20_usize.saturating_sub(constraint_name.len() + 2)),
                table_name,
                " ".repeat(12_usize.saturating_sub(table_name.len() + 2)),
                column_name,
                " ".repeat(13_usize.saturating_sub(column_name.len() + 2)),
                references,
                " ".repeat(22_usize.saturating_sub(references.len() + 2)),
            ));
        }
        
        markdown.push('\n');
    }

    markdown
}

/// Insert markdown text between sql-gen-doc tags in an existing file
fn insert_between_tags(file: &str, markdown: &str) -> Result<String> {
    let start_tag = "<!-- sql-gen-doc BEGIN -->";
    let end_tag = "<!-- sql-gen-doc END -->";

    if let Some(start_idx) = file.find(start_tag) {
        if let Some(end_idx) = file.find(end_tag) {
            if start_idx >= end_idx {
                return Err(anyhow::anyhow!(
                    "tags out of order! <!-- sql-gen-doc BEGIN --> was after <!-- sql-gen-doc END -->"
                ));
            }

            let start_idx_end = start_idx + start_tag.len();
            Ok(format!(
                "{}\n{}\n{}",
                &file[..start_idx_end],
                markdown,
                &file[end_idx..]
            ))
        } else {
            Err(anyhow::anyhow!("missing end tag <!-- sql-gen-doc END -->"))
        }
    } else if file.contains(end_tag) {
        Err(anyhow::anyhow!("missing start tag <!-- sql-gen-doc BEGIN -->"))
    } else {
        Ok(markdown.to_string())
    }
}

/// Write markdown to file, handling existing tag insertion
pub async fn write_to_file(filename: &str, markdown: &str) -> Result<()> {
    use std::fs;

    let existing_content = match fs::read_to_string(filename) {
        Ok(content) => content,
        Err(error) if error.kind() == std::io::ErrorKind::NotFound => String::new(),
        Err(error) => return Err(anyhow::anyhow!(
            "couldn't open {} for reading. reason: {}", 
            filename, 
            error
        )),
    };

    let processed_markdown = insert_between_tags(&existing_content, markdown)?;

    fs::write(filename, processed_markdown)?;
    Ok(())
}
