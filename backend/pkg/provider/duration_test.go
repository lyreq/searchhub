package provider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDurationToSeconds(t *testing.T) {
	require.Equal(t, 1515, parseDurationToSeconds("25:15"))
	require.Equal(t, 3723, parseDurationToSeconds("1:02:03"))
	require.Equal(t, 1080, parseDurationToSeconds("18:00"))
	require.Equal(t, 90, parseDurationToSeconds("90"))  // fallback: sayı
	require.Equal(t, 0, parseDurationToSeconds(""))     // boş
	require.Equal(t, 0, parseDurationToSeconds("oops")) // hatalı format → 0
}
