package cbor2primitive_test

import (
	"testing"

	"context"

	c2p "github.com/takanoriyanagitani/go-cbor2cbor/cbor2primitive"
)

func TestPrimitive(t *testing.T) {
	t.Parallel()

	t.Run("ConverterMapIxFromFormatString", func(t *testing.T) {
		t.Parallel()

		t.Run("empty", func(t *testing.T) {
			t.Parallel()

			var cmi c2p.ConverterMapIx = c2p.ConverterMapIxFromFormatString("")
			var cnv c2p.CborToAny = cmi(0)

			converted, e := cnv(context.Background(), nil)
			if nil != e {
				t.Fatalf("must be nil: %v\n", e)
			}

			switch b := converted.(type) {
			case []byte:
				switch len(b) {
				case 0:
					break
				default:
					t.Fatalf("unexpected len: %v\n", len(b))
				}
			default:
				t.Fatalf("unexpected type: %v\n", b)
			}
		})

		// ix: 012345
		t.Run("bQqfds", func(t *testing.T) {
			t.Parallel()

			var cmi c2p.ConverterMapIx = c2p.ConverterMapIxFromFormatString(
				"bQqfds",
			)

			t.Run("bool", func(t *testing.T) {
				t.Parallel()

				var cnv c2p.CborToAny = cmi(0)

				t.Run("false", func(t *testing.T) {
					t.Parallel()

					converted, e := cnv(context.Background(), []byte{
						0x00, 0x00,
					})

					if nil != e {
						t.Fatalf("unexpected err: %v\n", e)
					}

					switch b := converted.(type) {
					case bool:
						switch b {
						case true:
							t.Fatalf("expected true: %v\n", b)
						case false:
							break
						}
					default:
						t.Fatalf("unexpected type: %v\n", b)
					}
				})

				t.Run("true", func(t *testing.T) {
					t.Parallel()

					converted, e := cnv(context.Background(), []byte{
						// non-0 bit exists
						0x00, 0x30, 0x00,
					})

					if nil != e {
						t.Fatalf("unexpected err: %v\n", e)
					}

					switch b := converted.(type) {
					case bool:
						switch b {
						case true:
							break
						case false:
							t.Fatalf("expected false: %v\n", b)
						}
					default:
						t.Fatalf("unexpected type: %v\n", b)
					}
				})
			})

			t.Run("uint64", func(t *testing.T) {
				t.Parallel()

				var cnv c2p.CborToAny = cmi(1)

				t.Run("too few bits", func(t *testing.T) {
					t.Parallel()

					_, e := cnv(context.Background(), []byte{
						0x00, 0x00,
					})

					if nil == e {
						t.Fatalf("must fail\n")
					}
				})

				t.Run("zero", func(t *testing.T) {
					t.Parallel()

					converted, e := cnv(context.Background(), []byte{
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
					})

					if nil != e {
						t.Fatalf("unexpected err: %v\n", e)
					}

					switch u := converted.(type) {
					case uint64:
						switch u {
						case 0:
							break
						default:
							t.Fatalf("unexpected value: %v\n", u)
						}
						break
					default:
						t.Fatalf("unexpected type: %v\n", u)
					}
				})

				t.Run("fuji", func(t *testing.T) {
					t.Parallel()

					converted, e := cnv(context.Background(), []byte{
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
						0x37,
						0x76,
					})

					if nil != e {
						t.Fatalf("unexpected err: %v\n", e)
					}

					switch u := converted.(type) {
					case uint64:
						switch u {
						case 0x3776:
							break
						default:
							t.Fatalf("unexpected value: %v\n", u)
						}
						break
					default:
						t.Fatalf("unexpected type: %v\n", u)
					}
				})

				t.Run("sky", func(t *testing.T) {
					t.Parallel()

					converted, e := cnv(context.Background(), []byte{
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
						0x06,
						0x03,
						0x04,
					})

					if nil != e {
						t.Fatalf("unexpected err: %v\n", e)
					}

					switch u := converted.(type) {
					case uint64:
						switch u {
						case 0x060304:
							break
						default:
							t.Fatalf("unexpected value: %v\n", u)
						}
						break
					default:
						t.Fatalf("unexpected type: %v\n", u)
					}
				})
			})

			t.Run("int64", func(t *testing.T) {
				t.Parallel()

				var cnv c2p.CborToAny = cmi(2)

				t.Run("too few bits", func(t *testing.T) {
					t.Parallel()

					_, e := cnv(context.Background(), []byte{
						0x00, 0x00,
					})

					if nil == e {
						t.Fatalf("must fail\n")
					}
				})

				t.Run("zero", func(t *testing.T) {
					t.Parallel()

					converted, e := cnv(context.Background(), []byte{
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
					})

					if nil != e {
						t.Fatalf("unexpected err: %v\n", e)
					}

					switch i := converted.(type) {
					case int64:
						switch i {
						case 0:
							break
						default:
							t.Fatalf("unexpected value: %v\n", i)
						}
						break
					default:
						t.Fatalf("unexpected type: %v\n", i)
					}
				})

				t.Run("fuji", func(t *testing.T) {
					t.Parallel()

					converted, e := cnv(context.Background(), []byte{
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
						0x37,
						0x76,
					})

					if nil != e {
						t.Fatalf("unexpected err: %v\n", e)
					}

					switch i := converted.(type) {
					case int64:
						switch i {
						case 0x3776:
							break
						default:
							t.Fatalf("unexpected value: %v\n", i)
						}
						break
					default:
						t.Fatalf("unexpected type: %v\n", i)
					}
				})

				t.Run("neg", func(t *testing.T) {
					t.Parallel()

					converted, e := cnv(context.Background(), []byte{
						0xff,
						0xff,
						0xff,
						0xff,
						0xff,
						0xff,
						0xff,
						0xff,
					})

					if nil != e {
						t.Fatalf("unexpected err: %v\n", e)
					}

					switch i := converted.(type) {
					case int64:
						switch i {
						case -1:
							break
						default:
							t.Fatalf("unexpected value: %v\n", i)
						}
						break
					default:
						t.Fatalf("unexpected type: %v\n", i)
					}
				})
			})

			t.Run("float32", func(t *testing.T) {
				t.Parallel()

				var cnv c2p.CborToAny = cmi(3)

				t.Run("too few bits", func(t *testing.T) {
					t.Parallel()

					_, e := cnv(context.Background(), []byte{
						0x00, 0x00,
					})

					if nil == e {
						t.Fatalf("must fail\n")
					}
				})

				t.Run("zero", func(t *testing.T) {
					t.Parallel()

					converted, e := cnv(context.Background(), []byte{
						0x00,
						0x00,
						0x00,
						0x00,
					})

					if nil != e {
						t.Fatalf("unexpected err: %v\n", e)
					}

					switch f := converted.(type) {
					case float32:
						switch f {
						case 0:
							break
						default:
							t.Fatalf("unexpected value: %v\n", f)
						}
						break
					default:
						t.Fatalf("unexpected type: %v\n", f)
					}
				})

				t.Run("run", func(t *testing.T) {
					t.Parallel()

					converted, e := cnv(context.Background(), []byte{
						0x42,
						0x28,
						0xc7,
						0xae,
					})

					if nil != e {
						t.Fatalf("unexpected err: %v\n", e)
					}

					switch f := converted.(type) {
					case float32:
						switch f {
						case 42.195:
							break
						default:
							t.Fatalf("unexpected value: %v\n", f)
						}
						break
					default:
						t.Fatalf("unexpected type: %v\n", f)
					}
				})
			})

			t.Run("float64", func(t *testing.T) {
				t.Parallel()

				var cnv c2p.CborToAny = cmi(4)

				t.Run("too few bits", func(t *testing.T) {
					t.Parallel()

					_, e := cnv(context.Background(), []byte{
						0x00, 0x00, 0x00, 0x00,
					})

					if nil == e {
						t.Fatalf("must fail\n")
					}
				})

				t.Run("zero", func(t *testing.T) {
					t.Parallel()

					converted, e := cnv(context.Background(), []byte{
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
						0x00,
					})

					if nil != e {
						t.Fatalf("unexpected err: %v\n", e)
					}

					switch f := converted.(type) {
					case float64:
						switch f {
						case 0:
							break
						default:
							t.Fatalf("unexpected value: %v\n", f)
						}
						break
					default:
						t.Fatalf("unexpected type: %v\n", f)
					}
				})

				t.Run("light", func(t *testing.T) {
					t.Parallel()

					converted, e := cnv(context.Background(), []byte{
						0x41,
						0xb1,
						0xde,
						0x78,
						0x4a,
						0x00,
						0x00,
						0x00,
					})

					if nil != e {
						t.Fatalf("unexpected err: %v\n", e)
					}

					switch f := converted.(type) {
					case float64:
						switch f {
						case 299792458.0:
							break
						default:
							t.Fatalf("unexpected value: %v\n", f)
						}
						break
					default:
						t.Fatalf("unexpected type: %v\n", f)
					}
				})
			})

			t.Run("string", func(t *testing.T) {
				t.Parallel()

				var cnv c2p.CborToAny = cmi(5)

				t.Run("empty", func(t *testing.T) {
					t.Parallel()

					converted, e := cnv(context.Background(), []byte(""))

					if nil != e {
						t.Fatalf("unexpected err: %v\n", e)
					}

					switch s := converted.(type) {
					case string:
						if "" != s {
							t.Fatalf("unexpected value: %s\n", s)
						}
					default:
						t.Fatalf("unexpected type: %v\n", s)
					}
				})

				t.Run("helo", func(t *testing.T) {
					t.Parallel()

					converted, e := cnv(context.Background(), []byte("helo"))

					if nil != e {
						t.Fatalf("unexpected err: %v\n", e)
					}

					switch s := converted.(type) {
					case string:
						if "helo" != s {
							t.Fatalf("unexpected value: %s\n", s)
						}
					default:
						t.Fatalf("unexpected type: %v\n", s)
					}
				})
			})

		})
	})
}
