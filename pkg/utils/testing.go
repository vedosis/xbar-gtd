package utils

import (
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func GetFQFN(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func LoadTestFileData(t *testing.T, filename string) []byte {
	t.Helper()
	data, err := os.ReadFile(filename)
	require.NoError(t, err)
	return data
}
