use dicom::dictionary_std::tags;
use dicom::object::open_file;
use std::path::PathBuf;

// pub fn has_series_element(f: &File) -> bool {
//     if let Ok(obj) = from_reader(f) {
//         return match obj.element(tags::SERIES_INSTANCE_UID) {
//             Ok(_) => true,
//             Err(_) => false,
//         };
//     }
//     false
// }

pub fn get_series_element(path: &PathBuf) -> Option<String> {
    if let Ok(obj) = open_file(path) {
        return match obj.element(tags::SERIES_INSTANCE_UID) {
            Ok(id) => id.to_str().ok().map(|s| s.to_string()),
            Err(_) => None,
        };
    }
    None
}
