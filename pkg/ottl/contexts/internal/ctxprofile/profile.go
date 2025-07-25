// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ctxprofile // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/internal/ctxprofile"

import (
	"context"
	"encoding/hex"
	"errors"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pprofile"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/internal/ctxcommon"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/internal/ctxerror"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/internal/ctxutil"
)

func PathGetSetter[K Context](path ottl.Path[K]) (ottl.GetSetter[K], error) {
	if path == nil {
		return nil, ctxerror.New("nil", "nil", Name, DocRef)
	}
	switch path.Name() {
	case "sample_type":
		return accessSampleType[K](), nil
	case "sample":
		return accessSample[K](), nil
	case "location_indices":
		return accessLocationIndices[K](), nil
	case "time_unix_nano":
		return accessTimeUnixNano[K](), nil
	case "time":
		return accessTime[K](), nil
	case "duration_unix_nano":
		return accessDurationUnixNano[K](), nil
	case "duration":
		return accessDuration[K](), nil
	case "period_type":
		return accessPeriodType[K](), nil
	case "period":
		return accessPeriod[K](), nil
	case "comment_string_indices":
		return accessCommentStringIndices[K](), nil
	case "default_sample_type_index":
		return accessDefaultSampleTypeIndex[K](), nil
	case "profile_id":
		nextPath := path.Next()
		if nextPath != nil {
			if nextPath.Name() == "string" {
				return accessStringProfileID[K](), nil
			}
			return nil, ctxerror.New(nextPath.Name(), nextPath.String(), Name, DocRef)
		}
		return accessProfileID[K](), nil
	case "attribute_indices":
		return accessAttributeIndices[K](), nil
	case "dropped_attributes_count":
		return accessDroppedAttributesCount[K](), nil
	case "original_payload_format":
		return accessOriginalPayloadFormat[K](), nil
	case "original_payload":
		return accessOriginalPayload[K](), nil
	case "attributes":
		if path.Keys() == nil {
			return accessAttributes[K](), nil
		}
		return accessAttributesKey(path.Keys()), nil
	default:
		return nil, ctxerror.New(path.Name(), path.String(), Name, DocRef)
	}
}

func accessSampleType[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return tCtx.GetProfile().SampleType(), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			if v, ok := val.(pprofile.ValueTypeSlice); ok {
				v.CopyTo(tCtx.GetProfile().SampleType())
			}
			return nil
		},
	}
}

func accessSample[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return tCtx.GetProfile().Sample(), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			if v, ok := val.(pprofile.SampleSlice); ok {
				v.CopyTo(tCtx.GetProfile().Sample())
			}
			return nil
		},
	}
}

func accessLocationIndices[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return ctxutil.GetCommonIntSliceValues[int32](tCtx.GetProfile().LocationIndices()), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			return ctxutil.SetCommonIntSliceValues[int32](tCtx.GetProfile().LocationIndices(), val)
		},
	}
}

func accessTimeUnixNano[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return tCtx.GetProfile().Time().AsTime().UnixNano(), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			if i, ok := val.(int64); ok {
				tCtx.GetProfile().SetTime(pcommon.NewTimestampFromTime(time.Unix(0, i)))
			}
			return nil
		},
	}
}

func accessTime[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return tCtx.GetProfile().Time().AsTime(), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			if i, ok := val.(time.Time); ok {
				tCtx.GetProfile().SetTime(pcommon.NewTimestampFromTime(i))
			}
			return nil
		},
	}
}

func accessDurationUnixNano[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return tCtx.GetProfile().Duration().AsTime().UnixNano(), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			if t, ok := val.(int64); ok {
				tCtx.GetProfile().SetDuration(pcommon.NewTimestampFromTime(time.Unix(0, t)))
			}
			return nil
		},
	}
}

func accessDuration[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return tCtx.GetProfile().Duration().AsTime(), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			if t, ok := val.(time.Time); ok {
				tCtx.GetProfile().SetDuration(pcommon.NewTimestampFromTime(t))
			}
			return nil
		},
	}
}

func accessPeriodType[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return tCtx.GetProfile().PeriodType(), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			if v, ok := val.(pprofile.ValueType); ok {
				v.CopyTo(tCtx.GetProfile().PeriodType())
			}
			return nil
		},
	}
}

func accessPeriod[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return tCtx.GetProfile().Period(), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			if i, ok := val.(int64); ok {
				tCtx.GetProfile().SetPeriod(i)
			}
			return nil
		},
	}
}

func accessCommentStringIndices[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return ctxutil.GetCommonIntSliceValues[int32](tCtx.GetProfile().CommentStrindices()), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			return ctxutil.SetCommonIntSliceValues[int32](tCtx.GetProfile().CommentStrindices(), val)
		},
	}
}

func accessDefaultSampleTypeIndex[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return int64(tCtx.GetProfile().DefaultSampleTypeIndex()), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			if i, ok := val.(int64); ok {
				tCtx.GetProfile().SetDefaultSampleTypeIndex(int32(i))
			}
			return nil
		},
	}
}

func accessProfileID[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return tCtx.GetProfile().ProfileID(), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			if id, ok := val.(pprofile.ProfileID); ok {
				if id.IsEmpty() {
					return errors.New("profile ids must not be empty")
				}
				tCtx.GetProfile().SetProfileID(id)
			}
			return nil
		},
	}
}

func accessStringProfileID[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			id := tCtx.GetProfile().ProfileID()
			return hex.EncodeToString(id[:]), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			if s, ok := val.(string); ok {
				id, err := ctxcommon.ParseProfileID(s)
				if err != nil {
					return err
				}
				if id.IsEmpty() {
					return errors.New("profile ids must not be empty")
				}
				tCtx.GetProfile().SetProfileID(id)
			}
			return nil
		},
	}
}

func accessAttributeIndices[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return ctxutil.GetCommonIntSliceValues[int32](tCtx.GetProfile().AttributeIndices()), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			return ctxutil.SetCommonIntSliceValues[int32](tCtx.GetProfile().AttributeIndices(), val)
		},
	}
}

func accessDroppedAttributesCount[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return int64(tCtx.GetProfile().DroppedAttributesCount()), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			if i, ok := val.(int64); ok {
				tCtx.GetProfile().SetDroppedAttributesCount(uint32(i))
			}
			return nil
		},
	}
}

func accessOriginalPayloadFormat[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return tCtx.GetProfile().OriginalPayloadFormat(), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			if v, ok := val.(string); ok {
				tCtx.GetProfile().SetOriginalPayloadFormat(v)
			}
			return nil
		},
	}
}

func accessOriginalPayload[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return tCtx.GetProfile().OriginalPayload().AsRaw(), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			if v, ok := val.([]byte); ok {
				tCtx.GetProfile().OriginalPayload().FromRaw(v)
			}
			return nil
		},
	}
}

func accessAttributes[K Context]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(_ context.Context, tCtx K) (any, error) {
			return pprofile.FromAttributeIndices(tCtx.GetProfilesDictionary().AttributeTable(), tCtx.GetProfile()), nil
		},
		Setter: func(_ context.Context, tCtx K, val any) error {
			m, err := ctxutil.GetMap(val)
			if err != nil {
				return err
			}
			tCtx.GetProfile().AttributeIndices().FromRaw([]int32{})
			for k, v := range m.All() {
				if err := pprofile.PutAttribute(tCtx.GetProfilesDictionary().AttributeTable(), tCtx.GetProfile(), k, v); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func accessAttributesKey[K Context](key []ottl.Key[K]) ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(ctx context.Context, tCtx K) (any, error) {
			return ctxutil.GetMapValue[K](ctx, tCtx, pprofile.FromAttributeIndices(tCtx.GetProfilesDictionary().AttributeTable(), tCtx.GetProfile()), key)
		},
		Setter: func(ctx context.Context, tCtx K, val any) error {
			newKey, err := ctxutil.GetMapKeyName(ctx, tCtx, key[0])
			if err != nil {
				return err
			}
			v := getAttributeValue(tCtx, *newKey)
			if err = ctxutil.SetIndexableValue[K](ctx, tCtx, v, val, key[1:]); err != nil {
				return err
			}
			return pprofile.PutAttribute(tCtx.GetProfilesDictionary().AttributeTable(), tCtx.GetProfile(), *newKey, v)
		},
	}
}

func getAttributeValue[K Context](tCtx K, key string) pcommon.Value {
	// Find the index of the attribute in the profile's attribute indices
	// and return the corresponding value from the attribute table.
	table := tCtx.GetProfilesDictionary().AttributeTable()
	indices := tCtx.GetProfile().AttributeIndices().AsRaw()

	for _, tableIndex := range indices {
		attr := table.At(int(tableIndex))
		if attr.Key() == key {
			v := pcommon.NewValueEmpty()
			attr.Value().CopyTo(v)
			return v
		}
	}

	return pcommon.NewValueEmpty()
}
