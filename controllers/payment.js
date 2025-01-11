import { v4 as uuidv4 } from "uuid";
import { Payment } from "../models/payment";
import { Course } from "../models/course";
import { User } from "../models/user";
import { Instructor } from "../models/instructor";
import { createPaymentToken, checkTransactionStatus } from "../utils/midtrans";

// Create payment for a course
export const createPayment = async (req, res) => {
  try {
    const { courseId } = req.body;
    const userId = req.user._id;

    // Get course details
    const course = await Course.findById(courseId).populate("instructors");
    if (!course) {
      return res.status(404).json({ error: "Course not found" });
    }

    // Check if course is paid
    if (course.pricing !== "Paid") {
      return res.status(400).json({ error: "This course is free" });
    }

    // Check if user already enrolled
    const user = await User.findById(userId);
    if (user.enrolledCourses.some((c) => c.course.toString() === courseId)) {
      return res.status(400).json({ error: "Already enrolled in this course" });
    }

    // Create unique order ID
    const orderId = `ORDER-${uuidv4()}`;

    // Create payment record
    const payment = new Payment({
      orderId,
      course: courseId,
      user: userId,
      instructor: course.instructors[0]._id,
      amount: course.price,
      currency: course.currency,
    });

    // Get payment token from Midtrans
    const transaction = await createPaymentToken({
      orderId,
      amount: course.price,
      user,
      course,
    });

    // Update payment with token response
    payment.paymentLink = transaction.redirect_url;
    await payment.save();

    res.status(200).json({
      orderId,
      paymentLink: transaction.redirect_url,
      token: transaction.token,
    });
  } catch (error) {
    console.error("Create payment error:", error);
    res.status(500).json({ error: "Failed to create payment" });
  }
};

// Handle Midtrans webhook notification
export const handlePaymentNotification = async (req, res) => {
  try {
    const notification = req.body;
    const orderId = notification.order_id;
    const transactionStatus = notification.transaction_status;
    const fraudStatus = notification.fraud_status;

    const payment = await Payment.findOne({ orderId });
    if (!payment) {
      return res.status(404).json({ error: "Payment not found" });
    }

    // Check transaction status
    const statusResponse = await checkTransactionStatus(orderId);
    payment.midtransResponse = statusResponse;

    // Update payment status based on transaction status
    if (transactionStatus === "capture") {
      if (fraudStatus === "challenge") {
        payment.status = "pending";
      } else if (fraudStatus === "accept") {
        payment.status = "success";
      }
    } else if (transactionStatus === "settlement") {
      payment.status = "success";
    } else if (
      transactionStatus === "cancel" ||
      transactionStatus === "deny" ||
      transactionStatus === "expire"
    ) {
      payment.status = "failed";
    } else if (transactionStatus === "pending") {
      payment.status = "pending";
    }

    await payment.save();

    // If payment is successful, enroll user in course
    if (payment.status === "success") {
      // Update user's enrolled courses
      await User.findByIdAndUpdate(payment.user, {
        $addToSet: {
          enrolledCourses: {
            course: payment.course,
            enrolledOn: new Date(),
          },
        },
      });

      // Update course enrollments
      await Course.findByIdAndUpdate(payment.course, {
        $addToSet: {
          "meta.enrollments": { id: payment.user },
        },
      });

      // Update instructor earnings
      await Instructor.findByIdAndUpdate(payment.instructor, {
        $inc: { "meta.totalEarnings": payment.amount },
      });
    }

    res.status(200).json({ status: "OK" });
  } catch (error) {
    console.error("Payment notification error:", error);
    res.status(500).json({ error: "Failed to process payment notification" });
  }
};

// Get payment status
export const getPaymentStatus = async (req, res) => {
  try {
    const { orderId } = req.params;
    const payment = await Payment.findOne({ orderId })
      .populate("course")
      .populate("instructor");

    if (!payment) {
      return res.status(404).json({ error: "Payment not found" });
    }

    res.status(200).json(payment);
  } catch (error) {
    console.error("Get payment status error:", error);
    res.status(500).json({ error: "Failed to get payment status" });
  }
};

// Get instructor earnings
export const getInstructorEarnings = async (req, res) => {
  try {
    const instructorId = req.user.instructorProfile;

    const payments = await Payment.find({
      instructor: instructorId,
      status: "success",
    })
      .populate("course", "title")
      .populate("user", "name email")
      .sort("-createdAt");

    const totalEarnings = payments.reduce(
      (sum, payment) => sum + payment.amount,
      0
    );

    res.status(200).json({
      totalEarnings,
      payments,
    });
  } catch (error) {
    console.error("Get instructor earnings error:", error);
    res.status(500).json({ error: "Failed to get earnings" });
  }
};
