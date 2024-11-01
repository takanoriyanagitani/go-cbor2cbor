## Conversion Table

| input type | format char | output type          |
|:----------:|:-----------:|:--------------------:|
| [2]byte    | h           | int16                |
| [2]byte    | H           | uint16               |
| [4]byte    | i           | int32                |
| [4]byte    | I           | uint32               |
| [8]byte    | q           | int64                |
| [8]byte    | Q           | uint64               |
| [4]byte    | f           | float32              |
| [8]byte    | d           | float64              |
| [16]byte   | u           | [2]uint64(uuid)      |
| [16]byte   | U           | 32 byte string(uuid) |
| []byte     | s           | string(utf-8)        |
