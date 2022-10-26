package tinyserializer

import (
	"testing"
)

type Testie struct {
	Intlist              []int64                `tiny:"intlist"`
	StringList           []string               `tiny:"stringlist"`
	Structie             Structie               `tiny:"structie"`
	TestMap              map[string]string      `tiny:"testmap"`
	TestMapOfLists       map[string][]string    `tiny:"testmapoflists"`
	MapEmbeddedWithLists map[string][][][]int64 `tiny:"mapembeddedwithlists"`
	All                  *AllStruct             `tiny:"all"`
	ListMapStruct        []map[int]*AllStruct   `tiny:"listmapstruct"`
}

// Struct with every type of field
type AllStruct struct {
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

type Structie struct {
	IntList         []int64     `tiny:"intlist"`
	EmbeddedintList [][]int64   `tiny:"embeddedintlist"`
	EmbeddedstrList [][]string  `tiny:"embeddedstrlist"`
	DoubleEmbedded  [][][]int64 `tiny:"doubleembedded"`
}

var BasicIntList = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
var BasicEmbeddedIntList = [][]int64{BasicIntList, BasicIntList, BasicIntList}
var BasicDoubleEmbedded = [][][]int64{BasicEmbeddedIntList, BasicEmbeddedIntList, BasicEmbeddedIntList}
var BasicStringList = []string{"Hello", "World", "This", "Is", "A", "Test"}
var structmap = map[int]*AllStruct{
	1: All_S,
	2: All_S,
}

var All_S = &AllStruct{
	BoolField:     true,
	IntField:      123,
	FloatField:    123.456,
	StringField:   "Hello World",
	ListBool:      []bool{true, false, true, false},
	ListInt:       []int64{1, 2, 3, 4},
	ListFloat:     []float64{1.2, 3.4, 5.6, 7.8},
	ListString:    []string{"Hello", "World"},
	MapBool:       map[string]bool{"Hello": true, "World": false},
	MapInt:        map[string]int64{"Hello": 1, "World": 2},
	MapFloat:      map[string]float64{"Hello": 1.2, "World": 3.4},
	MapString:     map[string]string{"Hello": "World", "Foo": "Bar"},
	MapListBool:   map[string][]bool{"Hello": []bool{true, false}, "World": []bool{false, true}},
	MapListInt:    map[string][]int64{"Hello": []int64{1, 2}, "World": []int64{3, 4}},
	MapListFloat:  map[string][]float64{"Hello": []float64{1.2, 3.4}, "World": []float64{5.6, 7.8}},
	MapListString: map[string][]string{"Hello": []string{"Hello", "World"}, "World": []string{"Foo", "Bar"}},
	ListMapBool:   []map[int]bool{map[int]bool{1: true, 2: false}, map[int]bool{3: true, 4: false}},
	ListMapInt:    []map[int]int64{map[int]int64{1: 1, 2: 2}, map[int]int64{3: 3, 4: 4}},
	ListMapFloat:  []map[int]float64{map[int]float64{1: 1.2, 2: 3.4}, map[int]float64{3: 5.6, 4: 7.8}},
	ListMapString: []map[int]string{map[int]string{1: "Hello", 2: "World"}, map[int]string{3: "Foo", 4: "Bar"}},
}

// Create a test struct
var testStruct = Testie{
	Intlist: []int64{
		1,
		2,
		3,
		4,
	},
	StringList: BasicStringList,
	Structie: Structie{
		IntList:         BasicIntList,
		EmbeddedintList: BasicEmbeddedIntList,
		EmbeddedstrList: [][]string{
			BasicStringList,
			BasicStringList,
		},
		DoubleEmbedded: [][][]int64{
			BasicEmbeddedIntList,
			BasicEmbeddedIntList,
		},
	},
	TestMap: map[string]string{
		"Hello": "World",
		"Foo":   "Bar",
	},
	TestMapOfLists: map[string][]string{
		"Hello": BasicStringList,
		"Foo":   BasicStringList,
	},
	MapEmbeddedWithLists: map[string][][][]int64{
		"Hello": BasicDoubleEmbedded,
		"Foo":   BasicDoubleEmbedded,
	},
	All:           All_S,
	ListMapStruct: []map[int]*AllStruct{structmap, structmap},
}

func TestSerializer(t *testing.T) {
	// Create a new serializer
	s := NewSerializer()

	// Serialize the test struct
	serialized, err := s.SetCompress(true).Serialize(&testStruct)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Serialized length: ", len(string(serialized)))

	// Deserialize the test struct
	deserialized := Testie{}
	s = NewSerializer()
	s.SetData(serialized)
	err = s.SetCompress(true).Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(testStruct)
	t.Log(deserialized)

	// Check if the deserialized struct is equal to the original
	if len(deserialized.Intlist) == len(testStruct.Intlist) {
		if deserialized.Intlist[0] != testStruct.Intlist[0] {
			t.Fatal("deserialized struct is not equal to the original")
		}
		if deserialized.Intlist[1] != testStruct.Intlist[1] {
			t.Fatal("deserialized struct is not equal to the original")
		}
		if deserialized.Intlist[2] != testStruct.Intlist[2] {
			t.Fatal("deserialized struct is not equal to the original")
		}
		if deserialized.Intlist[3] != testStruct.Intlist[3] {
			t.Fatal("deserialized struct is not equal to the original")
		}
	} else {
		t.Fatal("deserialized struct is not equal to the original")
	}
	if len(deserialized.StringList) == len(testStruct.StringList) {
		if deserialized.StringList[0] != testStruct.StringList[0] {
			t.Fatal("deserialized struct is not equal to the original")
		}
		if deserialized.StringList[1] != testStruct.StringList[1] {
			t.Fatal("deserialized struct is not equal to the original")
		}
	} else {
		t.Fatal("deserialized struct is not equal to the original")
	}
	if len(deserialized.Structie.IntList) == len(testStruct.Structie.IntList) {
		if deserialized.Structie.IntList[0] != testStruct.Structie.IntList[0] {
			t.Fatal("deserialized struct is not equal to the original")
		}
		if deserialized.Structie.IntList[1] != testStruct.Structie.IntList[1] {
			t.Fatal("deserialized struct is not equal to the original")
		}
	} else {
		t.Fatal("deserialized struct is not equal to the original")
	}
	if len(deserialized.Structie.EmbeddedintList) == len(testStruct.Structie.EmbeddedintList) {
		if len(deserialized.Structie.EmbeddedintList[0]) == len(testStruct.Structie.EmbeddedintList[0]) {
			if deserialized.Structie.EmbeddedintList[0][0] != testStruct.Structie.EmbeddedintList[0][0] {
				t.Fatal("deserialized embedded int list inside of struct is not equal to the original")
			}
			if deserialized.Structie.EmbeddedintList[0][1] != testStruct.Structie.EmbeddedintList[0][1] {
				t.Fatal("deserialized embedded int list inside of struct is not equal to the original")
			}
		} else {
			t.Fatal("deserialized embedded int list inside of struct is not equal to the original")
		}
	} else {
		t.Fatal("deserialized embedded int list inside of struct is not equal to the original")
	}

	if len(deserialized.Structie.EmbeddedstrList) == len(testStruct.Structie.EmbeddedstrList) {
		if len(deserialized.Structie.EmbeddedstrList[0]) == len(testStruct.Structie.EmbeddedstrList[0]) {
			if deserialized.Structie.EmbeddedstrList[0][0] != testStruct.Structie.EmbeddedstrList[0][0] {
				t.Fatal("deserialized embedded str list inside of struct is not equal to the original")
			}
			if deserialized.Structie.EmbeddedstrList[0][1] != testStruct.Structie.EmbeddedstrList[0][1] {
				t.Fatal("deserialized embedded str list inside of struct is not equal to the original")
			}
		} else {
			t.Fatal("deserialized embedded str list inside of struct is not equal to the original")
		}
	} else {
		t.Fatal("deserialized embedded str list inside of struct is not equal to the original")
	}

	if len(deserialized.Structie.DoubleEmbedded) == len(testStruct.Structie.DoubleEmbedded) {
		if len(deserialized.Structie.DoubleEmbedded[0]) == len(testStruct.Structie.DoubleEmbedded[0]) {
			if len(deserialized.Structie.DoubleEmbedded[0][0]) == len(testStruct.Structie.DoubleEmbedded[0][0]) {
				if deserialized.Structie.DoubleEmbedded[0][0][0] != testStruct.Structie.DoubleEmbedded[0][0][0] {
					t.Fatal("deserialized double embedded int list inside of struct is not equal to the original")
				}
				if deserialized.Structie.DoubleEmbedded[0][0][1] != testStruct.Structie.DoubleEmbedded[0][0][1] {
					t.Fatal("deserialized double embedded int list inside of struct is not equal to the original")
				}
			} else {
				t.Fatal("deserialized double embedded int list inside of struct is not equal to the original")
			}
		} else {
			t.Fatal("deserialized double embedded int list inside of struct is not equal to the original")
		}
	} else {
		t.Fatal("deserialized double embedded int list inside of struct is not equal to the original")
	}
	// Validate Allstruct
	if deserialized.All != nil {
		all_s := deserialized.All
		if all_s.BoolField != testStruct.All.BoolField {
			t.Fatal("bool field is not equal to the original")
		}
		if all_s.IntField != testStruct.All.IntField {
			t.Fatal("int field is not equal to the original")
		}
		if all_s.FloatField != testStruct.All.FloatField {
			t.Fatal("float field is not equal to the original")
		}
		if all_s.StringField != testStruct.All.StringField {
			t.Fatal("string field is not equal to the original")
		}
		for i, v := range all_s.ListBool {
			if v != testStruct.All.ListBool[i] {
				t.Fatal("list bool field is not equal to the original")
			}
		}
		for i, v := range all_s.ListInt {
			if v != testStruct.All.ListInt[i] {
				t.Fatal("list int field is not equal to the original")
			}
		}

		for i, v := range all_s.ListFloat {
			if v != testStruct.All.ListFloat[i] {
				t.Fatal("list float field is not equal to the original")
			}
		}

		for i, v := range all_s.ListString {
			if v != testStruct.All.ListString[i] {
				t.Fatal("list string field is not equal to the original")
			}
		}

		for k, v := range all_s.MapBool {
			if v != testStruct.All.MapBool[k] {
				t.Fatal("map bool field is not equal to the original")
			}
		}

		for k, v := range all_s.MapInt {
			if v != testStruct.All.MapInt[k] {
				t.Fatal("map int field is not equal to the original")
			}
		}

		for k, v := range all_s.MapFloat {
			if v != testStruct.All.MapFloat[k] {
				t.Fatal("map float field is not equal to the original")
			}
		}

		for k, v := range all_s.MapString {
			if v != testStruct.All.MapString[k] {
				t.Fatal("map string field is not equal to the original")
			}
		}

		for k, v := range all_s.MapListBool {
			for i, vv := range v {
				if vv != testStruct.All.MapListBool[k][i] {
					t.Fatal("map list bool field is not equal to the original")
				}
			}
		}

		for k, v := range all_s.MapListInt {
			for i, vv := range v {
				if vv != testStruct.All.MapListInt[k][i] {
					t.Fatal("map list int field is not equal to the original")
				}
			}
		}

		for k, v := range all_s.MapListFloat {
			for i, vv := range v {
				if vv != testStruct.All.MapListFloat[k][i] {
					t.Fatal("map list float field is not equal to the original")
				}
			}
		}

		for k, v := range all_s.MapListString {
			for i, vv := range v {
				if vv != testStruct.All.MapListString[k][i] {
					t.Fatal("map list string field is not equal to the original")
				}
			}
		}

		for i, v := range all_s.ListMapBool {
			for k, vv := range v {
				if vv != testStruct.All.ListMapBool[i][k] {
					t.Fatal("list map bool field is not equal to the original")
				}
			}
		}

		for i, v := range all_s.ListMapInt {
			for k, vv := range v {
				if vv != testStruct.All.ListMapInt[i][k] {
					t.Fatal("list map int field is not equal to the original")
				}
			}
		}

		for i, v := range all_s.ListMapFloat {
			for k, vv := range v {
				if vv != testStruct.All.ListMapFloat[i][k] {
					t.Fatal("list map float field is not equal to the original")
				}
			}
		}

		for i, v := range all_s.ListMapString {
			for k, vv := range v {
				if vv != testStruct.All.ListMapString[i][k] {
					t.Fatal("list map string field is not equal to the original")
				}
			}
		}
	} else {
		t.Fatal("deserialized struct is not equal to the original")
	}
	if deserialized.ListMapStruct != nil {
		for i, v := range deserialized.ListMapStruct {
			for k, vv := range v {
				if vv.BoolField != testStruct.ListMapStruct[i][k].BoolField {
					t.Fatal("list map struct bool field is not equal to the original")
				}
				if vv.IntField != testStruct.ListMapStruct[i][k].IntField {
					t.Fatal("list map struct int field is not equal to the original")
				}
				if vv.FloatField != testStruct.ListMapStruct[i][k].FloatField {
					t.Fatal("list map struct float field is not equal to the original")
				}
				if vv.StringField != testStruct.ListMapStruct[i][k].StringField {
					t.Fatal("list map struct string field is not equal to the original")
				}
			}
		}
	} else {
		t.Fatal("deserialized struct is not equal to the original")
	}
}
