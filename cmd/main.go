package main

import (
	"fmt"

	"github.com/tobiasjaster/HomeEnergyManagement/external/weather/dwd"
)

func main() {
	// dat, _ := os.ReadFile("/tmp/Downloads/icon-d2_germany_regular-lat-lon_single-level_2024011221_000_2d_t_2m.grib2")

	ctx := dwd.CodesContextGetDefault()
	// handle := dwd.CodesHandleNewFromMessage(ctx, dat)
	handle, err := dwd.CodesHandleNewFromFile(ctx, "/tmp/Downloads/icon-d2_germany_regular-lat-lon_single-level_2024011221_000_2d_t_2m.grib2", 1)
	fmt.Println(err)
	nearest, err := dwd.CodesGribNearestNew(handle)
	fmt.Println(err)
	nearestValues, err2 := dwd.CodesGribNearestFind(nearest, handle, 51.02554675116857, 13.622345456814264, 0)
	fmt.Println(err2)
	fmt.Println(nearestValues)
	err3 := dwd.CodesHandleDelete(handle)
	fmt.Println(err3)
}
