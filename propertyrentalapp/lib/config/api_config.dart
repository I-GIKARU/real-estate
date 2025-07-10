class ApiConfig {
  static const String baseUrl = 'https://kenyan-real-estate-backend-671327858247.us-central1.run.app/api/v1';
  static const bool isDevelopment = true;
  
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
  };
  
  static Map<String, String> getAuthHeaders(String token) => {
    ...headers,
    'Authorization': 'Bearer $token',
  };
}
