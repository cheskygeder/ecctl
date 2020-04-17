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

package plan

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/pkg/util"
)

var errNull = errors.New(`{
  "errors": null
}`)

func TestCancel(t *testing.T) {
	type args struct {
		params CancelParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Succeeds",
			args: args{params: CancelParams{ClusterParams: util.ClusterParams{
				ClusterID: "8b43ed5e277f7ea6f13606fcf4027f9c",
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Body:       mock.NewStringBody(`{}`),
				}}),
			}}},
		},
		{
			name: "Fails due API error",
			args: args{params: CancelParams{ClusterParams: util.ClusterParams{
				ClusterID: "8b43ed5e277f7ea6f13606fcf4027f9c",
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
					Body:       mock.NewStringBody(`{}`),
				}}),
			}}},
			err: errNull,
		},
		{
			name: "Fails due to missing API parameter validation",
			args: args{params: CancelParams{ClusterParams: util.ClusterParams{
				ClusterID: "8b43ed5e277f7ea6f13606fcf4027f9c",
			}}},
			err: util.ErrAPIReq,
		},
		{
			name: "Fails due to missing cluster ID parameter validation",
			args: args{params: CancelParams{ClusterParams: util.ClusterParams{
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
					Body:       mock.NewStringBody(`{}`),
				}}),
			}}},
			err: errors.New("cluster id should have a length of 32 characters"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Cancel(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Cancel() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}