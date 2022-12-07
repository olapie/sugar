package mobilex

import (
	"code.olapie.com/sugar/mobilex/nomobile"
	"code.olapie.com/sugar/types"
)

type IntList struct {
	nomobile.List[int]
}

func NewIntList() *IntList {
	return new(IntList)
}

type IntSet struct {
	types.Set[int]
}

func NewIntSet() *IntSet {
	return new(IntSet)
}

type IntE struct {
	Value int
	Error *Error
}

type IntListE struct {
	Value *IntList
	Error *Error
}

type Int16List struct {
	nomobile.List[int16]
}

func NewInt16List() *Int16List {
	return new(Int16List)
}

type Int16Set struct {
	types.Set[int16]
}

func NewInt16Set() *Int16Set {
	return new(Int16Set)
}

type Int16E struct {
	Value int16
	Error *Error
}

type Int16ListE struct {
	Value *Int16List
	Error *Error
}

type Int32List struct {
	nomobile.List[int32]
}

func NewInt32List() *Int32List {
	return new(Int32List)
}

type Int32Set struct {
	types.Set[int32]
}

func NewInt32Set() *Int32Set {
	return new(Int32Set)
}

type Int32E struct {
	Value int32
	Error *Error
}

type Int32ListE struct {
	Value *Int32List
	Error *Error
}

type Int64List struct {
	nomobile.List[int64]
}

func NewInt64List() *Int64List {
	return new(Int64List)
}

type Int64Set struct {
	types.Set[int64]
}

func NewInt64Set() *Int64Set {
	return new(Int64Set)
}

type Int64E struct {
	Value int64
	Error *Error
}

type Int64ListE struct {
	Value *Int64List
	Error *Error
}

type Float64List struct {
	nomobile.List[float64]
}

func NewFloat64List() *Float64List {
	return new(Float64List)
}

type Float64Set struct {
	types.Set[float64]
}

func NewFloat64Set() *Float64Set {
	return new(Float64Set)
}

type Float64E struct {
	Value float64
	Error *Error
}

type Float64ListE struct {
	Value *Float64List
	Error *Error
}

type BoolList struct {
	nomobile.List[bool]
}

func NewBoolList() *BoolList {
	return new(BoolList)
}

type BoolSet struct {
	types.Set[bool]
}

func NewBoolSet() *BoolSet {
	return new(BoolSet)
}

type BoolE struct {
	Value bool
	Error *Error
}

type BoolListE struct {
	Value *BoolList
	Error *Error
}

type StringList struct {
	nomobile.List[string]
}

func NewStringList() *StringList {
	return new(StringList)
}

type StringSet struct {
	types.Set[string]
}

func NewStringSet() *StringSet {
	return new(StringSet)
}

type StringE struct {
	Value string
	Error *Error
}

type StringListE struct {
	Value *StringList
	Error *Error
}

type ByteArrayE struct {
	Value []byte
	Error *Error
}
type IntIntMap struct {
	nomobile.Map[int, int]
}

func NewIntIntMap() *IntIntMap {
	return &IntIntMap{
		Map: *nomobile.NewMap[int, int](),
	}
}

func (m *IntIntMap) Clone() *IntIntMap {
	return &IntIntMap{
		Map: *m.Map.Clone(),
	}
}

func (m *IntIntMap) InsertMap(v *IntIntMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *IntIntMap) Keys() *IntList {
	return &IntList{
		List: *m.Map.Keys(),
	}
}

type IntInt16Map struct {
	nomobile.Map[int, int16]
}

func NewIntInt16Map() *IntInt16Map {
	return &IntInt16Map{
		Map: *nomobile.NewMap[int, int16](),
	}
}

func (m *IntInt16Map) Clone() *IntInt16Map {
	return &IntInt16Map{
		Map: *m.Map.Clone(),
	}
}

func (m *IntInt16Map) InsertMap(v *IntInt16Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *IntInt16Map) Keys() *IntList {
	return &IntList{
		List: *m.Map.Keys(),
	}
}

type IntInt32Map struct {
	nomobile.Map[int, int32]
}

func NewIntInt32Map() *IntInt32Map {
	return &IntInt32Map{
		Map: *nomobile.NewMap[int, int32](),
	}
}

func (m *IntInt32Map) Clone() *IntInt32Map {
	return &IntInt32Map{
		Map: *m.Map.Clone(),
	}
}

func (m *IntInt32Map) InsertMap(v *IntInt32Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *IntInt32Map) Keys() *IntList {
	return &IntList{
		List: *m.Map.Keys(),
	}
}

type IntInt64Map struct {
	nomobile.Map[int, int64]
}

func NewIntInt64Map() *IntInt64Map {
	return &IntInt64Map{
		Map: *nomobile.NewMap[int, int64](),
	}
}

func (m *IntInt64Map) Clone() *IntInt64Map {
	return &IntInt64Map{
		Map: *m.Map.Clone(),
	}
}

func (m *IntInt64Map) InsertMap(v *IntInt64Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *IntInt64Map) Keys() *IntList {
	return &IntList{
		List: *m.Map.Keys(),
	}
}

type IntFloat64Map struct {
	nomobile.Map[int, float64]
}

func NewIntFloat64Map() *IntFloat64Map {
	return &IntFloat64Map{
		Map: *nomobile.NewMap[int, float64](),
	}
}

func (m *IntFloat64Map) Clone() *IntFloat64Map {
	return &IntFloat64Map{
		Map: *m.Map.Clone(),
	}
}

func (m *IntFloat64Map) InsertMap(v *IntFloat64Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *IntFloat64Map) Keys() *IntList {
	return &IntList{
		List: *m.Map.Keys(),
	}
}

type IntBoolMap struct {
	nomobile.Map[int, bool]
}

func NewIntBoolMap() *IntBoolMap {
	return &IntBoolMap{
		Map: *nomobile.NewMap[int, bool](),
	}
}

func (m *IntBoolMap) Clone() *IntBoolMap {
	return &IntBoolMap{
		Map: *m.Map.Clone(),
	}
}

func (m *IntBoolMap) InsertMap(v *IntBoolMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *IntBoolMap) Keys() *IntList {
	return &IntList{
		List: *m.Map.Keys(),
	}
}

type IntStringMap struct {
	nomobile.Map[int, string]
}

func NewIntStringMap() *IntStringMap {
	return &IntStringMap{
		Map: *nomobile.NewMap[int, string](),
	}
}

func (m *IntStringMap) Clone() *IntStringMap {
	return &IntStringMap{
		Map: *m.Map.Clone(),
	}
}

func (m *IntStringMap) InsertMap(v *IntStringMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *IntStringMap) Keys() *IntList {
	return &IntList{
		List: *m.Map.Keys(),
	}
}

type Int16IntMap struct {
	nomobile.Map[int16, int]
}

func NewInt16IntMap() *Int16IntMap {
	return &Int16IntMap{
		Map: *nomobile.NewMap[int16, int](),
	}
}

func (m *Int16IntMap) Clone() *Int16IntMap {
	return &Int16IntMap{
		Map: *m.Map.Clone(),
	}
}

func (m *Int16IntMap) InsertMap(v *Int16IntMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int16IntMap) Keys() *Int16List {
	return &Int16List{
		List: *m.Map.Keys(),
	}
}

type Int16Int16Map struct {
	nomobile.Map[int16, int16]
}

func NewInt16Int16Map() *Int16Int16Map {
	return &Int16Int16Map{
		Map: *nomobile.NewMap[int16, int16](),
	}
}

func (m *Int16Int16Map) Clone() *Int16Int16Map {
	return &Int16Int16Map{
		Map: *m.Map.Clone(),
	}
}

func (m *Int16Int16Map) InsertMap(v *Int16Int16Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int16Int16Map) Keys() *Int16List {
	return &Int16List{
		List: *m.Map.Keys(),
	}
}

type Int16Int32Map struct {
	nomobile.Map[int16, int32]
}

func NewInt16Int32Map() *Int16Int32Map {
	return &Int16Int32Map{
		Map: *nomobile.NewMap[int16, int32](),
	}
}

func (m *Int16Int32Map) Clone() *Int16Int32Map {
	return &Int16Int32Map{
		Map: *m.Map.Clone(),
	}
}

func (m *Int16Int32Map) InsertMap(v *Int16Int32Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int16Int32Map) Keys() *Int16List {
	return &Int16List{
		List: *m.Map.Keys(),
	}
}

type Int16Int64Map struct {
	nomobile.Map[int16, int64]
}

func NewInt16Int64Map() *Int16Int64Map {
	return &Int16Int64Map{
		Map: *nomobile.NewMap[int16, int64](),
	}
}

func (m *Int16Int64Map) Clone() *Int16Int64Map {
	return &Int16Int64Map{
		Map: *m.Map.Clone(),
	}
}

func (m *Int16Int64Map) InsertMap(v *Int16Int64Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int16Int64Map) Keys() *Int16List {
	return &Int16List{
		List: *m.Map.Keys(),
	}
}

type Int16Float64Map struct {
	nomobile.Map[int16, float64]
}

func NewInt16Float64Map() *Int16Float64Map {
	return &Int16Float64Map{
		Map: *nomobile.NewMap[int16, float64](),
	}
}

func (m *Int16Float64Map) Clone() *Int16Float64Map {
	return &Int16Float64Map{
		Map: *m.Map.Clone(),
	}
}

func (m *Int16Float64Map) InsertMap(v *Int16Float64Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int16Float64Map) Keys() *Int16List {
	return &Int16List{
		List: *m.Map.Keys(),
	}
}

type Int16BoolMap struct {
	nomobile.Map[int16, bool]
}

func NewInt16BoolMap() *Int16BoolMap {
	return &Int16BoolMap{
		Map: *nomobile.NewMap[int16, bool](),
	}
}

func (m *Int16BoolMap) Clone() *Int16BoolMap {
	return &Int16BoolMap{
		Map: *m.Map.Clone(),
	}
}

func (m *Int16BoolMap) InsertMap(v *Int16BoolMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int16BoolMap) Keys() *Int16List {
	return &Int16List{
		List: *m.Map.Keys(),
	}
}

type Int16StringMap struct {
	nomobile.Map[int16, string]
}

func NewInt16StringMap() *Int16StringMap {
	return &Int16StringMap{
		Map: *nomobile.NewMap[int16, string](),
	}
}

func (m *Int16StringMap) Clone() *Int16StringMap {
	return &Int16StringMap{
		Map: *m.Map.Clone(),
	}
}

func (m *Int16StringMap) InsertMap(v *Int16StringMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int16StringMap) Keys() *Int16List {
	return &Int16List{
		List: *m.Map.Keys(),
	}
}

type Int32IntMap struct {
	nomobile.Map[int32, int]
}

func NewInt32IntMap() *Int32IntMap {
	return &Int32IntMap{
		Map: *nomobile.NewMap[int32, int](),
	}
}

func (m *Int32IntMap) Clone() *Int32IntMap {
	return &Int32IntMap{
		Map: *m.Map.Clone(),
	}
}

func (m *Int32IntMap) InsertMap(v *Int32IntMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int32IntMap) Keys() *Int32List {
	return &Int32List{
		List: *m.Map.Keys(),
	}
}

type Int32Int16Map struct {
	nomobile.Map[int32, int16]
}

func NewInt32Int16Map() *Int32Int16Map {
	return &Int32Int16Map{
		Map: *nomobile.NewMap[int32, int16](),
	}
}

func (m *Int32Int16Map) Clone() *Int32Int16Map {
	return &Int32Int16Map{
		Map: *m.Map.Clone(),
	}
}

func (m *Int32Int16Map) InsertMap(v *Int32Int16Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int32Int16Map) Keys() *Int32List {
	return &Int32List{
		List: *m.Map.Keys(),
	}
}

type Int32Int32Map struct {
	nomobile.Map[int32, int32]
}

func NewInt32Int32Map() *Int32Int32Map {
	return &Int32Int32Map{
		Map: *nomobile.NewMap[int32, int32](),
	}
}

func (m *Int32Int32Map) Clone() *Int32Int32Map {
	return &Int32Int32Map{
		Map: *m.Map.Clone(),
	}
}

func (m *Int32Int32Map) InsertMap(v *Int32Int32Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int32Int32Map) Keys() *Int32List {
	return &Int32List{
		List: *m.Map.Keys(),
	}
}

type Int32Int64Map struct {
	nomobile.Map[int32, int64]
}

func NewInt32Int64Map() *Int32Int64Map {
	return &Int32Int64Map{
		Map: *nomobile.NewMap[int32, int64](),
	}
}

func (m *Int32Int64Map) Clone() *Int32Int64Map {
	return &Int32Int64Map{
		Map: *m.Map.Clone(),
	}
}

func (m *Int32Int64Map) InsertMap(v *Int32Int64Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int32Int64Map) Keys() *Int32List {
	return &Int32List{
		List: *m.Map.Keys(),
	}
}

type Int32Float64Map struct {
	nomobile.Map[int32, float64]
}

func NewInt32Float64Map() *Int32Float64Map {
	return &Int32Float64Map{
		Map: *nomobile.NewMap[int32, float64](),
	}
}

func (m *Int32Float64Map) Clone() *Int32Float64Map {
	return &Int32Float64Map{
		Map: *m.Map.Clone(),
	}
}

func (m *Int32Float64Map) InsertMap(v *Int32Float64Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int32Float64Map) Keys() *Int32List {
	return &Int32List{
		List: *m.Map.Keys(),
	}
}

type Int32BoolMap struct {
	nomobile.Map[int32, bool]
}

func NewInt32BoolMap() *Int32BoolMap {
	return &Int32BoolMap{
		Map: *nomobile.NewMap[int32, bool](),
	}
}

func (m *Int32BoolMap) Clone() *Int32BoolMap {
	return &Int32BoolMap{
		Map: *m.Map.Clone(),
	}
}

func (m *Int32BoolMap) InsertMap(v *Int32BoolMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int32BoolMap) Keys() *Int32List {
	return &Int32List{
		List: *m.Map.Keys(),
	}
}

type Int32StringMap struct {
	nomobile.Map[int32, string]
}

func NewInt32StringMap() *Int32StringMap {
	return &Int32StringMap{
		Map: *nomobile.NewMap[int32, string](),
	}
}

func (m *Int32StringMap) Clone() *Int32StringMap {
	return &Int32StringMap{
		Map: *m.Map.Clone(),
	}
}

func (m *Int32StringMap) InsertMap(v *Int32StringMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int32StringMap) Keys() *Int32List {
	return &Int32List{
		List: *m.Map.Keys(),
	}
}

type Int64IntMap struct {
	nomobile.Map[int64, int]
}

func NewInt64IntMap() *Int64IntMap {
	return &Int64IntMap{
		Map: *nomobile.NewMap[int64, int](),
	}
}

func (m *Int64IntMap) Clone() *Int64IntMap {
	return &Int64IntMap{
		Map: *m.Map.Clone(),
	}
}

func (m *Int64IntMap) InsertMap(v *Int64IntMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int64IntMap) Keys() *Int64List {
	return &Int64List{
		List: *m.Map.Keys(),
	}
}

type Int64Int16Map struct {
	nomobile.Map[int64, int16]
}

func NewInt64Int16Map() *Int64Int16Map {
	return &Int64Int16Map{
		Map: *nomobile.NewMap[int64, int16](),
	}
}

func (m *Int64Int16Map) Clone() *Int64Int16Map {
	return &Int64Int16Map{
		Map: *m.Map.Clone(),
	}
}

func (m *Int64Int16Map) InsertMap(v *Int64Int16Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int64Int16Map) Keys() *Int64List {
	return &Int64List{
		List: *m.Map.Keys(),
	}
}

type Int64Int32Map struct {
	nomobile.Map[int64, int32]
}

func NewInt64Int32Map() *Int64Int32Map {
	return &Int64Int32Map{
		Map: *nomobile.NewMap[int64, int32](),
	}
}

func (m *Int64Int32Map) Clone() *Int64Int32Map {
	return &Int64Int32Map{
		Map: *m.Map.Clone(),
	}
}

func (m *Int64Int32Map) InsertMap(v *Int64Int32Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int64Int32Map) Keys() *Int64List {
	return &Int64List{
		List: *m.Map.Keys(),
	}
}

type Int64Int64Map struct {
	nomobile.Map[int64, int64]
}

func NewInt64Int64Map() *Int64Int64Map {
	return &Int64Int64Map{
		Map: *nomobile.NewMap[int64, int64](),
	}
}

func (m *Int64Int64Map) Clone() *Int64Int64Map {
	return &Int64Int64Map{
		Map: *m.Map.Clone(),
	}
}

func (m *Int64Int64Map) InsertMap(v *Int64Int64Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int64Int64Map) Keys() *Int64List {
	return &Int64List{
		List: *m.Map.Keys(),
	}
}

type Int64Float64Map struct {
	nomobile.Map[int64, float64]
}

func NewInt64Float64Map() *Int64Float64Map {
	return &Int64Float64Map{
		Map: *nomobile.NewMap[int64, float64](),
	}
}

func (m *Int64Float64Map) Clone() *Int64Float64Map {
	return &Int64Float64Map{
		Map: *m.Map.Clone(),
	}
}

func (m *Int64Float64Map) InsertMap(v *Int64Float64Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int64Float64Map) Keys() *Int64List {
	return &Int64List{
		List: *m.Map.Keys(),
	}
}

type Int64BoolMap struct {
	nomobile.Map[int64, bool]
}

func NewInt64BoolMap() *Int64BoolMap {
	return &Int64BoolMap{
		Map: *nomobile.NewMap[int64, bool](),
	}
}

func (m *Int64BoolMap) Clone() *Int64BoolMap {
	return &Int64BoolMap{
		Map: *m.Map.Clone(),
	}
}

func (m *Int64BoolMap) InsertMap(v *Int64BoolMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int64BoolMap) Keys() *Int64List {
	return &Int64List{
		List: *m.Map.Keys(),
	}
}

type Int64StringMap struct {
	nomobile.Map[int64, string]
}

func NewInt64StringMap() *Int64StringMap {
	return &Int64StringMap{
		Map: *nomobile.NewMap[int64, string](),
	}
}

func (m *Int64StringMap) Clone() *Int64StringMap {
	return &Int64StringMap{
		Map: *m.Map.Clone(),
	}
}

func (m *Int64StringMap) InsertMap(v *Int64StringMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *Int64StringMap) Keys() *Int64List {
	return &Int64List{
		List: *m.Map.Keys(),
	}
}

type StringIntMap struct {
	nomobile.Map[string, int]
}

func NewStringIntMap() *StringIntMap {
	return &StringIntMap{
		Map: *nomobile.NewMap[string, int](),
	}
}

func (m *StringIntMap) Clone() *StringIntMap {
	return &StringIntMap{
		Map: *m.Map.Clone(),
	}
}

func (m *StringIntMap) InsertMap(v *StringIntMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *StringIntMap) Keys() *StringList {
	return &StringList{
		List: *m.Map.Keys(),
	}
}

type StringInt16Map struct {
	nomobile.Map[string, int16]
}

func NewStringInt16Map() *StringInt16Map {
	return &StringInt16Map{
		Map: *nomobile.NewMap[string, int16](),
	}
}

func (m *StringInt16Map) Clone() *StringInt16Map {
	return &StringInt16Map{
		Map: *m.Map.Clone(),
	}
}

func (m *StringInt16Map) InsertMap(v *StringInt16Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *StringInt16Map) Keys() *StringList {
	return &StringList{
		List: *m.Map.Keys(),
	}
}

type StringInt32Map struct {
	nomobile.Map[string, int32]
}

func NewStringInt32Map() *StringInt32Map {
	return &StringInt32Map{
		Map: *nomobile.NewMap[string, int32](),
	}
}

func (m *StringInt32Map) Clone() *StringInt32Map {
	return &StringInt32Map{
		Map: *m.Map.Clone(),
	}
}

func (m *StringInt32Map) InsertMap(v *StringInt32Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *StringInt32Map) Keys() *StringList {
	return &StringList{
		List: *m.Map.Keys(),
	}
}

type StringInt64Map struct {
	nomobile.Map[string, int64]
}

func NewStringInt64Map() *StringInt64Map {
	return &StringInt64Map{
		Map: *nomobile.NewMap[string, int64](),
	}
}

func (m *StringInt64Map) Clone() *StringInt64Map {
	return &StringInt64Map{
		Map: *m.Map.Clone(),
	}
}

func (m *StringInt64Map) InsertMap(v *StringInt64Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *StringInt64Map) Keys() *StringList {
	return &StringList{
		List: *m.Map.Keys(),
	}
}

type StringFloat64Map struct {
	nomobile.Map[string, float64]
}

func NewStringFloat64Map() *StringFloat64Map {
	return &StringFloat64Map{
		Map: *nomobile.NewMap[string, float64](),
	}
}

func (m *StringFloat64Map) Clone() *StringFloat64Map {
	return &StringFloat64Map{
		Map: *m.Map.Clone(),
	}
}

func (m *StringFloat64Map) InsertMap(v *StringFloat64Map) {
	m.Map.InsertMap(&v.Map)
}

func (m *StringFloat64Map) Keys() *StringList {
	return &StringList{
		List: *m.Map.Keys(),
	}
}

type StringBoolMap struct {
	nomobile.Map[string, bool]
}

func NewStringBoolMap() *StringBoolMap {
	return &StringBoolMap{
		Map: *nomobile.NewMap[string, bool](),
	}
}

func (m *StringBoolMap) Clone() *StringBoolMap {
	return &StringBoolMap{
		Map: *m.Map.Clone(),
	}
}

func (m *StringBoolMap) InsertMap(v *StringBoolMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *StringBoolMap) Keys() *StringList {
	return &StringList{
		List: *m.Map.Keys(),
	}
}

type StringStringMap struct {
	nomobile.Map[string, string]
}

func NewStringStringMap() *StringStringMap {
	return &StringStringMap{
		Map: *nomobile.NewMap[string, string](),
	}
}

func (m *StringStringMap) Clone() *StringStringMap {
	return &StringStringMap{
		Map: *m.Map.Clone(),
	}
}

func (m *StringStringMap) InsertMap(v *StringStringMap) {
	m.Map.InsertMap(&v.Map)
}

func (m *StringStringMap) Keys() *StringList {
	return &StringList{
		List: *m.Map.Keys(),
	}
}
