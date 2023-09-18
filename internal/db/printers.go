package db

import "fmt"

func (s *Safe) String() string {
	return fmt.Sprintf(
		"  {%s}\n  Address: %s\n  Threshold: %d\n  Owners: %v\n",
		s.ID.Tag,
		s.ID.Address,
		s.Threshold,
		s.Owners,
	)
}

func (t *AddressTag) String() string {
	return fmt.Sprintf(
		"  {%s}: %s\n",
		t.Tag,
		t.Address,
	)
}
