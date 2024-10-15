#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn mul6f(original: f64) -> f64 {
    original * 1e6
}
