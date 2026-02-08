import 'package:dio/dio.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_application/api/api_client.dart';
import 'package:flutter_application/api/api_service.dart';
import 'package:flutter_application/home_page.dart';
import 'package:flutter_application/login_screen.dart';
import 'package:flutter_application/services/session_manager.dart';

class AppStartup extends StatefulWidget {
  const AppStartup({super.key});

  @override
  State<AppStartup> createState() => _AppStartupState();
}

class _AppStartupState extends State<AppStartup> {
  @override
  @override
  void initState() {
    super.initState();
    validateAuth();
  }

  Future<void> validateAuth() async {
    final token = await SessionManager().getAuthToken();

    ApiClient.dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) {
          options.headers["Authorization"] = "Bearer $token";
          return handler.next(options);
        },
      ),
    );

    if (!mounted) return;

    if (token != null && token.isNotEmpty) {
      final result = await ApiService.validateAuth();
      if (result.$1) {
        if (!result.$2) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text("Account not verified. Please check your email.")),
            );
            Navigator.pushReplacement(
              context,
              MaterialPageRoute(builder: (_) => const HomePage()),
            );
        }
        else {
          Navigator.pushReplacement(
            context,
            MaterialPageRoute(builder: (_) => const HomePage()),
          );
        }
      }
      else {
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(builder: (_) => const LoginScreen()),
        );
      }
    } else {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(builder: (_) => const LoginScreen()),
      );
    }

  }

  @override
  Widget build(BuildContext context) {
    return const Scaffold (
      body: Center(
        child: CircularProgressIndicator.adaptive(),
      ),
    );
  }
}