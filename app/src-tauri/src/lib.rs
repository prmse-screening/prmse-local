use crate::error::{AppError, AppResult};
use crate::utils::compress::compress_files;
use crate::utils::dicom::get_series_element;
use anyhow::anyhow;
use futures_util::StreamExt;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::path::PathBuf;
use tauri_plugin_http::{reqwest, reqwest::multipart::Form};
use tokio::{fs, io};

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
            process_dir,
            export
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
async fn list_leaf_folders(roots: Vec<String>) -> AppResult<Vec<String>> {
    let mut result = Vec::new();
    let mut stack = roots.into_iter().map(PathBuf::from).collect::<Vec<PathBuf>>();

    while let Some(dir) = stack.pop() {
        let mut is_leaf = true;

        if let Ok(mut entries) = fs::read_dir(&dir).await {
            while let Some(entry) = entries.next_entry().await? {
                let path = entry.path();
                if let Ok(meta) = fs::metadata(&path).await {
                    if meta.is_dir() {
                        is_leaf = false;
                        stack.push(path);
                    }
                }
            }
        }

        if is_leaf {
            result.push(dir.to_string_lossy().to_string());
        }
    }

    if result.len() > 100 {
        result.truncate(100);
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

    if !dir_path.exists() || !dir_path.is_dir() {
        return Err(AppError::Anyhow(anyhow!("{:?} is not a directory or doesn't exist!", dir_path)));
    }
    let mut dir_entries = tokio::fs::read_dir(dir_path).await?;
    while let Some(entry) = dir_entries.next_entry().await? {
        let path = entry.path();
        if !path.is_file() {
            continue;
        }
        if let Some(uid) = get_series_element(&path) {
            return Ok(Some(uid));
        }
    }
    Ok(None)
}

#[tauri::command]
async fn export(path: String, url: String) -> AppResult<bool> {
    let client = reqwest::Client::new();
    let res = client.get(url).send().await?;

    if res.status().is_success() {
        let mut f = fs::File::create(path).await?;
        let mut stream = res.bytes_stream();
        while let Some(Ok(chunk)) = stream.next().await {
            io::copy(&mut chunk.as_ref(), &mut f).await?;
        }
        Ok(true)
    } else {
        Ok(false)
    }
}
