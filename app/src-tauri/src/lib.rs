use crate::error::{AppError, AppResult};
use crate::utils::compress::compress_files;
use crate::utils::dicom::get_series_element;
use anyhow::anyhow;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::path::PathBuf;
use tauri_plugin_http::{reqwest, reqwest::multipart::Form};
use tokio::fs;

mod error;
mod utils;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_http::init())
        .invoke_handler(tauri::generate_handler![
            upload,
            list_leaf_folders,
            process_dir
        ])
        .setup(|app| {
            if cfg!(debug_assertions) {
                app.handle().plugin(
                    tauri_plugin_log::Builder::default()
                        .level(log::LevelFilter::Info)
                        .build(),
                )?;
            }
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

#[tauri::command]
async fn upload(url: String, folder: String, form: HashMap<String, String>) -> AppResult<String> {
    let mut files = Vec::new();
    let mut entries = fs::read_dir(PathBuf::from(folder)).await?;

    while let Some(entry) = entries.next_entry().await? {
        let path = entry.path();
        files.push(path);
    }

    let (zip_path, _dir) = compress_files(files)?;
    println!("返回的目录是：{:?}", zip_path);

    let mut upload_form = Form::new();
    for (key, value) in form {
        upload_form = upload_form.text(key, value);
    }
    upload_form = upload_form.file("file", zip_path).await?;

    let client = reqwest::Client::new();
    let res = client.post(url).multipart(upload_form).send().await?;

    if res.status().is_success() {
        Ok("success".to_string())
    } else {
        Err(AppError::Anyhow(anyhow!("Upload failed!")))
    }
}

#[tauri::command]
async fn list_leaf_folders(root: String) -> AppResult<Vec<String>> {
    let mut result = Vec::new();
    let mut stack = vec![PathBuf::from(root)];

    while let Some(dir) = stack.pop() {
        let mut is_leaf = true;

        if let Ok(mut entries) = fs::read_dir(&dir).await {
            while let Some(entry) = entries.next_entry().await? {
                let path = entry.path();
                if fs::metadata(&path).await?.is_dir() {
                    is_leaf = false;
                    stack.push(path);
                }
            }
        }

        if is_leaf {
            result.push(dir.to_string_lossy().to_string());
        }
    }

    Ok(result)
}

#[derive(Debug, Serialize, Deserialize)]
struct ProcessDirResponse {
    files: Vec<String>,
    series: String,
}

#[tauri::command]
async fn process_dir(root: String) -> AppResult<Option<String>> {
    let dir_path = PathBuf::from(root);

    let series_uid = tokio::task::spawn_blocking(move || -> AppResult<Option<String>> {
        for entry in std::fs::read_dir(dir_path)? {
            let path = entry?.path();
            if let Some(uid) = get_series_element(&path) {
                return Ok(Some(uid));
            }
        }
        Ok(None)
    })
    .await?;
    series_uid
}
