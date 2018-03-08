package drum

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"strconv"
)

////  Drum machine File format
// 14: SPLICE + 8 byte length of rest of file
// 32: byte string HW Version
// 4: float Tempo
// 1 byte str/int track id
// 4 byte int length of track name
// Track name
// 16 bytes 0/1 based on track

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {
	pattern := &Pattern{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) < 6 || string(data[:6]) != "SPLICE" {
		return nil, fmt.Errorf("File '%s' is not a valid splice file, SPLICE header not found", path)
	}
	r := bytes.NewReader(data[6:]) // Skip SPLICE header

	// remaining file data size
	var dataSize int64
	binary.Read(r, binary.BigEndian, &dataSize)

	// Read version and Tempo
	var header struct {
		Version [32]byte
		Tempo   float32
	}
	binary.Read(r, binary.LittleEndian, &header)
	// get where version string ends
	pos := bytes.Index(header.Version[:], []byte("\x00"))
	pattern.Version = string(header.Version[:pos])
	pattern.Tempo = header.Tempo

	// keep remaining file bytes
	rest := dataSize - (32 + 8) // already read version and Tempo

	for rest > 0 { // read tracks
		var tr track

		// read track id
		var id uint8
		binary.Read(r, binary.BigEndian, &id)
		tr.ID = id

		// read track name
		var nameSize uint32
		binary.Read(r, binary.BigEndian, &nameSize)
		var name []byte
		for i := 0; i < int(nameSize); i++ {
			var c byte
			binary.Read(r, binary.BigEndian, &c)
			name = append(name, c)
		}
		tr.Name = string(name)

		// read track steps
		var steps [16]byte
		binary.Read(r, binary.BigEndian, &steps)
		tr.Steps = steps

		pattern.Tracks = append(pattern.Tracks, tr)
		rest -= int64(1 + 4 + nameSize + 16)
	}

	return pattern, nil
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	Version string
	Tempo   float32
	Tracks  []track
}

type track struct {
	ID    uint8
	Name  string
	Steps [16]byte
}

func (p Pattern) String() string {
	tempoMinPres := strconv.FormatFloat(float64(p.Tempo), 'f', -1, 32)
	header := fmt.Sprintf("Saved with HW Version: %s\nTempo: %s\n", p.Version, tempoMinPres)
	s := header
	for _, track := range p.Tracks {
		s += fmt.Sprintf("%s", track)
	}
	return s
}

func (t track) String() string {
	s := fmt.Sprintf("(%d) %s\t|", t.ID, t.Name)
	for m := 0; m < 4; m++ {
		for i := 0; i < 4; i++ {
			if t.Steps[m*4+i] != '\x00' {
				s += "x"
			} else {
				s += "-"
			}
		}
		s += "|"
	}
	s += "\n"
	return s
}
