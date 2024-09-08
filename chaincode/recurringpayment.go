package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type RecurringPaymentContract struct {
	contractapi.Contract
}

type Subscription struct {
	SubscriptionID   string `json:"subscriptionID"`
	PayerID          string `json:"payerID"`
	PayeeID          string `json:"payeeID"`
	Amount           float64 `json:"amount"`
	TotalInstallments int    `json:"totalInstallments"`
	PaymentsMade     int    `json:"paymentsMade"`
	Frequency        string `json:"frequency"`
	Completed        bool   `json:"completed"`
}

// CreateSubscription can only be invoked by Org1 users
func (rpc *RecurringPaymentContract) CreateSubscription(ctx contractapi.TransactionContextInterface, subscriptionID string, payerID string, payeeID string, amount float64, frequency string, totalInstallments int) error {
	// Verify that the transaction is submitted by Org1
	// Write the logic to CreateSubscription(subscriptionID, payerID, payeeID, amount, frequency, totalInstallments): This can only be done by Org1 users (payer).
}

// MakePayment can only be invoked by Org1 users
func (rpc *RecurringPaymentContract) MakePayment(ctx contractapi.TransactionContextInterface, subscriptionID string) error {
	// Verify that the transaction is submitted by Org1
	// Write the logic to MakePayment(subscriptionID): This can only be done by Org1 users (payer).
}

// ConfirmPayment can only be invoked by Org2 users
func (rpc *RecurringPaymentContract) ConfirmPayment(ctx contractapi.TransactionContextInterface, subscriptionID string) error {
	// Verify that the transaction is submitted by Org2
	// Write the logic to ConfirmPayment(subscriptionID): This can only be done by Org2 users (payee).
}

// QuerySubscriptionStatus can only be invoked by Org2 users
func (rpc *RecurringPaymentContract) QuerySubscriptionStatus(ctx contractapi.TransactionContextInterface, subscriptionID string) (*Subscription, error) {
	// Verify that the transaction is submitted by Org2
	// Write the logic to QuerySubscriptionStatus(subscriptionID): This can only be done by Org2 users (payee).
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(RecurringPaymentContract))
	if err != nil {
		fmt.Printf("Error creating recurring payment contract: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting recurring payment contract: %s", err.Error())
	}
}
