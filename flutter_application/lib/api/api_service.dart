import 'package:dio/dio.dart';
import 'package:flutter_application/api/api_client.dart';
import 'package:flutter_application/services/session_manager.dart';

class ApiService {
  static Future<List<dynamic>> fetchPosts() async {
    final response = await ApiClient.dio.get("/posts");
    return response.data;
  }

  static Future<(bool success, String token, String error)> login(String email, String password) async {
    try {
      final response = await ApiClient.dio.post("/api/auth/login", data: {
        "username": email,
        "password": password,
      });
      final String token = response.data['token'];
      final String error = response.data['error'] ?? "";
      if (error.isNotEmpty) {
        return (false, "", error);
      }
      return (true, token, "");
    } on DioException catch (e) {

      final String errorMessage =
          e.response?.data?['error'] ??
              "An error occurred during login.";

      return (false, "", errorMessage);
    }
  }
  static Future<(bool success, String token, String error)> signup(String username, String email, String password) async {
    try {
      final response = await ApiClient.dio.post("/api/auth/signup", data: {
        "username": username,
        "email": email,
        "password": password,
      });
      final String token = response.data['token'];
      final String error = response.data['error'] ?? "";
      if (error.isNotEmpty) {
        return (false, "", error);
      }
      return (true, token, "");
    } on DioException catch (e) {

      final String errorMessage =
          e.response?.data?['error'] ??
              "An error occurred during signup.";
      return (false, "", errorMessage);
    }
  }
  static Future<(bool valid, bool isAccVerified)> validateAuth() async {
    try {
      final response = await ApiClient.dio.get("/api/auth/validate");
      final String token = response.data['token'] ?? "";
      final String error = response.data['error'] ?? "";
      final bool isVerified = response.data['is_verified'] ?? false;
      if (error.isNotEmpty) {
        return (false, isVerified);
      }
      if (token.isEmpty) {
        // Current Token is valid
        return (true, isVerified);
      }
      if (token.isNotEmpty) {
        // New Token Issued, update session
        final sessionManager = SessionManager();
        sessionManager.setAuthToken(token);
        return (true, isVerified);
      }
      return (true, isVerified);
    } on DioException catch (e) {
      return (false, false);
    }
  }
  static Future<bool> sendVerificationEmail() async {
    try {
      final response = await ApiClient.dio.get("/api/user/send_verification_email");
      final code = response.statusCode;
      if (code == 200) {
        return true;
      }
      return false;
    } on DioException catch (e) {
      return false;
    }
  }
  static Future<(bool success, String message)> verifyEmail(String otp) async {
    try {
      final response = await ApiClient.dio.post("/api/user/verify_email"
      , data: {
        "otp": otp,
      });
      final code = response.statusCode;
      final String message = response.data['message'] ?? "";
      final String error = response.data['message'] ?? "Could not verify email.";
      if (code == 200) {
        return (true, message);
      }
      return (false, error);
    } on DioException catch (e) {
        final String errorMessage =
            e.response?.data?['error'] ??
                "An error occurred during email verification.";
      return (false, errorMessage);
    }
  }
}