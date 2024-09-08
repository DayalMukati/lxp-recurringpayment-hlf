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
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client org: %s", err.Error())
	}
	if clientOrgID != "Org1MSP" {
		return fmt.Errorf("only Org1 can create a subscription")
	}

	subscription := Subscription{
		SubscriptionID:   subscriptionID,
		PayerID:          payerID,
		PayeeID:          payeeID,
		Amount:           amount,
		TotalInstallments: totalInstallments,
		PaymentsMade:     0,
		Frequency:        frequency,
		Completed:        false,
	}

	subscriptionAsBytes, err := json.Marshal(subscription)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription: %s", err.Error())
	}

	return ctx.GetStub().PutState(subscriptionID, subscriptionAsBytes)
}

// MakePayment can only be invoked by Org1 users
func (rpc *RecurringPaymentContract) MakePayment(ctx contractapi.TransactionContextInterface, subscriptionID string) error {
	// Verify that the transaction is submitted by Org1
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client org: %s", err.Error())
	}
	if clientOrgID != "Org1MSP" {
		return fmt.Errorf("only Org1 can make payments")
	}

	subscriptionAsBytes, err := ctx.GetStub().GetState(subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %s", err.Error())
	}

	if subscriptionAsBytes == nil {
		return fmt.Errorf("subscription %s does not exist", subscriptionID)
	}

	subscription := new(Subscription)
	err = json.Unmarshal(subscriptionAsBytes, subscription)
	if err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %s", err.Error())
	}

	if subscription.Completed {
		return fmt.Errorf("subscription %s is already completed", subscriptionID)
	}

	// Check if there are remaining installments
	if subscription.PaymentsMade >= subscription.TotalInstallments {
		return fmt.Errorf("no remaining installments for subscription %s", subscriptionID)
	}

	// Increment the number of payments made
	subscription.PaymentsMade++

	// If all payments are made, mark the subscription as completed
	if subscription.PaymentsMade >= subscription.TotalInstallments {
		subscription.Completed = true
	}

	subscriptionAsBytes, err = json.Marshal(subscription)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription: %s", err.Error())
	}

	return ctx.GetStub().PutState(subscriptionID, subscriptionAsBytes)
}

// ConfirmPayment can only be invoked by Org2 users
func (rpc *RecurringPaymentContract) ConfirmPayment(ctx contractapi.TransactionContextInterface, subscriptionID string) error {
	// Verify that the transaction is submitted by Org2
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client org: %s", err.Error())
	}
	if clientOrgID != "Org2MSP" {
		return fmt.Errorf("only Org2 can confirm payments")
	}

	subscriptionAsBytes, err := ctx.GetStub().GetState(subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %s", err.Error())
	}

	if subscriptionAsBytes == nil {
		return fmt.Errorf("subscription %s does not exist", subscriptionID)
	}

	// If the payments made are confirmed, just print confirmation
	subscription := new(Subscription)
	err = json.Unmarshal(subscriptionAsBytes, subscription)
	if err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %s", err.Error())
	}

	fmt.Printf("Payment for subscription %s confirmed by payee.\n", subscriptionID)
	return nil
}

// QuerySubscriptionStatus can only be invoked by Org2 users
func (rpc *RecurringPaymentContract) QuerySubscriptionStatus(ctx contractapi.TransactionContextInterface, subscriptionID string) (*Subscription, error) {
	// Verify that the transaction is submitted by Org2
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("failed to get client org: %s", err.Error())
	}
	if clientOrgID != "Org2MSP" {
		return nil, fmt.Errorf("only Org2 can query the subscription status")
	}

	subscriptionAsBytes, err := ctx.GetStub().GetState(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %s", err.Error())
	}

	if subscriptionAsBytes == nil {
		return nil, fmt.Errorf("subscription %s does not exist", subscriptionID)
	}

	subscription := new(Subscription)
	err = json.Unmarshal(subscriptionAsBytes, subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription: %s", err.Error())
	}

	return subscription, nil
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
