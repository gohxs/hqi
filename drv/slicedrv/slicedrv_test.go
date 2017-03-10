package slicedrv

import (
	"log"

	"github.com/gohxs/hqi"

	"testing"
)

type Device struct {
	Brand string `json:",omitempty"`
	Model string `json:",omitempty"`
}

var (
	DeviceData = []Device{
		Device{"Apple", "IPhone 5"},
		Device{"Samsung", "Note 5"},
		Device{"Apple", "IPhone 4"},
		Device{"HTC", "One"},
	}
)

func TestColl(t *testing.T) {
	q := hqi.NewQuery(&Driver{&DeviceData})
	q.Insert(Device{"Sony", "Xperia"})
}

func TestRemoveOr(t *testing.T) {
	q := hqi.NewQuery(&Driver{&DeviceData})
	q.Find(hqi.M{"Brand": "Apple"}, hqi.M{"Brand": "HTC"}).Delete()
	log.Println("DeviceData:", DeviceData)
}
