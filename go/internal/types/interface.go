package types

const (
	ValueTypeList = iota
	ValueTypeString
	ValueTypeSet
	ValueTypeZSet
)

const (
	ValueTypeFancyList   = "list"
	ValueTypeFancyString = "string"
	ValueTypeFancySet    = "set"
	ValueTypeFancyZSet   = "zset"
)

// The item interface. An item is the value of a key.
type Item interface {
	// The pointer to the value.
	Value() interface{}

	// The id of the type of the Item.
	// This need to be constant for the type because it is
	// used when de-/serializing item from/to disk.
	Type() uint64
	TypeFancy() string
}
