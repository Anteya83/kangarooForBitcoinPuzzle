package state

import (
	"encoding/gob"
	"os"
)

type State struct {
	Version        int
	PublicKeyHex   string
	StartRangeHex  string
	EndRangeHex    string
	TamePosX       string
	TamePosY       string
	TameDist       string
	TotalTameSteps int64
	TotalWildSteps int64
	Map            map[string]string // key["x,y") -> distance(hex)
	FoundKey       string
}

func (s *State) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	return enc.Encode(s)
}

func Load(filename string) (*State, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var s State
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *State) IsValid(pubKeyHex, startHex, endHex string) bool {
	return s.PublicKeyHex == pubKeyHex &&
		s.StartRangeHex == startHex &&
		s.EndRangeHex == endHex
}
