#!/bin/bash

# Initialize score
score=0
source ./scripts/setOrgPeerContext.sh 2
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx

# Query the subscription status
echo "Querying the subscription status..."
QUERY_OUTPUT=$(peer chaincode query -C mychannel -n recurringpaymentorgs -c '{"Args":["QuerySubscriptionStatus","sub1"]}' 2>&1)

echo "Query output: $QUERY_OUTPUT"
# Check if the subscription exists
if [[ $QUERY_OUTPUT == *"Error"* ]]; then
    echo "Subscription not found. It was not created."
    # Output the final score
    echo "Final Score: $score/50"
    exit 0
else
    echo "Subscription found."
fi

# Check if the subscription was created (indicating subscription creation)
if [[ $QUERY_OUTPUT == *"subscriptionID\":\"sub1\""* ]]; then
    echo "Subscription creation successful."
    score=$((score + 15))
fi

# Check if a payment was made (indicating payment success)
if [[ $QUERY_OUTPUT == *"paymentsMade\":1"* ]]; then
    echo "Payment successful."
    score=$((score + 15))
fi

# Check if the payment was confirmed (indicating payment confirmation)
if [[ $QUERY_OUTPUT == *"completed\":true"* ]]; then
    echo "Payment confirmation successful."
    score=$((score + 15))
fi

# Check the full subscription status (for additional checks)
EXPECTED_OUTPUT='{"subscriptionID":"sub1","payerID":"payer1","payeeID":"payee1","amount":1000,"totalInstallments":6'
if [[ $QUERY_OUTPUT == *"$EXPECTED_OUTPUT"* ]]; then
    echo "Subscription status query successful."
    score=$((score + 5))
else
    echo "Subscription status query incomplete."
fi

# Final score output
echo "Final Score: $score/50"

# Exit with success
exit 0
