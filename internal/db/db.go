package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var d *Db

func Init() {
	if d == nil {
		d = getDb()
	}
}

func SearchSafe(id string) *Safe {
	Init()
	for _, s := range d.Safes {
		if s.ID.Address == id || s.ID.Tag == id {
			return &s
		}
	}
	return nil
}

func SearchTag(tag string) *AddressTag {
	Init()
	for _, t := range d.Tags {
		if t.Tag == tag {
			return &t
		}
	}
	return nil
}

func SwapAddress(address string) string {
	Init()
	for _, t := range d.Tags {
		if t.Address == address {
			return fmt.Sprintf("🏷️  (%s)", t.Tag)
		}
	}
	for _, s := range d.Safes {
		if s.ID.Address == address {
			return fmt.Sprintf("🏷️  (%s)", s.ID.Tag)
		}
	}
	return address
}

func AddSafe(safe Safe) error {
	Init()
	for _, s := range d.Safes {
		if s.ID.Address == safe.ID.Address {
			return fmt.Errorf(
				"safe already exists for address %s (%s)",
				safe.ID.Address,
				s.ID.Tag,
			)
		}
	}
	d.Safes = append(d.Safes, safe)
	return nil
}

func UpdateSafe(safe Safe) error {
	Init()
	for i, s := range d.Safes {
		if s.ID.Address == safe.ID.Address {
			d.Safes[i] = safe
			return nil
		}
	}
	return AddSafe(safe)
}

func AddTag(tag AddressTag) error {
	Init()
	for _, t := range d.Tags {
		if t.Address == tag.Address {
			return fmt.Errorf(
				"tag already exists for address %s (%s)",
				tag.Address,
				t.Tag,
			)
		}
	}
	d.Tags = append(d.Tags, tag)
	return nil
}

func UpdateTag(tag AddressTag) error {
	Init()
	for i, t := range d.Tags {
		if t.Address == tag.Address {
			d.Tags[i] = tag
			return nil
		}
	}
	return AddTag(tag)
}

func GetTags() []AddressTag {
	Init()
	return d.Tags
}

func GetSafes() []Safe {
	Init()
	return d.Safes
}

func Commit() error {
	Init()
	c := getDbConfig()
	err := os.MkdirAll(filepath.Dir(c.DbPath), 0700)
	if err != nil {
		return fmt.Errorf("error creating db directory: %v", err)
	}
	f, err := os.Create(c.DbPath)
	if err != nil {
		return fmt.Errorf("error creating db file: %v", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	err = enc.Encode(d)
	if err != nil {
		return fmt.Errorf("error encoding db: %v", err)
	}
	return nil
}
