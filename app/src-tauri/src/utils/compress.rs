use crate::utils::dicom::get_series_element;
use anyhow::{Context, Result};
use std::fs::File;
use std::io::{Read, Write};
use std::path::PathBuf;
use tempfile::{tempdir, TempDir};

const TEMP_FILE_NAME: &'static str = "TEMP_COMPRESS";
pub fn compress_files(files: Vec<PathBuf>) -> Result<(PathBuf, TempDir)> {
    let dir = tempdir().context("Failed to create temp dir")?;
    println!("临时目录路径: {:?}", dir.path());

    let file_path = dir.path().join(TEMP_FILE_NAME);

    let file = File::create(&file_path)?;

    let mut zip = zip::ZipWriter::new(&file);
    let options = zip::write::SimpleFileOptions::default()
        .compression_method(zip::CompressionMethod::DEFLATE)
        .unix_permissions(0o755);
    for path_str in files {
        let name = path_str
            .file_name()
            .and_then(|n| n.to_str())
            .ok_or_else(|| anyhow::anyhow!("Invalid file name: {:?}", path_str))?;

        if let Some(_) = get_series_element(&path_str) {
            let mut f = File::open(&path_str)?;
            let metadata = file.metadata()?;
            let mut buffer = Vec::with_capacity(metadata.len() as usize);
            f.read_to_end(&mut buffer)?;

            zip.start_file(name, options)?;
            zip.write_all(&buffer)?;
        }
    }

    Ok((file_path, dir))
}
