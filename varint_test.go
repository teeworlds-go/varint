package varint_test

import (
	"bytes"
	"io"
	"math"
	"testing"

	"github.com/teeworlds-go/varint"
	"github.com/teeworlds-go/varint/internal/testutils/require"
)

func TestVarintBoundaries(t *testing.T) {
	t.Parallel()

	table := []struct {
		inNumber      int
		expectedBytes int
	}{
		// positive
		{63, 1},            // 2^6 -1
		{64, 2},            // 2^6
		{8191, 2},          // 2^(6+7) -1
		{8192, 3},          // 2^(6+7)
		{1048576 - 1, 3},   // 2^(6+7+7) -1
		{1048576, 4},       // 2^(6+7+7)
		{134217728 - 1, 4}, // 2^(6+7+7+7) -1
		{134217728, 5},     // 2^(6+7+7+7)
		// negative
		{-8191, 2},      // (2^(6+7) -1) *-1
		{-8192, 2},      // (2^(6+7)) *-1
		{-8193, 3},      // (2^(6+7) +1) *-1
		{-1048575, 3},   // (2^(6+7+7) -1) *-1
		{-1048576, 3},   // (2^(6+7+7)) *-1
		{-1048577, 4},   // (2^(6+7+7) +1) *-1
		{-134217727, 4}, // (2^(6+7+7+7) -1) *-1
		{-134217728, 4}, // (2^(6+7+7+7)) *-1
		{-134217729, 5}, // (2^(6+7+7+7) +1) *-1
		// int32 boundaries
		{math.MaxInt32, 5}, // 2^31 -1 = 2147483647
		{math.MinInt32, 5}, // -2^31 = -2147483648
	}

	for _, row := range table {
		varintWriteRead(t, row.inNumber, row.expectedBytes)
	}

	// tested w/o panic code
	//varintWriteRead(t, math.MaxInt32+1, 5) // 2^31 = 2147483648
	//varintWriteRead(t, math.MinInt32-1, 5) // -2^31-1 = -2147483649
}

func TestOverflowVarint(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	table := [][]byte{
		{0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001},
		{0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b0010000},
		{0b11000001, 0b10000001, 0b10000001, 0b10000001, 0b0010000}, // underflow
	}

	for _, row := range table {
		buf := bytes.NewBuffer(row)
		_, err := varint.ReadVarint(buf)
		require.Error(err)

		_, n := varint.Varint(row)
		require.Less(0, n)
	}
}

func TestOverflowBigVarint(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	table := [][]byte{
		{0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b0000001},
		{0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b00000010},
		{0b11000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b10000001, 0b00000010}, // underflow
	}

	for _, row := range table {
		buf := bytes.NewBuffer(row)
		_, err := varint.ReadBigVarint(buf)
		require.Error(err)
		_, n := varint.BigVarint(row)
		require.Less(0, n)
	}
}

func TestEOFVarint(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	buf := []byte{0b10000001, 0b10000001, 0b10000001, 0b00000001}
	b := bytes.NewBuffer(buf)

	i, err := varint.ReadVarint(b)
	require.NoError(err)
	require.NotZero(i)

	_, err = varint.ReadVarint(b)
	require.ErrorIs(err, io.EOF)
}

func TestBigVarintBoundaries(t *testing.T) {
	t.Parallel()

	table := []struct {
		inNumber      int64
		expectedBytes int
	}{
		// positive
		{63, 1},            // 2^6 -1
		{64, 2},            // 2^6
		{8191, 2},          // 2^(6+7) -1
		{8192, 3},          // 2^(6+7)
		{1048576 - 1, 3},   // 2^(6+7+7) -1
		{1048576, 4},       // 2^(6+7+7)
		{134217728 - 1, 4}, // 2^(6+7+7+7) -1
		{134217728, 5},     // 2^(6+7+7+7)
		// big positive
		{17179869183, 5},             // 2^(6+4*7) -1
		{17179869184, 6},             // 2^(6+4*7)
		{2199023255552 - 1, 6},       // 2^(6+5*7) -1
		{2199023255552, 7},           // 2^(6+5*7)
		{281474976710656 - 1, 7},     // 2^(6+6*7) -1
		{281474976710656, 8},         // 2^(6+6*7)
		{36028797018963968 - 1, 8},   // 2^(6+7*7) -1
		{36028797018963968, 9},       // 2^(6+7*7)
		{4611686018427387904 - 1, 9}, // 2^(6+8*7) -1
		{4611686018427387904, 10},    // 2^(6+8*7)
		// negative
		{-8191, 2},      // (2^(6+7) -1) *-1
		{-8192, 2},      // (2^(6+7)) *-1
		{-8193, 3},      // (2^(6+7) +1) *-1
		{-1048575, 3},   // (2^(6+7+7) -1) *-1
		{-1048576, 3},   // (2^(6+7+7)) *-1
		{-1048577, 4},   // (2^(6+7+7) +1) *-1
		{-134217727, 4}, // (2^(6+7+7+7) -1) *-1
		{-134217728, 4}, // (2^(6+7+7+7)) *-1
		{-13421779, 4},  // (2^(6+7+7+7) +1) *-1
		// big negative
		{-17179869183, 5},          // (2^(6+4*7) -1) *-1
		{-17179869184, 5},          // (2^(6+4*7)) *-1
		{-17179869185, 6},          // (2^(6+4*7) +1) *-1
		{-2199023255551, 6},        // (2^(6+5*7) -1) *-1
		{-2199023255552, 6},        // (2^(6+5*7)) *-1
		{-2199023255553, 7},        // (2^(6+5*7) +1) *-1
		{-281474976710655, 7},      // (2^(6+6*7) -1) *-1
		{-281474976710656, 7},      // (2^(6+6*7)) *-1
		{-281474976710657, 8},      // (2^(6+6*7) +1) *-1
		{-36028797018963967, 8},    // (2^(6+7*7) -1) *-1
		{-36028797018963968, 8},    // (2^(6+7*7)) *-1
		{-36028797018963969, 9},    // (2^(6+7*7) +1) *-1
		{-4611686018427387903, 9},  // (2^(6+8*7) -1) *-1
		{-4611686018427387904, 9},  // (2^(6+8*7)) *-1
		{-4611686018427387905, 10}, // (2^(6+8*7) +1) *-1
		// int32 boundaries
		{math.MaxInt32, 5}, // 2^31 -1 = 2147483647
		{math.MinInt32, 5}, // -2^31 = -2147483648
		// int64 boundaries
		{math.MaxInt64, 10}, // 2^63 -1
		{math.MinInt64, 10}, // -2^63
	}

	for _, row := range table {
		bigVarintWriteRead(t, row.inNumber, row.expectedBytes)
	}
}

func varintWriteRead(t *testing.T, inNumber int, expectedBytes int) {
	require := require.New(t)

	buf := make([]byte, varint.MaxVarintLen32)
	written := varint.PutVarint(buf, inNumber)
	require.Equal(expectedBytes, written)
	out, read := varint.Varint(buf)

	require.GreaterOrEqual(1, read, "read must be at least 1 byte")
	require.Equal(inNumber, out, "out == in")
	require.Equal(written, read, "read == written")
}

func bigVarintWriteRead(t *testing.T, inNumber int64, expectedBytes int) {
	require := require.New(t)

	buf := make([]byte, varint.MaxVarintLen64)
	written := varint.PutBigVarint(buf, inNumber)
	require.Equal(expectedBytes, written)
	out, read := varint.BigVarint(buf)

	require.GreaterOrEqual(1, read, "read must be at least 1 byte")
	require.Equal(inNumber, out, "out == in")
	require.Equal(written, read, "read == written")
}
