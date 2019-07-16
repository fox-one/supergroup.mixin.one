package models

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	bot "github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/lib/pq"
)

const referral_DDL = `
CREATE TABLE IF NOT EXISTS referrals (
	code         			VARCHAR(36) PRIMARY KEY,
	inviter_id        VARCHAR(36) NOT NULL,
	invitee_id	      VARCHAR(36),
	is_used       	  BOOLEAN NOT NULL DEFAULT FALSE,
	created_at       	TIMESTAMP WITH TIME ZONE NOT NULL,
	used_at        		TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS referrals_inviterx ON referrals(inviter_id);
`

type Referral struct {
	Code      string
	InviterID string
	InviteeID string
	IsUsed    bool 	
	CreatedAt time.Time
	UsedAt    pq.NullTime
}

var referralColumns = []string{"code", "inviter_id", "invitee_id", "is_used", "created_at", "used_at"}

func (r *Referral) values() []interface{} {
	return []interface{}{r.Code, r.InviterID, r.InviteeID, r.IsUsed, r.CreatedAt, r.UsedAt}
}

func referralFromRow(row durable.Row) (*Referral, error) {
	var r Referral
	err := row.Scan(&r.Code, &r.InviterID, &r.InviteeID, &r.IsUsed, &r.CreatedAt, &r.UsedAt)
	return &r, err
}

func (user *User) Referrals(ctx context.Context) ([]*Referral, error) {
	var referrals []*Referral
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM referrals WHERE inviter_id = '%s' AND is_used = %t", strings.Join(referralColumns, ","), user.UserId, false)
		rows, err := tx.QueryContext(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			referral, err := referralFromRow(rows)
			if err != nil {
				return err
			}
			referrals = append(referrals, referral)
		}
		return nil
	})
	if err != nil {
		if sessionErr, ok := err.(session.Error); ok {
			return nil, sessionErr
		}
		return nil, session.TransactionError(ctx, err)
	}
	return referrals, nil
}

func (user *User) CreateReferrals(ctx context.Context) ([]*Referral, error) {
	var referrals []*Referral
	currentReferrals, err := user.Referrals(ctx)
	if err != nil {
		return nil, err
	} else if referralCount:= len(currentReferrals); referralCount > 0 {
		return nil, fmt.Errorf("There are %d unused codes, can't create new one", referralCount)
	} else {
		var values bytes.Buffer
		for i := 1;  i<=3; i++ {
			referral := &Referral{InviterID: user.UserId, Code: bot.UuidNewV4().String(), CreatedAt: time.Now(), IsUsed: false}
			if i > 1 {
				values.WriteString(",")
			}
			values.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%t', '%s')", referral.Code, referral.InviterID, referral.InviteeID, referral.IsUsed, string(pq.FormatTimestamp(referral.CreatedAt))))
			referrals = append(referrals, referral)
		}
		query := fmt.Sprintf("INSERT INTO referrals (code,inviter_id,invitee_id,is_used,created_at) VALUES %s", values.String())
		_, err := session.Database(ctx).ExecContext(ctx, query)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		return referrals, nil
	}
}

func (user *User) ApplyReferral(ctx context.Context, referralCode string) (*Referral, error) {
	// return err if user is already joined
	var referral *Referral
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		referral, err = findReferralByCode(ctx, tx, referralCode)
		if err != nil {
			return err
		}

		referral.InviteeID = user.UserId
		referral.IsUsed = true
		referral.UsedAt = pq.NullTime{Time: time.Now(), Valid: true}
		query := fmt.Sprintf("UPDATE referrals SET (invitee_id,is_used,used_at)=($1,$2,$3) WHERE code=$4")
		_, err = tx.ExecContext(ctx, query, referral.InviteeID, referral.IsUsed, referral.UsedAt, referralCode) 
		if err != nil {
			return err
		}

		query = fmt.Sprintf("UPDATE users SET state=$1 WHERE user_id=$2")
		_, err = tx.ExecContext(ctx, query, PaymentStateInvited, user.UserId) 
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if sessionErr, ok := err.(session.Error); ok {
			return nil, sessionErr
		}
		return nil, session.TransactionError(ctx, err)
	}
	
	return referral, nil
}

func findReferralByCode(ctx context.Context, tx *sql.Tx, code string) (*Referral, error) {
	query := fmt.Sprintf("SELECT %s FROM referrals WHERE code = $1 LIMIT 1", strings.Join(referralColumns, ","))
	row := tx.QueryRowContext(ctx, query, code)
	referral, err := referralFromRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return referral, err
}