package stdlib

import (
	"testing"
)

// Test in a map of struct value, if the struct also contains map/slice, what would happen when the value
// is modified after it is added to the map, get the added value from map should not have the updated change
// I suppose, but due to the nature of slice, it might cause problem when append

type registry struct {
	childRegistries map[string]registry
	childLoggers    []int
}

type registryPtr struct {
	childRegistries map[string]*registryPtr
	childLoggers    []int
}

func TestMapReference(t *testing.T) {
	rTop := registry{
		childRegistries: make(map[string]registry),
	}
	r1 := registry{
		childRegistries: make(map[string]registry),
		childLoggers:    []int{0},
	}
	t.Logf("r1 registry len %d loggers len %d", len(r1.childRegistries), len(r1.childLoggers))
	rTop.childRegistries["r1"] = r1
	r1.childRegistries["r2"] = registry{}
	r1.childLoggers = append(r1.childLoggers, 1)
	t.Logf("r1 registry len %d loggers len %d", len(r1.childRegistries), len(r1.childLoggers))
	r1Cp := rTop.childRegistries["r1"]
	t.Logf("r1Cp registry len %d loggers len %d", len(r1Cp.childRegistries), len(r1Cp.childLoggers))

	// We can see map is a bit different than slice,
	//
	// for slice even if you change the underlying array by append
	// the copy of slice will still see the old range and backing array,
	//
	// for map change it in one place will will impact all the copy of that map
	// I think the underlying representation of map is pointer to a struct
	// https://github.com/golang/go/blob/master/src/runtime/map.go#L305
	// func makemap(t *maptype, hint int, h *hmap) *hmap {
	// }
	//
	//    stdlib_map_test.go:29: r1 registry len 0 loggers len 1
	//    stdlib_map_test.go:33: r1 registry len 1 loggers len 2
	//    stdlib_map_test.go:35: r1Cp registry len 1 loggers len 1
}
