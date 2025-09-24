#[derive(Debug, Clone, PartialEq)]
pub struct ForeignDescription {
    pub table_name: String,
    pub column_name: String,
    pub constraint_name: String,
    pub referenced_table_name: String,
    pub referenced_column_name: String,
}

#[derive(Debug)]
pub struct ForeignDescriptions(pub Vec<ForeignDescription>);


