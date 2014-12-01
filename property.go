package devatlas

import "encoding/json"

const trueValue uint32 = 1

// A property is a node attribute
type Property struct {
	Kind uint8
	Name string
}

// Convert converts a value to the `kind`interface
func (p *Property) Convert(v interface{}) interface{} {
	switch p.Kind {
	case KIND_BOOL:
		return v == trueValue
	}
	return v
}

// UnmarshalJSON custom implementation
func (p *Property) UnmarshalJSON(b []byte) error {
	var name string
	if err := json.Unmarshal(b, &name); err != nil {
		return err
	}

	kind := KIND_STRING
	switch name[0] {
	case 'b':
		kind = KIND_BOOL
	case 'i':
		kind = KIND_INT
	}
	*p = Property{kind, name[1:]}
	return nil
}
