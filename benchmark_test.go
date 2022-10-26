package tinyserializer

import (
	"testing"
	"time"
)

type A struct {
	Name     string    `tiny:"name"`
	BirthDay time.Time `tiny:"birthday"`
	Phone    string    `tiny:"phone"`
	Siblings int       `tiny:"siblings"`
	Spouse   bool      `tiny:"spouse"`
	Money    float64   `tiny:"money"`
}

var ASTRUCT = GetA()

func GetA() A {
	return A{
		Name:     "John Doe",
		BirthDay: time.Now(),
		Phone:    "123456789",
		Siblings: 2,
		Spouse:   true,
		Money:    123.45,
	}
}

func BenchmarkSerializer(b *testing.B) {
	ser := NewSerializer()
	for i := 0; i < b.N; i++ {
		data, err := ser.Serialize(&ASTRUCT)
		if err != nil {
			b.Error(err)
		}
		_ = data

	}
}

func BenchmarkDeserializer(b *testing.B) {
	var newt_struct A = ASTRUCT
	ser := NewSerializer()
	data, err := ser.Serialize(&newt_struct)
	if err != nil {
		b.Error(err)
	}
	for i := 0; i < b.N; i++ {
		err := ser.Deserialize(data, &newt_struct)
		if err != nil {
			b.Error(err)
		}
	}
}
