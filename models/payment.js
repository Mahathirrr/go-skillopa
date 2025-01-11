import mongoose from "mongoose";
const { Schema } = mongoose;
const { ObjectId } = Schema;

const paymentSchema = new Schema(
  {
    orderId: {
      type: String,
      required: true,
      unique: true,
    },
    course: {
      type: ObjectId,
      ref: "Course",
      required: true,
    },
    user: {
      type: ObjectId,
      ref: "User",
      required: true,
    },
    instructor: {
      type: ObjectId,
      ref: "Instructor",
      required: true,
    },
    amount: {
      type: Number,
      required: true,
    },
    currency: {
      type: String,
      required: true,
      enum: ["IDR", "USD"],
    },
    status: {
      type: String,
      enum: ["pending", "success", "failed", "expired"],
      default: "pending",
    },
    paymentMethod: String,
    midtransResponse: Schema.Types.Mixed,
    paymentLink: String,
  },
  {
    timestamps: true,
  }
);

// Index for faster queries
paymentSchema.index({ orderId: 1 });
paymentSchema.index({ user: 1 });
paymentSchema.index({ instructor: 1 });
paymentSchema.index({ course: 1 });
paymentSchema.index({ status: 1 });

export const Payment = mongoose.model("Payment", paymentSchema);
