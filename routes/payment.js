import express from "express";
import { authenticate, isInstructor } from "../middlewares";
import {
  createPayment,
  handlePaymentNotification,
  getPaymentStatus,
  getInstructorEarnings,
} from "../controllers/payment";

const router = express.Router();

// Payment routes
router.post("/payment/create", authenticate, createPayment);
router.post("/payment/notification", handlePaymentNotification);
router.get("/payment/status/:orderId", authenticate, getPaymentStatus);
router.get(
  "/instructor/earnings",
  authenticate,
  isInstructor,
  getInstructorEarnings
);

module.exports = router;
