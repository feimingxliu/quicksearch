package test

import (
	"fmt"
	"github.com/feimingxliu/quicksearch/pkg/util/random"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"testing"
)

func TestBolt(t *testing.T) {
	db, err := bolt.Open("testdata/bolt.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	db.Info()
	// Start a writable transaction.
	tx, err := db.Begin(true)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("++++++ A transaction started ++++++")
	// Use the transaction...
	bucketName := random.RandomString(8)
	bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
	log.Printf("++++++ Bucket '%s' created ++++++\n", bucketName)
	if err != nil {
		t.Fatal(err)
	}
	var key, value string
	for i := 0; i < 10000; i++ {
		key = random.RandomString(4)
		value = random.RandomString(10)
		err = bucket.Put([]byte(key), []byte(value))
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("PUT key: %s, value: %s \n", key, value)
	}
	// Commit the transaction and check for error.
	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}
	log.Println("++++++ The transaction committed ++++++")

	log.Printf("++++++ Start retriving bucket '%s' ++++++", bucketName)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()
		var total uint
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s \n", k, v)
			total++
		}
		log.Printf("++++++ End retrieve, total %d ++++++\n", total)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	//close the DB
	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
	wd, _ := os.Getwd()
	fmt.Println("work dir: ", wd)
}
