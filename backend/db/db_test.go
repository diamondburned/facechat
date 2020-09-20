// +build integration

package db

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/internal/dotenv"
)

func TestDB(t *testing.T) {
	db, err := Open(dotenv.Getenv("TEST_SQL_ADDRESS"))
	if err != nil {
		t.Fatal("failed to open database:", err)
	}
	t.Cleanup(db.dropTestDatabase)
	var user *facechat.User
	err = db.Acquire(new5sCtx(t), 0,
		func(tx *Tx) error {
			var err error
			user, _, err = tx.Register("dmr", "p/q2-q4!", "dmr@bell-labs.com")
			return err
		},
	)
	if err != nil {
		t.Fatal("failed to register user:", err)
	}

	err = db.Acquire(new5sCtx(t), 0,
		func(tx *Tx) error {
			_, _, err := tx.Register("dmr", "p/q2-q4!", "dmr@bell-labs.com")
			return err
		},
	)
	if err == nil {
		t.Fatal("duplicate user creation didn't fail")
	}

	err = db.Acquire(new5sCtx(t), user.ID, testAddAccounts)
	if err != nil {
		t.Fatal("failed to add accounts:", err)
	}
	err = db.Acquire(new5sCtx(t), 0, testLogin)
	if err != nil {
		t.Fatal(err)
	}
	err = db.RAcquire(new5sCtx(t), user.ID, testGetUser)
	if err != nil {
		t.Fatal("failed to get user:", err)
	}
	err = db.Acquire(new5sCtx(t), user.ID, testCreateJoinRoom)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Acquire(new5sCtx(t), user.ID, testJoinedRoomsAndCreateMessage)
	if err != nil {
		t.Fatal(err)
	}
}

func testAddAccounts(tx *Tx) error {
	if err := tx.AddAccount(facechat.Account{Name: "GitHub"}); err != nil {
		return err
	}
	return tx.AddAccount(facechat.Account{Name: "Twitter"})
}

func testLogin(tx *Tx) error {
	_, err := tx.Login("dmr@bell-labs.com", "p/q2-q4!")
	if err != nil {
		return errors.Wrap(err, "login failed:")
	}
	_, err = tx.Login("dmr@bell-labs.com", "incorrectpassword")
	if err == nil {
		return errors.New("login with incorrect password succeeded:")
	}
	return nil
}

func testGetUser(tx *ReadTx) error {
	_, err := tx.User(tx.UserID)
	return err
}

func testCreateJoinRoom(tx *Tx) error {
	room, err := tx.CreatePublicLobby("roomname", facechat.Anonymous)
	if err != nil {
		return errors.Wrap(err, "error creating room")
	}
	if err = tx.JoinRoom(room.ID); err != nil {
		return errors.Wrap(err, "error joining room")
	}
	return nil
}

func testJoinedRoomsAndCreateMessage(tx *Tx) error {
	rooms, err := tx.JoinedRooms()
	if err != nil {
		return err
	}
	if rooms == nil {
		return errors.New("no joined rooms")
	}
	_, err = tx.CreateMessage(rooms[0].ID, "Hello")
	return err
}

func new5sCtx(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func (db *DB) dropTestDatabase() {
	_, err := db.db.Exec(`DO $$ DECLARE
	r RECORD;
BEGIN
	FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
		EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
	END LOOP;
END $$;`)
	if err != nil {
		log.Fatalln(err)
	}
}
