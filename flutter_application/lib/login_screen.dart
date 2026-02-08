import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter_application/app_constants.dart';
import 'package:flutter_application/services/auth_service.dart';
import 'package:flutter_application/services/loader_service.dart';
import 'package:flutter_application/signup_screen.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

  class _LoginScreenState extends State<LoginScreen> {
  bool _obscurePassword = true;
  double borderRadius = AppConstants.defaultBorderRadius;

  final emailController = TextEditingController();
  final passwordController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text(
              "Login",
              style: TextStyle(fontSize: 28, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 32),

            TextField(
              controller: emailController,
              decoration: InputDecoration(
                labelText: "Username or Email",
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(borderRadius),
                ),
              ),
            ),
            const SizedBox(height: 16),

            TextField(
              controller: passwordController,
              obscureText: _obscurePassword,
              decoration: InputDecoration(
                labelText: "Password",
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(borderRadius),
                ),
                suffixIcon: IconButton(onPressed: () {
                  setState(() {
                    _obscurePassword = !_obscurePassword;
                  });
                }, icon: Icon(_obscurePassword ? Icons.visibility_off : Icons.visibility))
              ),
            ),
            const SizedBox(height: 24),

            SizedBox(
              width: double.infinity,
              height: 50,
              child: ElevatedButton(
                onPressed: () async {

                  LoaderService.show(context, "Logging in...");

                  final String email = emailController.text.trim();
                  final String password = passwordController.text.trim();

                  final result = await AuthService.Login(email, password);

                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text(result.$2))
                  );

                  if (result.$1) {
                    // Navigate to the home screen or dashboard
                    // Navigator.pushReplacement(context, MaterialPageRoute(builder: (_) => HomeScreen()));
                  }
                  LoaderService.hide(context);
                },
                child: const Text("Login"),
              ),
            ),

            const SizedBox(height: 24),

            RichText(
              text: TextSpan(
                text: "Don't have an account? ",
                style: TextStyle(fontSize: 16, color: Colors.black),
                children: [
                  TextSpan(
                    text: "Sign up",
                    style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Colors.blue),
                    recognizer: TapGestureRecognizer()..onTap = () {
                      Navigator.push(context,
                        MaterialPageRoute(
                          builder: (context) => const SignupScreen()
                        )
                      );
                    },
                  )
                ],
              ),
            )
          ],
        ),
      ),
    );
  }
}
