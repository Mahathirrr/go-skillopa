import midtransClient from "midtrans-client";

// Create Snap API instance
const snap = new midtransClient.Snap({
  isProduction: process.env.MIDTRANS_IS_PRODUCTION === "true",
  serverKey: process.env.MIDTRANS_SERVER_KEY,
  clientKey: process.env.MIDTRANS_CLIENT_KEY,
});

// Create Core API instance
const core = new midtransClient.CoreApi({
  isProduction: process.env.MIDTRANS_IS_PRODUCTION === "true",
  serverKey: process.env.MIDTRANS_SERVER_KEY,
  clientKey: process.env.MIDTRANS_CLIENT_KEY,
});

export const createPaymentToken = async ({ orderId, amount, user, course }) => {
  try {
    const parameter = {
      transaction_details: {
        order_id: orderId,
        gross_amount: amount,
      },
      credit_card: {
        secure: true,
      },
      customer_details: {
        first_name: user.name.split(" ")[0],
        last_name: user.name.split(" ").slice(1).join(" "),
        email: user.email,
      },
      item_details: [
        {
          id: course._id,
          price: amount,
          quantity: 1,
          name: course.title,
          category: course.category,
        },
      ],
      callbacks: {
        finish: `${process.env.CLIENT_URL}/payment/finish`,
        error: `${process.env.CLIENT_URL}/payment/error`,
        pending: `${process.env.CLIENT_URL}/payment/pending`,
      },
    };

    const transaction = await snap.createTransaction(parameter);
    return transaction;
  } catch (error) {
    console.error("Midtrans createPaymentToken error:", error);
    throw error;
  }
};

export const checkTransactionStatus = async (orderId) => {
  try {
    const response = await core.transaction.status(orderId);
    return response;
  } catch (error) {
    console.error("Midtrans checkTransactionStatus error:", error);
    throw error;
  }
};
