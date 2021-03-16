# v0.4.7 - 2021/02/22

### Fix decoder

* Fix decoding of deep recursive structure
* Fix decoding of embedded unexported pointer field
* Fix invalid test case
* Fix decoding of invalid value
* Fix decoding of prefilled value
* Fix not being able to return UnmarshalTypeError when it should be returned
* Fix decoding of null value
* Fix decoding of type of null string
* Use pre allocated pointer if exists it at decoding

### Reduce memory usage at compile

* Integrate int/int8/int16/int32/int64 and uint/uint8/uint16/uint32/uint64 operation to reduce memory usage at compile

### Remove unnecessary optype
