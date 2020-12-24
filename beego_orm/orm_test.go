package beego_orm

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/astaxie/beego/orm"
	"testing"
	"time"
	"won-models"
)

func TestTransactionWithRepeatableReadIsolation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	o := orm.NewOrm()
	o.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	a := &models.Accounts{MemberId: 1, CurrencyId: "won"}
	err := o.ReadForUpdate(a, "member_id", "currency_id")
	fmt.Println(a, err)

	a.UpdatedAt = time.Now()

	_, err = o.Update(a)
	fmt.Println(a, err)

	err = o.Commit()
	if err != nil && err != orm.ErrTxDone && err.Error() != "context canceled" {
		fmt.Printf("commit: %T\n", err)
		fmt.Printf("commit: %#v\n", err)
	}
}