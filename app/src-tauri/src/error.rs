use serde::{ser::Serializer, Serialize};
use tauri_plugin_http::reqwest;

// create the error type that represents all errors possible in our program
#[derive(Debug, thiserror::Error)]
pub enum AppError {
    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),

    #[error("Zip error: {0}")]
    Zip(#[from] zip::result::ZipError),

    #[error("Tauri error: {0}")]
    Tauri(#[from] tauri::Error),

    #[error("Tokio join error: {0}")]
    TokioJoin(#[from] tokio::task::JoinError),

    #[error("Json error: {0}")]
    JsonError(#[from] serde_json::Error),

    #[error("Anyhow error: {0}")]
    Anyhow(#[from] anyhow::Error),

    #[error("HTTP request error: {0}")]
    Reqwest(#[from] reqwest::Error),
}
impl Serialize for AppError {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        serializer.serialize_str(self.to_string().as_ref())
    }
}

// Todo: error handling
pub type AppResult<T, E = AppError> = anyhow::Result<T, E>;
