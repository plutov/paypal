package paypal

import (
	"testing"
	"time"
)

func TestFilter_AddTextField(t *testing.T) {
	filter := &Filter{}
	filter.AddTextField("sort_by").Is = "create_time"
	filter.AddTextField("count").Is = "30"
	filter.AddTextField("sort_order").Is = "desc"

	expected := "?sort_by=create_time&count=30&sort_order=desc"
	if filter.String() != expected {
		t.Errorf("filter string was %s, wanted %s", filter.String(), expected)
	}
}

func TestFilter_AddTimeField(t *testing.T) {
	filter := &Filter{}
	startTime := time.Time{}
	endTime := startTime.Add(time.Hour * 24 * 30)
	filter.AddTimeField("start_time").Is = startTime
	filter.AddTimeField("stop_time").Is = endTime

	expected := "?start_time=0001-01-01T00:00:00Z&stop_time=0001-01-31T00:00:00Z"
	if filter.String() != expected {
		t.Errorf("filter string was %s, wanted %s", filter.String(), expected)
	}
}

func TestFilter_AddMixedFields(t *testing.T) {
	filter := &Filter{}
	startTime := time.Time{}
	endTime := startTime.Add(time.Hour * 24 * 30)
	filter.AddTimeField("stop_time").Is = endTime
	filter.AddTextField("count").Is = "30"

	expected := "?stop_time=0001-01-31T00:00:00Z&count=30"
	if filter.String() != expected {
		t.Errorf("filter string was %s, wanted %s", filter.String(), expected)
	}
}
