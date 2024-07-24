package persistentinmemory

import (
	"encoding/json"
	"io"

	"github.com/red-serenity/goptuna"
)

// NewStorage returns new Persistent InMemory storage.
func NewStorage(r io.Reader) (*Storage, error) {
	s := &Storage{}

	// If reader is nil, we don't load anything
	if r == nil {
		return s, nil
	}

	// Otherwise we try to reload data
	dec := json.NewDecoder(r)
	err := dec.Decode(s)
	if err != nil {
		return nil, err
	}

	if err := s.SetStudyDirection(s.StudyID, s.Direction); err != nil {
		return nil, err
	}

	for _, trial := range s.Trials {
		if _, err := s.CloneTrial(s.StudyID, trial); err != nil {
			return nil, err
		}
	}

	for k, v := range s.UserAttrs {
		if err := s.SetStudyUserAttr(s.StudyID, k, v); err != nil {
			return nil, err
		}
	}

	for k, v := range s.SystemAttrs {
		if err := s.SetStudySystemAttr(s.StudyID, k, v); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Storage stores data in your relational databases.
type Storage struct {
	goptuna.InMemoryStorage

	Direction   goptuna.StudyDirection
	Trials      []goptuna.FrozenTrial
	UserAttrs   map[string]string
	SystemAttrs map[string]string
	StudyName   string
	StudyID     int
}

// Allow to store StudyName
func (s *Storage) CreateNewStudy(name string) (int, error) {
	var err error
	s.StudyName = name
	s.StudyID, err = s.InMemoryStorage.CreateNewStudy(name)
	return s.StudyID, err
}

func (s *Storage) Dump(w io.Writer) error {
	var err error

	if s.StudyID, err = s.GetStudyIDFromName(s.StudyName); err != nil {
		return err
	}

	if s.Direction, err = s.GetStudyDirection(s.StudyID); err != nil {
		return err
	}

	s.Trials = make([]goptuna.FrozenTrial, 0)
	trials, err := s.GetAllTrials(s.StudyID)
	if err != nil {
		return err
	}
	s.Trials = append(s.Trials, trials...)

	s.UserAttrs = make(map[string]string)
	userAttrs, err := s.GetStudyUserAttrs(s.StudyID)
	if err != nil {
		return err
	}
	for k, v := range userAttrs {
		s.UserAttrs[k] = v
	}

	s.SystemAttrs = make(map[string]string)
	sysAttrs, err := s.GetStudySystemAttrs(s.StudyID)
	if err != nil {
		return err
	}
	for k, v := range sysAttrs {
		s.SystemAttrs[k] = v
	}

	enc := json.NewEncoder(w)
	return enc.Encode(s)
}
