import express from "express";
import { readdirSync } from "fs";
import cors from "cors";
import dotenv from "dotenv";
import mongoose from "mongoose";

const logger = require("morgan");

// Create express app
const app = express();
dotenv.config();

mongoose
  .connect(process.env.DATABASE)
  .then(() => console.log("DB Connected"))
  .catch((e) => {
    console.error("MongoDB connection error:");
    console.error(e);
    process.exit(1); // Exit process with failure
  });

// Apply Middlewares
app.use(logger("dev"));
app.use(express.json());

// Update CORS configuration to allow credentials and specific origin
app.use(
  cors({
    origin: process.env.CLIENT_URL,
    credentials: true,
    methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    allowedHeaders: ["Content-Type", "Authorization"],
  })
);

app.get("/", (req, res) => {
  res.send("API available at /api");
});

readdirSync("./routes").map((r) => app.use("/api", require(`./routes/${r}`)));

const PORT = process.env.PORT || 8000;

app.listen(PORT, () => console.log(`Server is running on port ${PORT}`));
