package date

// Copyright 2017 Microsoft Corporation
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// UnixTime provides marshalling and unmarshalling facilities of time.Time into or from Unix time format.
// The format is the number of seconds elapsed from the Unix epoch.
type UnixTime time.Time

// Duration returns the time.Time as a Duration since the Unix epoch.
func (t UnixTime) Duration() time.Duration {
	return time.Duration(time.Time(t).Unix())
}

// NewUnixTimeFromSeconds creates a UnixTime as a number of seconds from the UnixEpoch.
// Deprecated: please use Unix time facilities in the standard library.
func NewUnixTimeFromSeconds(seconds float64) UnixTime {
	return NewUnixTimeFromDuration(time.Duration(seconds * float64(time.Second)))
}

// NewUnixTimeFromNanoseconds creates a UnixTime as a number of nanoseconds from the UnixEpoch.
// Deprecated: please use Unix time facilities in the standard library.
func NewUnixTimeFromNanoseconds(nanoseconds int64) UnixTime {
	return NewUnixTimeFromDuration(time.Duration(nanoseconds))
}

// NewUnixTimeFromDuration creates a UnixTime as a duration of time since the UnixEpoch.
// Deprecated: please use Unix time facilities in the standard library.
func NewUnixTimeFromDuration(dur time.Duration) UnixTime {
	return UnixTime(UnixEpoch().Add(dur))
}

// UnixEpoch retreives the moment considered the Unix Epoch. I.e. The time represented by '0'
func UnixEpoch() time.Time {
	return time.Unix(0, 0)
}

// MarshalJSON preserves the UnixTime as a JSON number conforming to Unix Timestamp requirements.
// (i.e. the number of seconds since midnight January 1st, 1970 not considering leap seconds.)
func (t UnixTime) MarshalJSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	enc := json.NewEncoder(buffer)
	err := enc.Encode(time.Time(t).Unix())
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// UnmarshalJSON reconstitures a UnixTime saved as a JSON number of the number of seconds since
// midnight January 1st, 1970.
func (t *UnixTime) UnmarshalJSON(text []byte) error {
	dec := json.NewDecoder(bytes.NewReader(text))

	var secondsSinceEpoch int64
	if err := dec.Decode(&secondsSinceEpoch); err != nil {
		return err
	}

	*t = UnixTime(time.Unix(secondsSinceEpoch, 0))

	return nil
}

// MarshalText returns the number of seconds since the Unix Epoch as an int64 in text format.
func (t UnixTime) MarshalText() ([]byte, error) {
	s := fmt.Sprintf("%d", time.Time(t).Unix())
	return []byte(s), nil
}

// UnmarshalText populates a UnixTime with a value stored textually as a floating point number of seconds since the Unix Epoch.
func (t *UnixTime) UnmarshalText(raw []byte) error {
	ut, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil {
		return err
	}
	*t = UnixTime(time.Unix(ut, 0))
	return nil
}

// MarshalBinary converts a UnixTime into a binary.LittleEndian int64 of seconds since the epoch.
func (t UnixTime) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}

	payload := int64(t.Duration())

	if err := binary.Write(buf, binary.LittleEndian, &payload); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary converts a from a binary.LittleEndian int64 of seconds since the epoch into a UnixTime.
func (t *UnixTime) UnmarshalBinary(raw []byte) error {
	var nanosecondsSinceEpoch int64

	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, &nanosecondsSinceEpoch); err != nil {
		return err
	}
	*t = NewUnixTimeFromNanoseconds(nanosecondsSinceEpoch)
	return nil
}
