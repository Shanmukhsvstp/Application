import 'package:shared_preferences/shared_preferences.dart';

class SessionManager {
  static final SessionManager _instance = SessionManager._internal();
  static const authTokenKey = "auth_token_key";

  factory SessionManager() {
    return _instance;
  }

  SessionManager._internal();


  void setAuthToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(authTokenKey, token);
  }


  Future<String?> getAuthToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(authTokenKey);
  }

  void clearSession() {
    SharedPreferences.getInstance().then((prefs) {
      prefs.remove(authTokenKey);
    });
  }
}