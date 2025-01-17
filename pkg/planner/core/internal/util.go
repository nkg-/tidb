// Copyright 2023 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"github.com/pingcap/tidb/pkg/expression/aggregation"
	"github.com/pingcap/tidb/pkg/sessionctx"
)

// WrapCastForAggFuncs wraps the args of an aggregate function with a cast function.
// If the mode is FinalMode or Partial2Mode, we do not need to wrap cast upon the args,
// since the types of the args are already the expected.
func WrapCastForAggFuncs(sctx sessionctx.Context, aggFuncs []*aggregation.AggFuncDesc) {
	for i := range aggFuncs {
		if aggFuncs[i].Mode != aggregation.FinalMode && aggFuncs[i].Mode != aggregation.Partial2Mode {
			aggFuncs[i].WrapCastForAggArgs(sctx)
		}
	}
}
