package paypal

import (
	"fmt"
	"time"
)

const format = "2006-01-02T15:04:05Z"

// Filter type
type Filter struct {
	fields []fmt.Stringer
}

func (s *Filter) String() string {
	var filter string
	for i, f := range s.fields {
		if i == 0 {
			filter = "?" + f.String()
		} else {
			filter = filter + "&" + f.String()
		}
	}

	return filter
}

// TextField type
type TextField struct {
	name string
	Is   string
}

func (d TextField) String() string {
	return fmt.Sprintf("%s=%s", d.name, d.Is)
}

// TimeField type
type TimeField struct {
	name string
	Is   time.Time
}

// String .
func (d TimeField) String() string {
	return fmt.Sprintf("%s=%s", d.name, d.Is.UTC().Format(format))
}

// AddTextField .
func (s *Filter) AddTextField(field string) *TextField {
	f := &TextField{name: field}
	s.fields = append(s.fields, f)
	return f
}

// AddTimeField .
func (s *Filter) AddTimeField(field string) *TimeField {
	f := &TimeField{name: field}
	s.fields = append(s.fields, f)
	return f
}
