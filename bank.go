package main

import (
	"errors"
	"fmt"
	"math"

	"github.com/aerospike/aerospike-client-go"
)

type Bank struct {
	Id      string
	Balance float64
}

var (
	ErrBalanceNotFound = errors.New("Balance not found")
	WPolicy            *aerospike.WritePolicy
)

func init() {
	WPolicy := aerospike.NewWritePolicy(0, 0)
	WPolicy.SendKey = true
}

func BankSaveStats(amount float64, clnt *aerospike.Client) error {
	key, err := aerospike.NewKey(DBNs, DBTblStat, "stat")
	if err != nil {
		return err
	}

	intAmount := int(math.Floor((amount * MonetaryShift) + 0.5))

	_, err = clnt.Operate(WPolicy,
		key,
		aerospike.AddOp(aerospike.NewBin("funds", intAmount)),
		aerospike.AddOp(aerospike.NewBin("counter", 1)))

	return err
}

func BankGetStats(clnt *aerospike.Client) (float64, int, error) {
	key, err := aerospike.NewKey(DBNs, DBTblStat, "stat")
	if err != nil {
		return 0, 0, err
	}

	rec, err := Clnt.Get(nil, key)
	if err != nil {
		return 0, 0, err
	}
	if rec == nil {
		return 0, 0, errors.New("No stats found")
	}

	var counter int
	var funds float64

	if _, ok := rec.Bins["funds"]; ok && rec.Bins["funds"] != nil {
		funds = float64(rec.Bins["funds"].(int)) / MonetaryShift
	}
	if _, ok := rec.Bins["counter"]; ok && rec.Bins["counter"] != nil {
		counter = rec.Bins["counter"].(int)
	}

	return funds, counter, err
}

func NewBank(id string, clnt *aerospike.Client) (*Bank, error) {
	key, err := aerospike.NewKey(DBNs, DBTblAccounts, fmt.Sprintf("%s", id))
	if err != nil {
		return nil, err
	}

	rec, err := Clnt.Get(nil, key)
	if err != nil {
		return nil, err
	}

	var balance float64

	if rec != nil {
		if _, ok := rec.Bins["balance"]; ok && rec.Bins["balance"] != nil {
			balance = float64(rec.Bins["balance"].(int)) / MonetaryShift
		}
	}

	bankModel := Bank{
		Id:      id,
		Balance: balance,
	}

	return &bankModel, nil
}

func (self *Bank) addFunds(amount float64, Clnt *aerospike.Client) (int32, error) {

	intAmount := int(math.Floor((amount * MonetaryShift) + 0.5))

	key, err := aerospike.NewKey(DBNs, DBTblAccounts, self.Id)
	if err != nil {
		return 0, err
	}

	rec, err := Clnt.Operate(WPolicy, key, aerospike.AddOp(aerospike.NewBin("balance", intAmount)))

	return int32(rec.Generation), err
}
