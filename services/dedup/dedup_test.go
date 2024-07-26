package dedup_test

import (
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/rudderlabs/rudder-go-kit/config"
	"github.com/rudderlabs/rudder-go-kit/logger"
	"github.com/rudderlabs/rudder-go-kit/stats"
	"github.com/rudderlabs/rudder-go-kit/testhelper/rand"

	"github.com/rudderlabs/rudder-server/services/dedup"
	"github.com/rudderlabs/rudder-server/services/dedup/types"
	"github.com/rudderlabs/rudder-server/utils/misc"
)

func Test_Dedup(t *testing.T) {
	config.Reset()
	logger.Reset()
	misc.Init()

	dbPath := os.TempDir() + "/dedup_test"
	defer func() { _ = os.RemoveAll(dbPath) }()
	_ = os.RemoveAll(dbPath)
	conf := config.New()
	t.Setenv("RUDDER_TMPDIR", dbPath)
	d, err := dedup.New(conf, stats.Default)
	require.Nil(t, err)
	defer d.Close()

	t.Run("if message id is not present in cache and badger db", func(t *testing.T) {
		found, _, err := d.Set(types.KeyValue{Key: "a", Value: 1})
		require.Nil(t, err)
		require.Equal(t, true, found)

		// Checking it again should give us the previous value from the cache
		found, value, err := d.Set(types.KeyValue{Key: "a", Value: 2})
		require.Nil(t, err)
		require.Equal(t, false, found)
		require.Equal(t, int64(1), value)
	})

	t.Run("if message is committed, previous value should always return", func(t *testing.T) {
		found, _, err := d.Set(types.KeyValue{Key: "b", Value: 1})
		require.Nil(t, err)
		require.Equal(t, true, found)

		err = d.Commit(map[string]types.KeyValue{"a": {Key: "a", Value: 1}})
		require.NoError(t, err)

		found, value, err := d.Set(types.KeyValue{Key: "b", Value: 2})
		require.Nil(t, err)
		require.Equal(t, false, found)
		require.Equal(t, int64(1), value)
	})

	t.Run("committing a messageid not present in cache", func(t *testing.T) {
		found, _, err := d.Set(types.KeyValue{Key: "c", Value: 1})
		require.Nil(t, err)
		require.Equal(t, true, found)

		err = d.Commit(map[string]types.KeyValue{"d": {Key: "d", Value: 2}})
		require.NotNil(t, err)
	})
}

func Test_Dedup_Window(t *testing.T) {
	config.Reset()
	logger.Reset()
	misc.Init()

	dbPath := os.TempDir() + "/dedup_test"
	conf := config.New()
	defer func() { _ = os.RemoveAll(dbPath) }()
	_ = os.RemoveAll(dbPath)
	conf.Set("Dedup.dedupWindow", "1s")
	t.Setenv("RUDDER_TMPDIR", dbPath)
	d, err := dedup.New(conf, stats.Default)
	require.Nil(t, err)
	defer d.Close()

	found, _, err := d.Set(types.KeyValue{Key: "to be deleted", Value: 1})
	require.Nil(t, err)
	require.Equal(t, true, found)

	key := "to be deleted"
	value := types.KeyValue{Key: key, Value: 2}
	err = d.Commit(map[string]types.KeyValue{key: value})
	require.NoError(t, err)

	found, _, err = d.Set(types.KeyValue{Key: "to be deleted", Value: 2})
	require.Nil(t, err)
	require.Equal(t, false, found)

	require.Eventually(t, func() bool {
		found, _, err = d.Set(types.KeyValue{Key: "to be deleted", Value: 3})
		require.Nil(t, err)
		return found
	}, 2*time.Second, 100*time.Millisecond)
}

func Test_Dedup_ErrTxnTooBig(t *testing.T) {
	config.Reset()
	logger.Reset()
	misc.Init()

	dbPath := os.TempDir() + "/dedup_test_errtxntoobig"
	defer os.RemoveAll(dbPath)
	os.RemoveAll(dbPath)
	conf := config.New()
	t.Setenv("RUDDER_TMPDIR", dbPath)
	d, err := dedup.New(conf, stats.Default)
	require.Nil(t, err)
	defer d.Close()

	size := 105_000
	messages := map[string]types.KeyValue{}
	for i := 0; i < size; i++ {
		key := uuid.New().String()
		messages[key] = types.KeyValue{Key: key, Value: int64(i + 1)}
		_, _, _ = d.Set(messages[key])
	}
	err = d.Commit(messages)
	require.NoError(t, err)
}

func Benchmark_Dedup(b *testing.B) {
	config.Reset()
	logger.Reset()
	misc.Init()

	dbPath := path.Join("./testdata", "tmp", rand.String(10), "/DB_Benchmark_Dedup")
	b.Logf("using path %s, since tmpDir has issues in macOS\n", dbPath)
	defer func() { _ = os.RemoveAll(dbPath) }()
	_ = os.MkdirAll(dbPath, 0o750)
	conf := config.New()
	b.Setenv("RUDDER_TMPDIR", dbPath)
	d, err := dedup.New(conf, stats.Default)
	require.NoError(b, err)

	b.Run("no duplicates 1000 batch unique", func(b *testing.B) {
		batchSize := 1000

		msgIDs := make([]types.KeyValue, batchSize)
		messages := map[string]types.KeyValue{}

		for i := 0; i < b.N; i++ {
			key := uuid.New().String()
			msgIDs[i%batchSize] = types.KeyValue{
				Key:   key,
				Value: int64(i + 1),
			}
			messages[key] = msgIDs[i%batchSize]
			if i%batchSize == batchSize-1 || i == b.N-1 {
				for _, msgID := range msgIDs[:i%batchSize] {
					_, _, _ = d.Set(msgID)
				}
				err := d.Commit(messages)
				require.NoError(b, err)
				messages = nil
			}
		}
		b.ReportMetric(float64(b.N), "events")
		b.ReportMetric(float64(b.N*len(uuid.New().String())), "bytes")
	})
	d.Close()

	cmd := exec.Command("du", "-sh", dbPath)
	out, err := cmd.Output()
	if err != nil {
		b.Log(err)
	}

	b.Log("db size:", string(out))
}
