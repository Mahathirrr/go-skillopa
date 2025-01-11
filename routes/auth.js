import express from 'express';
import { register, login, googleLogin, logout } from '../controllers/auth';

const router = express.Router();

router.post('/auth/register', register);
router.post('/auth/login', login);
router.post('/auth/google-login', googleLogin);
router.post('/auth/logout', logout);

module.exports = router;