package utils

import (
	"learnlit/models"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type Transaction struct {
	Token       string
	RedirectURL string
}

func CreatePaymentToken(orderID string, amount float64, user models.User, course models.Course) (*Transaction, error) {
	// Configure Midtrans client
	var client snap.Client
	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	// Create transaction details
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
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
				Price: int64(course.Price),
				Qty:   1,
				Name:  course.Title,
			},
		},
		Callbacks: &snap.Callbacks{
			Finish:  os.Getenv("CLIENT_URL") + "/payment/finish",
			Error:   os.Getenv("CLIENT_URL") + "/payment/error",
			Pending: os.Getenv("CLIENT_URL") + "/payment/pending",
		},
	}

	snapResp, err := client.CreateTransaction(req)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Token:       snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
	}, nil
}

