package storage

import (
	"reflect"
	"testing"
	"time"
)

// TestNew tests creating new storage.
func TestNew(t *testing.T) {
	tests := []struct {
		Timeout       time.Duration
		ExpectedError bool
	}{
		{1 * time.Second, false},
		{100 * time.Minute, false},
		{-10 * time.Hour, true},
		{time.Duration(0), true},
	}

	for _, test := range tests {
		_, err := New(test.Timeout)
		if (err != nil) != test.ExpectedError {
			t.Errorf("Timeout: %v. Expected error: %t. Got: %v\n", test.Timeout, test.ExpectedError, err)
		}
	}
}

// TestSetGet tests setting and getting values with storage.
func TestSetGet(t *testing.T) {
	tests := []struct {
		Key           string
		Value         interface{}
		Timeout       time.Duration
		WaitTime      time.Duration // WaitTime is time before getting value by key.
		GetAgain      bool          // GetAgain gets value again (in order to check if value deleting after being got)
		ExpectedValue interface{}
		ExpectedError bool
	}{
		{"Hello", "World", 1 * time.Second, 100 * time.Millisecond, false, "World", false},
		{"Ages", []int{1, 2, 3}, 1 * time.Second, 100 * time.Millisecond, false, []int{1, 2, 3}, false},
		{"Amount", 1000, 100 * time.Millisecond, 1 * time.Second, false, nil, true},
		{"ProductID", "12AB8", 1 * time.Second, 100 * time.Millisecond, true, nil, true},
	}

	for _, test := range tests {
		s, _ := New(test.Timeout)
		s.Set(test.Key, test.Value)
		<-time.After(test.WaitTime)
		value, err := s.Get(test.Key)
		if test.GetAgain {
			<-time.After(test.WaitTime)
			value, err = s.Get(test.Key)
		}
		if ((err != nil) != test.ExpectedError) || !reflect.DeepEqual(value, test.ExpectedValue) {
			t.Errorf("Timeout: %v. WaitTime: %v. Expected value: %v. Got value: %v. Get again: %v. Expected error: %t. Got: %v\n", test.Timeout, test.WaitTime, test.ExpectedValue, value, test.GetAgain, test.ExpectedError, err)
		}
	}
}
