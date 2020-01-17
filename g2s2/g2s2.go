package g2s2

import (
    "fmt"
    "github.com/golang/geo/s2"
    geojson "github.com/paulmach/go.geojson"
    "github.com/urfave/cli"
)

func getCoordinates(fc *geojson.FeatureCollection) []s2.Point {
    var points []s2.Point
    for _, coordinate := range fc.Features[0].Geometry.Polygon[0] {
        latLong := s2.PointFromLatLng(s2.LatLngFromDegrees(coordinate[1], coordinate[0]))
        points = append(points, latLong)
    }

    return points
}

func getChildren(cellID s2.CellID, level int) []s2.CellID {
    var children []s2.CellID

    if cellID.Level() >= level {
        return []s2.CellID{cellID.Parent(level)}
    }

    i := cellID.ChildBeginAtLevel(level)
    for i != cellID.ChildEndAtLevel(level) {
        children = append(children, i)
    }
    return children
}

func unique(intSlice []uint64) []uint64 {
    keys := make(map[uint64]bool)
    var list []uint64
    for _, entry := range intSlice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}

func getS2Ids(points []s2.Point, level int, maxCells int) []uint64 {

    regionCoverer := s2.RegionCoverer{
        MinLevel: level,
        MaxLevel: level,
        MaxCells: maxCells}

    var loops []*s2.Loop
    loops = append(loops, s2.LoopFromPoints(points))

    var coverings []s2.CellUnion
    coverings  = append(coverings, regionCoverer.Covering(s2.PolygonFromLoops(loops)))

    for _, point := range points {
        coverings = append(coverings, regionCoverer.CellUnion(point))
    }

    var coveringsUpdated []s2.CellUnion
    if level > 0 {
        for _, cells := range coverings {
            for _, cell := range cells {
                coveringsUpdated = append(coveringsUpdated, getChildren(cell, level))
            }
        }
    }

    var s2IDS []uint64
    for _, cells := range coveringsUpdated {
        for _, cell := range cells {
            s2IDS = append(s2IDS, uint64(cell))
        }
    }

    return unique(s2IDS)
}

// Run as CLI entry point
func Run(_ *cli.Context) {

    const S2IDLevel = 20
    const MaxCells = 100

    rawFeatureJSON := []byte(`{
				 "type": "FeatureCollection",
				 "features": [
				   {
					 "type": "Feature",
					 "properties": {},
					 "geometry": {
					   "type": "Polygon",
					   "coordinates": [
						 [
						   [
							  107.60782241821289,
							  -6.893148077890368
							],
							[
							  107.60822474956512,
							  -6.894474160996839
							],
							[
							  107.60981798171997,
							  -6.8941492976071075
							],
							[
							  107.60939955711365,
							  -6.892748654540805
							],
							[
							  107.60782241821289,
							  -6.893148077890368
							]
						 ]
					   ]
					 }
				   }
				 ]
				}`)

    fc, err := geojson.UnmarshalFeatureCollection(rawFeatureJSON)
    if err != nil {
        fmt.Printf("Bad GEOJSON file: %v", err)
    }

    points := getCoordinates(fc)
    s2IDs := getS2Ids(points, S2IDLevel, MaxCells)

    fmt.Printf("\n\nS2IDs \n%v \n", s2IDs)

}
