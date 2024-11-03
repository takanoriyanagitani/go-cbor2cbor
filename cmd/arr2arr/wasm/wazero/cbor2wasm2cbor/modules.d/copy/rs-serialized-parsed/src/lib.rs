use std::sync::RwLock;

use ciborium::Value;

static INPUT_CBOR: RwLock<Vec<u8>> = RwLock::new(vec![]);
static OUTPUT_CBOR: RwLock<Vec<u8>> = RwLock::new(vec![]);

pub fn array2writer<W>(arr: &[Value], wtr: W) -> Result<(), &'static str>
where
    W: ciborium_io::Write,
    W::Error: std::fmt::Debug,
{
    ciborium::into_writer(arr, wtr).map_err(|_| "unable to write the array")
}

pub fn array2vec(arr: &[Value], v: &mut Vec<u8>) -> Result<(), &'static str> {
    array2writer(arr, v)
}

pub fn set_size(v: &mut Vec<u8>, sz: usize, init: u8) -> usize {
    v.resize(sz, init);
    v.capacity()
}

pub fn _set_input_size(sz: u32, init: u8) -> Result<usize, &'static str> {
    let mut guard = INPUT_CBOR.try_write().map_err(|_| "unable to write lock")?;
    let v: &mut Vec<_> = &mut guard;
    Ok(set_size(v, sz as usize, init))
}

pub fn _set_output_size(sz: u32) -> Result<usize, &'static str> {
    let init: u8 = 0;
    let mut guard = OUTPUT_CBOR
        .try_write()
        .map_err(|_| "unable to write lock")?;
    let v: &mut Vec<_> = &mut guard;
    Ok(set_size(v, sz as usize, init))
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn set_input_size(sz: u32, init: u8) -> i32 {
    _set_input_size(sz, init)
        .ok()
        .and_then(|u| u.try_into().ok())
        .unwrap_or(-1)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn set_output_size(sz: u32) -> i32 {
    _set_output_size(sz)
        .ok()
        .and_then(|u| u.try_into().ok())
        .unwrap_or(-1)
}

pub fn _estimate_output_size() -> Result<usize, &'static str> {
    let guard = INPUT_CBOR.try_read().map_err(|_| "unable to read lock")?;
    let s: &[u8] = &guard;
    Ok(s.len())
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn estimate_output_size() -> i32 {
    _estimate_output_size()
        .ok()
        .and_then(|u| u.try_into().ok())
        .unwrap_or(-1)
}

pub fn slice2ptr(s: &[u8]) -> *const u8 {
    s.as_ptr()
}

pub fn _input_offset() -> Result<*const u8, &'static str> {
    let guard = INPUT_CBOR.try_read().map_err(|_| "unable to read lock")?;
    let s: &[u8] = &guard;
    Ok(slice2ptr(s))
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn input_offset() -> *const u8 {
    _input_offset().ok().unwrap_or_else(std::ptr::null)
}

pub fn _output_offset() -> Result<*const u8, &'static str> {
    let guard = OUTPUT_CBOR.try_read().map_err(|_| "unable to read lock")?;
    let s: &[u8] = &guard;
    Ok(slice2ptr(s))
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn output_offset() -> *const u8 {
    _output_offset().ok().unwrap_or_else(std::ptr::null)
}

pub fn serialize_parsed(parsed: &[Value], out: &mut Vec<u8>) -> Result<usize, &'static str> {
    array2vec(parsed, out)?;
    Ok(out.len())
}

pub fn input2output(i: &[u8], out: &mut Vec<u8>) -> Result<usize, &'static str> {
    let val: Value = ciborium::from_reader(i).map_err(|_| "unable to parse")?;
    match val {
        Value::Array(v) => serialize_parsed(&v, out),
        _ => Err("not an array"),
    }
}

pub fn _convert() -> Result<usize, &'static str> {
    let guard = INPUT_CBOR.try_read().map_err(|_| "unable to read lock")?;
    let i: &[u8] = &guard;

    let mut mg = OUTPUT_CBOR
        .try_write()
        .map_err(|_| "unable to write lock")?;
    let mv: &mut Vec<u8> = &mut mg;

    mv.clear();

    input2output(i, mv)
}

#[allow(unsafe_code)]
#[no_mangle]
pub fn convert() -> i32 {
    _convert()
        .ok()
        .and_then(|u| u.try_into().ok())
        .unwrap_or(-1)
}
