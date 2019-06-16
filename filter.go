package paypalsdk

import (
	"fmt"
	"time"
)

const format = "2006-01-02T15:04:05Z"


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

type TextField struct {
	name string
	Is   string
}

func (d TextField) String() string {
	return fmt.Sprintf("%s=%s", d.name, d.Is)
}

type TimeField struct {
	name string
	Is   time.Time
}

func (d TimeField) String() string {
	return fmt.Sprintf("%s=%s", d.name, d.Is.UTC().Format(format))
}

func (s *Filter) AddTextField(field string) *TextField {
	f := &TextField{name: field}
	s.fields = append(s.fields, f)
	return f
}

func (s *Filter) AddTimeField(field string) *TimeField {
	f := &TimeField{name: field}
	s.fields = append(s.fields, f)
	return f
}
