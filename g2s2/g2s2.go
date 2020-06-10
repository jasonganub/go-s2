package g2s2

import (
    "errors"
    "fmt"
    "github.com/golang/geo/s2"
    geojson "github.com/paulmach/go.geojson"
    "github.com/urfave/cli"
    "io/ioutil"
    "strconv"
)

func getCoordinates(fc *geojson.Feature) []s2.Point {
    var points []s2.Point
    for _, coordinate := range fc.Geometry.Polygon[0] {
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
    coverings = append(coverings, regionCoverer.Covering(s2.PolygonFromLoops(loops)))
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

func parseArgument(c *cli.Context) (fc *geojson.Feature, level int, delimiter string, err error) {
    level, err = strconv.Atoi(c.Args().Get(1))
    if err != nil {
        return nil, -1, "", err
    }

    delimiter = c.Args().Get(2)

    filePath := c.Args().First()
    if len(filePath) <= 0 {
        return nil, -1, "", errors.New("give me the file")
    }

    raw, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, -1, "", err
    }

    fc, err = geojson.UnmarshalFeature(raw)
    if err != nil {
        return nil, -1, "", err
    }

    if fc == nil {
        return nil, -1, "", errors.New("looks like bad GeoJson file. Make sure the root is a 'Feature'")
    }

    return fc, level, delimiter, nil
}

// Run as CLI entry point
func Run(c *cli.Context) error {

    const MaxCells = 100

    fc, level, delimiter, err := parseArgument(c)
    if err != nil {
        return err
    }

    points := getCoordinates(fc)
    s2IDs := getS2Ids(points, level, MaxCells)

    fmt.Printf("S2IDs Level %v :\n", level)
    //for _, s2id := range s2IDs {
    //    fmt.Printf("%v ", s2id)
    //}
    
    for i, s2id := range s2IDs {
        fmt.Print(s2id)
        if i < (len(s2IDs) -1 ){
            fmt.Print(fmt.Sprintf("%s ", delimiter))
        }
    }

    fmt.Println("")
    return nil
}
