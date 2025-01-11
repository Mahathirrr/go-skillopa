import jwt from "jsonwebtoken";
import { User } from "../models/user";
import { comparePassword, hashPassword } from "../utils";
import { interests } from "../data/interests";
import { uploadToCloudinary } from "../utils/cloudinary";
import { Course } from "../models/course";

const currentUser = async (req, res) => {
  try {
    const user = await User.findById(req.user._id).select("-password").exec();

    if (user) {
      return res.json(user);
    }

    return res.status(404).send("User not found.");
  } catch (error) {
    console.error(error);
    res.status(401).json(error);
  }
};

// @desc    Add to cart
// @route   POST /api/user/cart
// @access  Private
const addToCart = async (req, res) => {
  try {
    const { _id } = req.user;

    if (!_id) {
      return res.status(400).send("Missing user id");
    }

    const data = await User.findByIdAndUpdate(
      _id,
      {
        $addToSet: {
          cart: req.body.id,
        },
      },
      {
        new: true,
        strict: false,
      }
    ).select("cart");

    return res.status(200).send(data.cart);
  } catch (error) {
    res.status(401).json(error);
  }
};

// @desc    Delete from cart
// @route   DELETE /api/user/cart
// @access  Private
const removeFromCart = async (req, res) => {
  try {
    const { _id } = req.user;

    if (!_id) {
      return res.status(400).send("Missing user id");
    }

    await User.findByIdAndUpdate(_id, {
      $pull: { cart: req.params.id },
    }).lean();

    return res.status(200).json({ success: true });
  } catch (error) {
    res.status(401).json(error);
  }
};

// @desc    Fetch cart
// @route   GET /api/user/cart
// @access  Private
const getCart = async (req, res) => {
  try {
    const { _id } = req.user;

    if (!_id) {
      return res.status(400).send("Missing user id");
    }

    const data = await User.findById(_id)
      .lean()
      .populate({
        path: "cart",
        select: "title instructors level pricing currency price duration slug",
        populate: {
          path: "instructors",
          select: "name",
        },
      })
      .select("cart");

    return res.status(200).send(data.cart);
  } catch (error) {
    res.status(401).json(error);
  }
};

// @desc    Add to wishlist
// @route   POST /api/user/wishlist
// @access  Private
const addToWishlist = async (req, res) => {
  try {
    const { _id } = req.user;

    if (!_id) {
      return res.status(400).send("Missing user id");
    }

    const data = await User.findByIdAndUpdate(
      _id,
      {
        $addToSet: {
          wishlist: req.body.id,
        },
      },
      {
        new: true,
        strict: false,
      }
    ).select("wishlist");

    return res.status(200).json(data.wishlist);
  } catch (error) {
    res.status(401).json(error);
  }
};

// @desc    Delete from wishlist
// @route   DELETE /api/user/wishlist
// @access  Private
const removeFromWishlist = async (req, res) => {
  try {
    const { _id } = req.user;

    if (!_id) {
      return res.status(400).send("Missing user id");
    }

    await User.findByIdAndUpdate(_id, {
      $pull: { wishlist: req.params.id },
    }).lean();

    return res.status(200).json({ success: true });
  } catch (error) {
    res.status(401).json(error);
  }
};

// @desc    Fetch wishlist
// @route   GET /api/user/wishlist
// @access  Private
const getWishlist = async (req, res) => {
  try {
    const { _id } = req.user;

    if (!_id) {
      return res.status(400).send("Missing user id");
    }

    const data = await User.findById(_id)
      .lean()
      .populate({
        path: "wishlist",
        select: "-meta -curriculum",
        populate: {
          path: "instructors",
          select: "name",
        },
      })
      .select("wishlist");

    return res.status(200).send(data.wishlist);
  } catch (error) {
    res.status(401).json(error);
  }
};

const checkout = async (req, res) => {
  try {
    const { _id } = req.user;
    const { ids } = req.body;

    if (!_id) {
      return res.status(400).send("Missing user id");
    }

    if (!ids || !ids.length) {
      return res.status(404).send("No courses to checkout");
    }

    const enrolledCourses = ids.map((id) => ({
      course: id,
      enrolledOn: new Date().toISOString(),
    }));

    const [_, user] = await Promise.all([
      Course.updateMany(
        { _id: { $in: ids }, pricing: { $eq: "Free" } },
        {
          $addToSet: {
            "meta.enrollments": { id: _id },
          },
        },
        {
          new: true,
          strict: false,
        }
      ).lean(),
      User.findByIdAndUpdate(
        _id,
        {
          $addToSet: {
            enrolledCourses,
          },
          $pull: { cart: { $in: ids }, wishlist: { $in: ids } },
        },
        {
          new: true,
          strict: false,
        }
      )
        .lean()
        .select("cart wishlist enrolledCourses"),
    ]);

    return res.status(200).json(user);
  } catch (error) {
    console.error("Checkout error:", error);
    return res.status(401).json(error);
  }
};

const getEnrolledCourses = async (req, res) => {
  try {
    const { _id } = req.user;

    const data = await User.findById(_id)
      .lean()
      .populate({ path: "enrolledCourses.course", select: "-meta -curriculum" })
      .select("enrolledCourses");

    return res.status(200).send(data.enrolledCourses);
  } catch (error) {
    console.error(error);
    return res.status(401).json(error);
  }
};

const updateProfile = async (req, res) => {
  try {
    const { _id } = req.user;
    const { name } = req.body;
    const { files } = req;
    const avatar = files ? files.avatar : undefined;

    if (!_id) {
      return res.status(400).send("Missing user id");
    }

    let updateData = {};

    // Handle name update if provided
    if (name) {
      updateData.name = name;
    }

    // Handle avatar update if provided
    if (avatar) {
      try {
        const result = await uploadToCloudinary(avatar.tempFilePath);
        updateData.avatar = result.secure_url;
      } catch (error) {
        console.error("Cloudinary upload error:", error);
        return res.status(400).json({ errorMessage: "Error uploading image" });
      }
    }

    // Only update if we have changes
    if (Object.keys(updateData).length === 0) {
      return res.status(400).json({ errorMessage: "No update data provided" });
    }

    const user = await User.findByIdAndUpdate(_id, updateData, {
      new: true,
    }).select("-password");

    if (!user) {
      return res.status(404).json({ errorMessage: "User not found" });
    }

    return res.status(200).json(user);
  } catch (error) {
    console.error(error);
    return res.status(400).json({ errorMessage: error.message });
  }
};

export {
  currentUser,
  addToCart,
  removeFromCart,
  getCart,
  addToWishlist,
  getWishlist,
  removeFromWishlist,
  checkout,
  getEnrolledCourses,
  updateProfile,
};
