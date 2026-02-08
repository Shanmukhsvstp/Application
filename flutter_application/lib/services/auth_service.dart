import 'dart:ffi';

import 'package:flutter_application/api/api_service.dart';
import 'package:flutter_application/services/session_manager.dart';

class AuthService {

  static Future<(bool, String)> Login(String email, String password) async {
    if (email.isEmpty || password.isEmpty) {
      return (false, "Email and password cannot be empty.");
    }
    if (password.length < 8) {
      return (false, "Password must be at least 8 characters long.");
    }


    final result = await ApiService.login(email, password);

    if (result.$1) {
      final session = SessionManager();
      session.setAuthToken(result.$2);
      return (true, "Login successful.");
    } else {
      return (false, result.$3);
    }

  }

  static Future<(bool, String)> Signup(String username, String email, String password) async {
    if (username.isEmpty || email.isEmpty || password.isEmpty) {
      return (false, "Fields cannot be empty.");
    }
    if (password.length < 8) {
      return (false, "Password must be at least 8 characters long.");
    }


    final result = await ApiService.signup("username", "email", "password");

    if (result.$1) {
      final session = SessionManager();
      session.setAuthToken(result.$2);
      return (true, "Signup Successful!");
    } else {
      return (false, result.$3);
    }

  }

}