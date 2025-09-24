use serde::Deserialize;

#[derive(Debug, Clone, PartialEq)]
pub struct TableInfo {
    pub name: String,
    pub comment: String,
}

#[derive(Debug, Clone, PartialEq)]
pub struct ColumnDescription {
    pub field: String,
    pub type_name: String,
    pub null: String,
    pub key: String,
    pub default_value: Option<String>,
    pub extra: String,
    pub comment: String,
    pub collation: Option<String>,
    pub privileges: String,
}

impl<'de> Deserialize<'de> for ColumnDescription {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        #[derive(Deserialize)]
        struct Helper {
            field: String,
            #[serde(rename = "type")]
            type_name: String,
            null: String,
            key: String,
            #[serde(rename = "default")]
            default_value: Option<String>,
            extra: String,
            comment: String,
            collation: Option<String>,
            privileges: String,
        }

        let helper = Helper::deserialize(deserializer)?;
        
        Ok(ColumnDescription {
            field: helper.field,
            type_name: helper.type_name,
            null: helper.null,
            key: helper.key,
            default_value: helper.default_value,
            extra: helper.extra,
            comment: helper.comment,
            collation: helper.collation,
            privileges: helper.privileges,
        })
    }
}

impl<'de> Deserialize<'de> for TableInfo {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        #[derive(Deserialize)]
        struct Helper {
            name: String,
            comment: String,
        }

        let helper = Helper::deserialize(deserializer)?;
        
        Ok(TableInfo {
            name: helper.name,
            comment: helper.comment,
        })
    }
}


