package txdb

import (
	"bytes"
	"os"
	"testing"

	dbm "github.com/tendermint/tmlibs/db"

	"github.com/bytom/protocol/bc"
)

func TestCleanSnapshotDB(t *testing.T) {
	testDB := dbm.NewDB("testdb", "leveldb", "temp")
	defer os.RemoveAll("temp")

	// Insert the test data
	hash := &bc.Hash{}
	for i := uint64(0); i <= uint64(10); i++ {
		hash.V0 = i
		testDB.Set(calcSnapshotKey(hash), nil)
	}
	testDB.SetSync(nil, nil)

	// run the test function
	cleanSnapshotDB(testDB, hash)

	// check the clean result
	iter := testDB.IteratorPrefix([]byte(snapshotPreFix))
	defer iter.Release()

	if !iter.Next() || !bytes.Equal(iter.Key(), calcSnapshotKey(hash)) {
		t.Errorf("latest snapshot get deleted from db")
	}
	if iter.Next() {
		t.Errorf("more than one snapshot still saved in the db")
	}
}
