package utils

import (
	"learnlit/models"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func CreatePaymentToken(orderId string, amount int64, user models.User, course models.Course) (*snap.Response, error) {
	var client snap.Client
	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderId,
			GrossAmt: amount,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    course.ID.Hex(),
				Price: amount,
				Qty:   1,
				Name:  course.Title,
			},
		},
		EnabledPayments: []snap.SnapPaymentType{
			snap.PaymentTypeGopay,
			snap.PaymentTypeBankTransfer,
			snap.PaymentTypeCreditCard,
		},
	}

	resp, err := client.CreateTransaction(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
