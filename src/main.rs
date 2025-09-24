use anyhow::Result;
use clap::Parser;
use sql_gen_doc::format::DatabaseAnalyzer;
use sql_gen_doc::DatabaseConnection;

#[derive(Parser)]
#[command(name = "sql-gen-doc")]
#[command(about = "A tool to automatically generate SQL documentation")]
struct Args {
    /// A data source name for the database, e.g. mysql://user:password@localhost:3306/database_name
    #[arg(long)]
    dsn: String,
    
    /// The output file to write the documentation to. If not specified, output is written to stdout
    #[arg(short = 'o')]
    output: Option<String>,
    
    /// Outputs tables in alphabetical order
    #[arg(long)]
    sort_tables: bool,
}

#[tokio::main]
async fn main() -> Result<()> {
    let args = Args::parse();
    
    // Validate DSN
    if args.dsn.is_empty() {
        eprintln!("the --dsn flag must be provided");
        std::process::exit(1);
    }

    let db = DatabaseConnection::new(&args.dsn).await?;
    let analyzer = DatabaseAnalyzer::new(db);
    
    let database_name = extract_database_name(&args.dsn)?;
    let mut tables = analyzer.get_tables(&database_name).await?;
    
    if args.sort_tables {
        tables.sort_by(|a, b| a.name.cmp(&b.name));
    }
    
    let mut markdown = String::new();
    for (idx, table) in tables.iter().enumerate() {
        let table_name = &table.name;
        let columns = analyzer.describe_table(table_name).await?;
        let index_data = analyzer.get_index_descriptions(table_name).await?;
        let logical_indexes = index_data.convert_to_logical_indexes()?;
        let foreign_key_data = analyzer.get_foreign_key_descriptions(table_name).await?;
        
        let table_markdown = sql_gen_doc::format::create_table_markdown(
            table_name,
            &table.comment,
            &columns,
            &logical_indexes,
            &foreign_key_data,
        );
        
        markdown.push_str(&table_markdown);
        
        if idx != tables.len() - 1 {
            markdown.push('\n');
        }
    }
    
    if let Some(output_file) = &args.output {
        sql_gen_doc::format::write_to_file(output_file, &markdown).await?;
    } else {
        print!("{}", markdown);
    }
    
    Ok(())
}

fn extract_database_name(dsn: &str) -> Result<String> {
    if let Ok(url) = url::Url::parse(dsn) {
        if let Some(segments) = url.path_segments() {
            let path = segments.collect::<Vec<_>>().join("");
            if !path.is_empty() {
                return Ok(path);
            }
        }
    }
    
    // Fallback: try to extract from MySQL-style DSN
    // Format: mysql://user:password@host:port/database
    if let Ok(url) = url::Url::parse(dsn) {
        let path = url.path();
        if path.starts_with('/') && path.len() > 1 {
            return Ok(path[1..].to_string());
        }
    }
    
    Err(anyhow::anyhow!("Could not extract database name from DSN"))
}


