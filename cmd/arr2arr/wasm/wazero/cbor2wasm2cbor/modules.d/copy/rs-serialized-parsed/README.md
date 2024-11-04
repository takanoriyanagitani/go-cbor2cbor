# CBOR -> CBOR(lossless conversion)

- input: CBOR(e.g, CBOR bytes serialized by fxamacker/cbor/v2): can be bigger

  types:
	- bool
    - numbers: uint64, int64, float64
	- []byte
	- string
	- list([]any)
	- map(map[any]any)
	- null
	- time.Time
	- big.Int
	- cbor.Tag
- ouput: CBOR(CBOR bytes serialized by ciborium): can be smaller

  types:
    - Bool(bool)
	- numbers:
	  - Integer
	    - signed: i8, i16, i32, i64, i128
		- unsigned: u8, u16, u32, u64, u128
	  - Float
	- Bytes(Vec<u8>)
	- Text(String)
	- Array(Vec<Value>)
	- Map(Vec<(Value, Value)>)
	- Null
