#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn add6i(original: i64) -> i64 {
    original + 1980
}
