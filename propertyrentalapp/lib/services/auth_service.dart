import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import '../config/api_config.dart';
import '../models/user_model.dart';

class AuthService {
  static const String _tokenKey = 'auth_token';
  static const String _refreshTokenKey = 'refresh_token';
  static const String _userKey = 'user_data';

  // Login with email and password
  Future<AuthResponse?> login(String email, String password) async {
    try {
      print('Attempting login with URL: ${ApiConfig.baseUrl}${ApiConfig.login}');
      
      final response = await http.post(
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.login}'),
        headers: ApiConfig.headers,
        body: jsonEncode({
          'email': email,
          'password': password,
        }),
      );

      print('Login response status: ${response.statusCode}');
      print('Login response body: ${response.body}');

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        // Backend returns user and token directly, not nested in 'data'
        final authResponse = AuthResponse(
          token: data['token'],
          refreshToken: '', // Backend doesn't provide refresh token yet
          user: User.fromJson(data['user']),
        );
        
        await _saveAuthData(authResponse);
        return authResponse;
      } else {
        final error = jsonDecode(response.body);
        throw Exception(error['error'] ?? 'Login failed');
      }
    } catch (e) {
      print('Login error: $e');
      if (e.toString().contains('SocketException') || e.toString().contains('Connection refused')) {
        throw Exception('Unable to connect to server. Please check your internet connection.');
      }
      rethrow;
    }
  }

  // Register new user
  Future<AuthResponse?> register({
    required String email,
    required String password,
    required String firstName,
    required String lastName,
    required String phoneNumber,
    required String userType,
  }) async {
    try {
      print('Attempting registration with URL: ${ApiConfig.baseUrl}${ApiConfig.register}');
      print('Request body: ${jsonEncode({
        'email': email,
        'password': password,
        'first_name': firstName,
        'last_name': lastName,
        'phone_number': phoneNumber,
        'user_type': userType,
      })}');
      
      final response = await http.post(
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.register}'),
        headers: ApiConfig.headers,
        body: jsonEncode({
          'email': email,
          'password': password,
          'first_name': firstName,
          'last_name': lastName,
          'phone_number': phoneNumber,
          'user_type': userType,
        }),
      );

      print('Response status: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 201) {
        final data = jsonDecode(response.body);
        // Backend returns user and token directly, not nested in 'data'
        final authResponse = AuthResponse(
          token: data['token'],
          refreshToken: '', // Backend doesn't provide refresh token yet
          user: User.fromJson(data['user']),
        );
        
        await _saveAuthData(authResponse);
        return authResponse;
      } else {
        final error = jsonDecode(response.body);
        throw Exception(error['error'] ?? 'Registration failed');
      }
    } catch (e) {
      print('Registration error: $e');
      if (e.toString().contains('SocketException') || e.toString().contains('Connection refused')) {
        throw Exception('Unable to connect to server. Please check your internet connection.');
      }
      rethrow;
    }
  }

  // Get current user profile
  Future<User?> getProfile() async {
    try {
      final token = await getToken();
      if (token == null) return null;

      final response = await http.get(
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.profile}'),
        headers: ApiConfig.getAuthHeaders(token),
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        // Backend returns user directly, not nested in 'data'
        return User.fromJson(data['user']);
      }
      return null;
    } catch (e) {
      print('Get profile error: $e');
      return null;
    }
  }

  // Update user profile
  Future<User?> updateProfile({
    String? firstName,
    String? lastName,
    String? phoneNumber,
  }) async {
    try {
      final token = await getToken();
      if (token == null) return null;

      final body = <String, dynamic>{};
      if (firstName != null) body['first_name'] = firstName;
      if (lastName != null) body['last_name'] = lastName;
      if (phoneNumber != null) body['phone_number'] = phoneNumber;

      final response = await http.put(
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.profile}'),
        headers: ApiConfig.getAuthHeaders(token),
        body: jsonEncode(body),
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        final user = User.fromJson(data['data']);
        await _saveUser(user);
        return user;
      }
      return null;
    } catch (e) {
      print('Update profile error: $e');
      return null;
    }
  }

  // Refresh token
  Future<String?> refreshToken() async {
    try {
      final refreshToken = await getRefreshToken();
      if (refreshToken == null) return null;

      final response = await http.post(
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.refreshToken}'),
        headers: ApiConfig.headers,
        body: jsonEncode({
          'refresh_token': refreshToken,
        }),
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        final newToken = data['data']['token'];
        await _saveToken(newToken);
        return newToken;
      }
      return null;
    } catch (e) {
      print('Refresh token error: $e');
      return null;
    }
  }

  // Send email verification
  Future<bool> sendVerificationEmail() async {
    try {
      final token = await getToken();
      if (token == null) return false;

      final response = await http.post(
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.sendVerificationEmail}'),
        headers: ApiConfig.getAuthHeaders(token),
      );

      return response.statusCode == 200;
    } catch (e) {
      print('Send verification email error: $e');
      return false;
    }
  }

  // Verify email with token
  Future<bool> verifyEmail(String verificationToken) async {
    try {
      final response = await http.post(
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.verifyEmail}'),
        headers: ApiConfig.headers,
        body: jsonEncode({
          'token': verificationToken,
        }),
      );

      return response.statusCode == 200;
    } catch (e) {
      print('Verify email error: $e');
      return false;
    }
  }

  // Check verification status
  Future<bool> getVerificationStatus() async {
    try {
      final token = await getToken();
      if (token == null) return false;

      final response = await http.get(
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.verificationStatus}'),
        headers: ApiConfig.getAuthHeaders(token),
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        return data['data']['is_verified'] ?? false;
      }
      return false;
    } catch (e) {
      print('Get verification status error: $e');
      return false;
    }
  }

  // Sign out
  Future<void> signOut() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_tokenKey);
    await prefs.remove(_refreshTokenKey);
    await prefs.remove(_userKey);
  }

  // Get stored token
  Future<String?> getToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_tokenKey);
  }

  // Get stored refresh token
  Future<String?> getRefreshToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_refreshTokenKey);
  }

  // Get stored user
  Future<User?> getStoredUser() async {
    final prefs = await SharedPreferences.getInstance();
    final userData = prefs.getString(_userKey);
    if (userData != null) {
      return User.fromJson(jsonDecode(userData));
    }
    return null;
  }

  // Check if user is authenticated
  Future<bool> isAuthenticated() async {
    final token = await getToken();
    return token != null;
  }

  // Private methods
  Future<void> _saveAuthData(AuthResponse authResponse) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_tokenKey, authResponse.token);
    await prefs.setString(_refreshTokenKey, authResponse.refreshToken);
    await prefs.setString(_userKey, jsonEncode(authResponse.user.toJson()));
  }

  Future<void> _saveToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_tokenKey, token);
  }

  Future<void> _saveUser(User user) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_userKey, jsonEncode(user.toJson()));
  }
}
