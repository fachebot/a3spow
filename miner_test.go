package a3spow

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLongRepeatedFilter_Filter(t *testing.T) {
	//{
	//	filter := &LongRepeatedFilter{
	//		MinLength: 3,
	//	}
	//	require.False(t, filter.Filter("0123456789012345678901234567890123456789"))
	//	require.True(t, filter.Filter("0003456789012345678901234567890123456789"))
	//	require.True(t, filter.Filter("0123456789012345678901234567890123450000"))
	//}
	//
	//{
	//	char := byte('a')
	//	filter := &LongRepeatedFilter{
	//		MinLength: 3,
	//		Char:      &char,
	//	}
	//	require.False(t, filter.Filter("0123456789012345678901234567890123456789"))
	//	require.False(t, filter.Filter("0003456789012345678901234567890123456789"))
	//	require.True(t, filter.Filter("012345678901234567890123456789012345aaaa"))
	//}
	//
	//{
	//	char := byte('a')
	//	maxStartIndex := 3
	//	filter := &LongRepeatedFilter{
	//		MinLength:     3,
	//		MaxStartIndex: &maxStartIndex,
	//	}
	//	require.False(t, filter.Filter("012345678901234567890123456789012345aaaa"))
	//	require.True(t, filter.Filter("aaa3456789012345678901234567890123456789"))
	//	require.True(t, filter.Filter("012bbb6789012345678901234567890123456789"))
	//
	//	filter.Char = &char
	//	require.True(t, filter.Filter("aaa3456789012345678901234567890123456789"))
	//	require.False(t, filter.Filter("012bbb6789012345678901234567890123456789"))
	//}

	{
		maxStartIndex := 3
		filter := &LongRepeatedFilter{
			MinLength:     3,
			Reverse:       true,
			MaxStartIndex: &maxStartIndex,
		}
		//require.False(t, filter.Filter("aaa3456789012345678901234567890123456789"))
		//require.False(t, filter.Filter("012bbb6789012345678901234567890123456789"))
		require.True(t, filter.Filter("012345678901234567890123456789012345aaaa"))
		//	require.True(t, filter.Filter("0123456789012345678901234567890120000789"))
	}
}
