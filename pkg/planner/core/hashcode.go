// Copyright 2019 PingCAP, Inc.
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

package core

import (
	"bytes"
	"encoding/binary"
	"slices"

	"github.com/pingcap/tidb/pkg/util/plancodec"
)

func encodeIntAsUint32(result []byte, value int) []byte {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(value))
	return append(result, buf[:]...)
}

// HashCode implements LogicalPlan interface.
func (p *baseLogicalPlan) HashCode() []byte {
	// We use PlanID for the default hash, so if two plans do not have
	// the same id, the hash value will never be the same.
	result := make([]byte, 0, 4)
	result = encodeIntAsUint32(result, p.ID())
	return result
}

// HashCode implements LogicalPlan interface.
func (p *LogicalProjection) HashCode() []byte {
	// PlanType + SelectOffset + ExprNum + [Exprs]
	// Expressions are commonly `Column`s, whose hashcode has the length 9, so
	// we pre-alloc 10 bytes for each expr's hashcode.
	result := make([]byte, 0, 12+len(p.Exprs)*10)
	result = encodeIntAsUint32(result, plancodec.TypeStringToPhysicalID(p.TP()))
	result = encodeIntAsUint32(result, p.SelectBlockOffset())
	result = encodeIntAsUint32(result, len(p.Exprs))
	for _, expr := range p.Exprs {
		exprHashCode := expr.HashCode()
		result = encodeIntAsUint32(result, len(exprHashCode))
		result = append(result, exprHashCode...)
	}
	return result
}

// HashCode implements LogicalPlan interface.
func (p *LogicalTableDual) HashCode() []byte {
	// PlanType + SelectOffset + RowCount
	result := make([]byte, 0, 12)
	result = encodeIntAsUint32(result, plancodec.TypeStringToPhysicalID(p.TP()))
	result = encodeIntAsUint32(result, p.SelectBlockOffset())
	result = encodeIntAsUint32(result, p.RowCount)
	return result
}

// HashCode implements LogicalPlan interface.
func (p *LogicalSelection) HashCode() []byte {
	// PlanType + SelectOffset + ConditionNum + [Conditions]
	// Conditions are commonly `ScalarFunction`s, whose hashcode usually has a
	// length larger than 20, so we pre-alloc 25 bytes for each expr's hashcode.
	result := make([]byte, 0, 12+len(p.Conditions)*25)
	result = encodeIntAsUint32(result, plancodec.TypeStringToPhysicalID(p.TP()))
	result = encodeIntAsUint32(result, p.SelectBlockOffset())
	result = encodeIntAsUint32(result, len(p.Conditions))

	condHashCodes := make([][]byte, len(p.Conditions))
	for i, expr := range p.Conditions {
		condHashCodes[i] = expr.HashCode()
	}
	// Sort the conditions, so `a > 1 and a < 100` can equal to `a < 100 and a > 1`.
	slices.SortFunc(condHashCodes, func(i, j []byte) int { return bytes.Compare(i, j) })

	for _, condHashCode := range condHashCodes {
		result = encodeIntAsUint32(result, len(condHashCode))
		result = append(result, condHashCode...)
	}
	return result
}

// HashCode implements LogicalPlan interface.
func (p *LogicalLimit) HashCode() []byte {
	// PlanType + SelectOffset + Offset + Count
	result := make([]byte, 24)
	binary.BigEndian.PutUint32(result, uint32(plancodec.TypeStringToPhysicalID(p.TP())))
	binary.BigEndian.PutUint32(result[4:], uint32(p.SelectBlockOffset()))
	binary.BigEndian.PutUint64(result[8:], p.Offset)
	binary.BigEndian.PutUint64(result[16:], p.Count)
	return result
}
