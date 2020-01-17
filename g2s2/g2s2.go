package g2s2

import (
    "fmt"
    "github.com/golang/geo/s2"
    geojson "github.com/paulmach/go.geojson"
    "github.com/urfave/cli"
)

func getCoordinatesAsPoints(fc *geojson.FeatureCollection) []s2.Point {
    var points []s2.Point

    for _, coordinate := range fc.Features[0].Geometry.Polygon[0] {
        fmt.Println(coordinate)
        latLong := s2.PointFromLatLng(s2.LatLngFromDegrees(coordinate[1], coordinate[0]))
        points = append(points, latLong)
    }

    fmt.Println("points")
    fmt.Println(points)

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

func getS2Ids(points []s2.Point, level int) []uint64 {
    maxCells := 100
    regionCoverer := s2.RegionCoverer{MinLevel: level, MaxLevel: level, MaxCells: maxCells}
    var coverings []s2.CellUnion

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

func getFeature(s int) *geojson.Feature {
    fmt.Println("\ns is ")
    fmt.Println(s)

    c := s2.CellFromCellID(s2.CellID(s))

    fmt.Println("\ncell is ")
    fmt.Println(c)

    var y []s2.LatLng
    for _, i := range []int{0, 1, 2, 3, 0} {
        vertex := c.Vertex(i)
        y = append(y, s2.LatLngFromPoint(vertex))
    }

    fmt.Println("\ny is ")
    fmt.Println(y)

    var x [][][]float64
    for _, ll := range y {
        lat := []float64{ll.Lat.Degrees()}
        lng := []float64{ll.Lng.Degrees()}
        row := [][]float64{lat, lng}
        x = append(x, row)
    }

    fmt.Println("\nx is")
    fmt.Println(x)

    geometry := &geojson.Geometry{Polygon: x}

    fmt.Println("\ngeometry is")
    fmt.Println(geometry)

    fmt.Println("string s and string c level")
    fmt.Println(s, c.Level())
    fmt.Println(string(s))
    fmt.Println(string(c.Level()))

    properties := map[string]interface{}{
        "s2id": s,
        "lvl":  c.Level(),
    }

    fmt.Println("\nproperties")
    fmt.Println(properties)

    f := geojson.Feature{Geometry: geometry, Properties: properties}

    fmt.Println("\nfeatures")
    fmt.Println(f)

    return nil
}

// Run as CLI entry point
func Run(c *cli.Context)  {

    const S2IDLevel = 15

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
        fmt.Printf("error from unmarshalling: %v", err)
    }

    points := getCoordinatesAsPoints(fc)

    s2IDs := getS2Ids(points, S2IDLevel)
    fmt.Println("\ns2 ids")
    fmt.Println(s2IDs)

}
