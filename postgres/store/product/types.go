package product

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type ImageUrls map[int]string

func (i *ImageUrls) Scan(src any) error {
	switch s := src.(type) {
	case []byte:
		if err := json.Unmarshal(s, i); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(s), i); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported scan type for ImageUrls: %T", src)
	}

	return nil
}

func (i ImageUrls) Value() (driver.Value, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return b, nil
}

type Specs map[string]string

func (i *Specs) Scan(src any) error {
	switch s := src.(type) {
	case []byte:
		if err := json.Unmarshal(s, i); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(s), i); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported scan type for Specs: %T", src)
	}

	return nil
}

func (i Specs) Value() (driver.Value, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return b, nil
}
