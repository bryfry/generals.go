package main

import (
	"fmt"
	"github.com/fatih/color"
)

type Map struct {
	Scores      []PlayerScore
	Turn        int
	AttackIndex int
	Generals    []int
	MapArray    []int
	CitiesArray []int
	Width       int
	Height      int
	Size        int
	Tiles       [][]Tile
}

type Tile struct {
	Owner   int // Player Index
	Armies  int
	Terrain TerrainType // Terrain Type
	Seen    bool
	Fog     bool
}

type TerrainType int

const (
	EMPTY TerrainType = iota + 1
	MOUNTAIN
	OBSTACLE
	CITY
	CAPITAL // game calls this a GENERAL
)

func (m *Map) Print() {
	fmt.Println("TILES", m.Height, m.Width, m.Size)
	for i := 0; i < m.Height; i++ {
		fmt.Println()
		for j := 0; j < m.Width; j++ {
			m.Tiles[i][j].Print()
		}
	}
	fmt.Println()
}

func (m *Map) Patch(u Update) {
	m.MapArray = patch(m.MapArray, u.MapDiff)
	m.CitiesArray = patch(m.CitiesArray, u.CitiesDiff)
	m.Generals = u.Generals
	// Init
	if m.Tiles == nil {
		m.Width = m.MapArray[0]
		m.Height = m.MapArray[1]
		m.Tiles = make([][]Tile, m.Height)
		m.Size = m.Width * m.Height
		for i := 0; i < m.Height; i++ {
			m.Tiles[i] = make([]Tile, m.Width)
		}
	}
	// Apply patch to Map/Tiles
	for i := 0; i < m.Height; i++ {
		for j := 0; j < m.Width; j++ {
			m.Tiles[i][j].Armies = m.MapArray[(i*m.Width)+j+2]
			m.Tiles[i][j].DecodeTerrain(m.MapArray[(i*m.Width)+j+2+m.Size])
			for k := range m.Generals {
				if i*m.Width+j == m.Generals[k] {
					m.Tiles[i][j].Owner = k
					m.Tiles[i][j].Terrain = CAPITAL
				}
			}
			for l := range m.CitiesArray {
				if i*m.Width+j == m.CitiesArray[l] {
					m.Tiles[i][j].Terrain = CITY
				}
			}
		}
	}
}

func patch(old []int, diff []int) (new []int) {
	i := 0
	for i < len(diff) {
		if diff[i] > 0 {
			new = append(new, old[len(new):(len(new)+diff[i])]...)
		}
		i++
		if i < len(diff) && diff[i] > 0 {
			new = append(new, diff[i+1:(i+1+diff[i])]...)
			i += diff[i]
		}
		i++
	}
	return new
}

func (t *Tile) DecodeTerrain(terrain int) {
	switch terrain {
	case GENIO_EMPTY:
		t.Seen = true
		t.Terrain = EMPTY
		t.Owner = 0
	case GENIO_MOUNTAIN:
		t.Seen = true
		t.Terrain = MOUNTAIN
	case GENIO_FOG:
		t.Fog = true
		t.Terrain = EMPTY
	case GENIO_FOG_OBSTACLE:
		// Retain knowledge.
		// Obstacle can become a Mountain or City,
		// not the other way around
		if t.Terrain == 0 {
			t.Terrain = OBSTACLE
		}
		t.Fog = true
	default:
		t.Seen = true
		t.Owner = terrain
	}
}

func (t *Tile) Print() {
	p := color.New(color.FgWhite)
	if t.Fog {
		if t.Seen {
			p.Add(color.BgHiBlue)
		} else {
			p.Add(color.BgHiBlue)
		}
	}
	switch t.Terrain {
	case EMPTY:
		if t.Fog || t.Armies == 0 {
			p.Printf("  _")
		} else {
			// TODO if owner != PlayerID of bot = ReD
			p.Printf("%3d", t.Armies)
		}
	case MOUNTAIN:
		p.Printf("  M")
	case OBSTACLE:
		p.Printf("  ?")
	case CITY:
		p.Add(color.BgCyan)
		p.Printf("%3d", t.Armies)
	case CAPITAL:
		p.Add(color.Bold, color.BgHiMagenta)
		p.Printf("%3d", t.Armies)
	default:
		p.Add(color.FgRed)
		p.Printf("%3d", t.Armies)
	}

}
