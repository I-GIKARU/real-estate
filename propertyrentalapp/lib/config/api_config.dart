import 'package:flutter_dotenv/flutter_dotenv.dart';

class ApiConfig {
  static String get baseUrl => dotenv.env['API_BASE_URL'] ?? ' ';
  static bool get isDevelopment => dotenv.env['DEBUG_MODE'] == 'true';
  static String get apiKey => dotenv.env['API_KEY'] ?? '';
  
  // Auth endpoints
  static const String register = '/register';
  static const String login = '/login';
  static const String refreshToken = '/refresh-token';
  static const String profile = '/profile';
  
  // Property endpoints
  static const String properties = '/properties';
  static const String myProperties = '/my-properties';
  
  // Location endpoints
  static const String counties = '/counties';
  static const String subCounties = '/sub-counties';
  
  // Email verification
  static const String sendVerificationEmail = '/send-verification-email';
  static const String verifyEmail = '/verify-email';
  static const String verificationStatus = '/verification-status';
  
  // Password reset
  static const String forgotPassword = '/auth/forgot-password';
  static const String resetPassword = '/auth/reset-password';
  static const String validateResetToken = '/auth/validate-reset-token';
  static const String changePassword = '/auth/change-password';
  
  // Kenyan features
  static const String amenities = '/amenities';
  static const String propertyTypes = '/property-types';
  static const String utilities = '/utilities';
  static const String rentalTerms = '/rental-terms';
  static const String popularAreas = '/popular-areas';
  static const String validatePhone = '/validate-phone';
  static const String formatCurrency = '/format-currency';
  
  // Headers
  static Map<String, String> get headers => {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
    if (apiKey.isNotEmpty) 'X-API-Key': apiKey,
  };
  
  static Map<String, String> getAuthHeaders(String token) => {
    ...headers,
    'Authorization': 'Bearer $token',
  };
}
