// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package parse_vehicle_trace2trace

// Config for parse_parse_vehicle_trace2trace processor.
type Config struct {
	Field           string `config:"field"`            // log message field
	TimeField       string `config:"time_field"`       // specified the time field
	IgnoreMissing   bool   `config:"ignore_missing"`   // Skip field when From field is missing.
	IgnoreMalformed bool   `config:"ignore_malformed"` // Skip log when From log is incorrect.
	DropOrigin      bool   `config:"drop_origin"`      // Drop the origin field
}

func defaultConfig() Config {
	return Config{
		Field:           "message",
		TimeField:       "logtime",
		IgnoreMissing:   true,
		IgnoreMalformed: true,
		DropOrigin:      true,
	}
}
