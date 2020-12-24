package gorm_tx_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
    "fmt"
    "reflect"
)

var (
	db *gorm.DB
)

func init() {
	var err error

	db, err = gorm.Open("mysql", "root:@(0.0.0.0:3306)/peatio_dev?charset=utf8&parseTime=True&loc=Asia%2FShanghai")
	if err != nil {
		log.Fatal(err.Error())
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(10)
	db.DB().SetConnMaxLifetime(14400 * time.Second)
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	db.LogMode(true)
}

const (
	TransactionRequired  = 0 // Support a current transaction; create a new one if none exists.
	TransactionSupported = 1 // Support a current transaction; execute without transaction if none exists.
)

type Session struct {
	db     *gorm.DB
	cancel context.CancelFunc
	isTx   bool
}

func (session *Session) Begin() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	tx := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	session.db = tx
	session.cancel = cancel
	session.isTx = true
	return
}

func (session *Session) Commit() (err error) {
    if session.isTx {
        return session.db.Commit().Error
    }
    return nil
}

func (session *Session) End() {
    if session.isTx {
        session.cancel()
        session.reset()
    }
}

func (session *Session) reset() {
    session.db = nil
    session.cancel = nil
    session.isTx = false
}

func (session *Session) BeginTrans(transactionDefinition ...int) (*Transaction, error) {
	var tx *Transaction
	if len(transactionDefinition) == 0 {
		tx = session.transaction(TransactionRequired)
	} else {
		tx = session.transaction(transactionDefinition[0])
	}

	err := tx.BeginTrans()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func ConvertToGormDB(session *Session) (*gorm.DB) {
    if reflect.ValueOf(session).IsNil() {
        return db
    }
    
    tx, err := session.BeginTrans(TransactionSupported)
    if err != nil {
        return db
    }

    return tx.txSession.db
}

func (session *Session) transaction(transactionDefinition int) *Transaction {
	if transactionDefinition > 1 || transactionDefinition < 0 {
		return &Transaction{txSession: session, transactionDefinition: TransactionRequired}
	}
	return &Transaction{txSession: session, transactionDefinition: transactionDefinition}
}

type Transaction struct {
	txSession             *Session
	transactionDefinition int
	isNested              bool
}

func (transaction *Transaction) TransactionDefinition() int {
	return transaction.transactionDefinition
}

func (transaction *Transaction) IsExistingTransaction() bool {
	return transaction.txSession.isTx
}

func (transaction *Transaction) Session() *Session {
	return transaction.txSession
}

func (transaction *Transaction) BeginTrans() error {
	switch transaction.transactionDefinition {
	case TransactionRequired:
		if !transaction.IsExistingTransaction() {
			if err := transaction.txSession.Begin(); err != nil {
				return err
			}
		} else {
			transaction.isNested = true
		}
		return nil
	case TransactionSupported:
		if transaction.IsExistingTransaction() {
			transaction.isNested = true
		} else {
            transaction.txSession.db = db
        }
		return nil
	default:
		return errors.New("transaction definition error")
	}
}

func (transaction *Transaction) CommitTrans() error {
	switch transaction.transactionDefinition {
	case TransactionRequired:
		if !transaction.IsExistingTransaction() {
			return errors.New("not in transaction")
		}
		if !transaction.isNested {
			err := transaction.txSession.Commit()
			if err != nil {
				return err
			}
		}
		return nil
	case TransactionSupported:
		// nothing to do
		return nil
	default:
		return errors.New("transaction definition error")
	}
}

func (transaction *Transaction) End() error {
	switch transaction.transactionDefinition {
	case TransactionRequired:
		if !transaction.IsExistingTransaction() {
			return errors.New("not in transaction")
		}
		if !transaction.isNested {
			transaction.txSession.End()
		}
		return nil
	case TransactionSupported:
        // nothing to do
		return nil
	default:
		return errors.New("transaction definition error")
	}
}

func NestedTxRequired(s *Session) error {
    tx, err := s.BeginTrans(TransactionRequired)
    if err != nil {
        return err
    }
    defer tx.End()

    // do sth. using tx.txSession.db

    return tx.CommitTrans()
}

func TxRequired(s *Session) error {
	tx, err := s.BeginTrans(TransactionRequired)
	if err != nil {
		return err
	}
	defer tx.End()

    // do sth. using tx.txSession.db

    err = NestedTxRequired(s)
    if err != nil {
        return err
    }

    // do sth. using tx.txSession.db

    return tx.CommitTrans()
}

func TestTxRequiredSessionBegin(t *testing.T) {
	var session Session
	err := session.Begin()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer session.End()

	err = TxRequired(&session)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = session.Commit()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestTxRequiredWithoutSessionBegin(t *testing.T) {
    var session Session
    err := TxRequired(&session)
    if err != nil {
        t.Fatal(err.Error())
    }
}

func TxSupported(s *Session) error {
    tx, err := s.BeginTrans(TransactionSupported)
    if err != nil {
        return err
    }
    defer tx.End()

    // do sth. using tx.txSession.db

    tx.CommitTrans()
    return nil
}

func TestTxSupportedSessionBegin(t *testing.T) {
    var session Session
    err := session.Begin()
    if err != nil {
        t.Fatal(err.Error())
    }
    defer session.End()

	err = TxSupported(&session)
	if err != nil {
		t.Fatal(err.Error())
	}

    err = session.Commit()
    if err != nil {
        t.Fatal(err.Error())
    }
}

func TestTxSupportedWithoutSessionBegin(t *testing.T) {
    var session Session
    err := TxSupported(&session)
    if err != nil {
        t.Fatal(err.Error())
    }
}

func TxSupportedIncludeTxRequired(s *Session) error {
    tx, err := s.BeginTrans(TransactionSupported)
    if err != nil {
        return err
    }
    defer tx.End()

    db := tx.txSession.db
    // do sth. using db
    fmt.Println(db)

    err = NestedTxRequired(s)
    if err != nil {
        return err
    }

    // [注意] 如果在调用NestedTxRequired之前，事务没有启动，此时的tx.txSession.db会发生变化，而且不可用
    // do sth. using db above instead of tx.txSession.db

    tx.CommitTrans()
    return nil
}

func TestTxSupportedIncludeTxRequiredWithoutSessionBegin(t *testing.T) {
    var session Session
    err := TxSupportedIncludeTxRequired(&session)
    if err != nil {
        t.Fatal(err.Error())
    }
}
