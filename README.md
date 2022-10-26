# Why TinySerializer?

### It's Tiny!
As the title says, the package is not that big. It's written with about 500 lines of code, and only works on structs.
It does not work on slices, maps or other types by themselves. It needs a struct to work on.
```go
type TestStruct struct {
	ListMapStruct []map[int]*AllStruct `tiny:"listmapstruct"`
	BoolField     bool                 `tiny:"boolfield"`
	IntField      int64                `tiny:"intfield"`
	FloatField    float64              `tiny:"floatfield"`
	StringField   string               `tiny:"stringfield"`
	ListBool      []bool               `tiny:"listbool"`
	ListInt       []int64              `tiny:"listint"`
	ListFloat     []float64            `tiny:"listfloat"`
	ListString    []string             `tiny:"liststring"`
	MapBool       map[string]bool      `tiny:"mapbool"`
	MapInt        map[string]int64     `tiny:"mapint"`
	MapFloat      map[string]float64   `tiny:"mapfloat"`
	MapString     map[string]string    `tiny:"mapstring"`
	MapListBool   map[string][]bool    `tiny:"maplistbool"`
	MapListInt    map[string][]int64   `tiny:"maplistint"`
	MapListFloat  map[string][]float64 `tiny:"maplistfloat"`
	MapListString map[string][]string  `tiny:"mapliststring"`
	ListMapBool   []map[int]bool       `tiny:"listmapbool"`
	ListMapInt    []map[int]int64      `tiny:"listmapint"`
	ListMapFloat  []map[int]float64    `tiny:"listmapfloat"`
	ListMapString []map[int]string     `tiny:"listmapstring"`
}
```

### Supports GZIP compression
Easily shrink your data by using GZIP compression. It's disabled by default, but can be enabled by using ```Serializer.SetCompress(true)```

### Example:
Create a serializer like so:
```go
var err error
s := NewSerializer()
s, err = s.SetCompress(true)
if err != nil {
	panic(err)
}
data := s.Serialize(&teststruct) // Serialized data
```

And deserialize it like so:
```go
deserialized := MyStruct{}
s = NewSerializer()
s.SetData(serialized)
s = s.SetCompress(true)
s.Deserialize(serialized, &TestStruct)
if err != nil {
	panic(err)
}
```