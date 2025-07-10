import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiService {
  // Base URL for the backend API
  static const String baseUrl = 'https://api.example.com/api/v1';
  
  // Auth token storage
  static String? _authToken;
  
  // Set auth token
  static void setAuthToken(String token) {
    _authToken = token;
  }
  
  // Get auth token
  static String? getAuthToken() {
    return _authToken;
  }
  
  // Clear auth token (logout)
  static void clearAuthToken() {
    _authToken = null;
  }
  
  // Headers with auth token
  static Map<String, String> get _headers {
    final headers = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };
    
    if (_authToken != null) {
      headers['Authorization'] = 'Bearer $_authToken';
    }
    
    return headers;
  }
  
  // Login
  static Future<Map<String, dynamic>> login(String email, String password) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/auth/login'),
        headers: _headers,
        body: jsonEncode({
          'email': email,
          'password': password,
        }),
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 200) {
        // Save auth token
        setAuthToken(data['token']);
        return data;
      } else {
        throw Exception(data['message'] ?? 'Login failed');
      }
    } catch (e) {
      throw Exception('Login failed: $e');
    }
  }
  
  // Register
  static Future<Map<String, dynamic>> register(Map<String, dynamic> userData) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/auth/register'),
        headers: _headers,
        body: jsonEncode(userData),
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 201) {
        // Save auth token
        setAuthToken(data['token']);
        return data;
      } else {
        throw Exception(data['message'] ?? 'Registration failed');
      }
    } catch (e) {
      throw Exception('Registration failed: $e');
    }
  }
  
  // Get user profile
  static Future<Map<String, dynamic>> getUserProfile() async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/user/profile'),
        headers: _headers,
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 200) {
        return data;
      } else {
        throw Exception(data['message'] ?? 'Failed to get user profile');
      }
    } catch (e) {
      throw Exception('Failed to get user profile: $e');
    }
  }
  
  // Update user profile
  static Future<Map<String, dynamic>> updateUserProfile(Map<String, dynamic> userData) async {
    try {
      final response = await http.put(
        Uri.parse('$baseUrl/user/profile'),
        headers: _headers,
        body: jsonEncode(userData),
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 200) {
        return data;
      } else {
        throw Exception(data['message'] ?? 'Failed to update profile');
      }
    } catch (e) {
      throw Exception('Failed to update profile: $e');
    }
  }
  
  // Get properties
  static Future<List<dynamic>> getProperties({
    String? location,
    String? type,
    String? priceRange,
    String? sort,
    int page = 1,
    int limit = 10,
  }) async {
    try {
      final queryParams = {
        'page': page.toString(),
        'limit': limit.toString(),
      };
      
      if (location != null) queryParams['location'] = location;
      if (type != null) queryParams['type'] = type;
      if (priceRange != null) queryParams['price_range'] = priceRange;
      if (sort != null) queryParams['sort'] = sort;
      
      final uri = Uri.parse('$baseUrl/properties').replace(queryParameters: queryParams);
      
      final response = await http.get(
        uri,
        headers: _headers,
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 200) {
        return data['properties'];
      } else {
        throw Exception(data['message'] ?? 'Failed to get properties');
      }
    } catch (e) {
      throw Exception('Failed to get properties: $e');
    }
  }
  
  // Get property details
  static Future<Map<String, dynamic>> getPropertyDetails(String propertyId) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/properties/$propertyId'),
        headers: _headers,
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 200) {
        return data;
      } else {
        throw Exception(data['message'] ?? 'Failed to get property details');
      }
    } catch (e) {
      throw Exception('Failed to get property details: $e');
    }
  }
  
  // Save property to favorites
  static Future<void> saveProperty(String propertyId) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/favorites'),
        headers: _headers,
        body: jsonEncode({
          'property_id': propertyId,
        }),
      );
      
      if (response.statusCode != 200 && response.statusCode != 201) {
        final data = jsonDecode(response.body);
        throw Exception(data['message'] ?? 'Failed to save property');
      }
    } catch (e) {
      throw Exception('Failed to save property: $e');
    }
  }
  
  // Remove property from favorites
  static Future<void> unsaveProperty(String propertyId) async {
    try {
      final response = await http.delete(
        Uri.parse('$baseUrl/favorites/$propertyId'),
        headers: _headers,
      );
      
      if (response.statusCode != 200) {
        final data = jsonDecode(response.body);
        throw Exception(data['message'] ?? 'Failed to remove property from favorites');
      }
    } catch (e) {
      throw Exception('Failed to remove property from favorites: $e');
    }
  }
  
  // Get saved properties
  static Future<List<dynamic>> getSavedProperties() async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/favorites'),
        headers: _headers,
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 200) {
        return data['properties'];
      } else {
        throw Exception(data['message'] ?? 'Failed to get saved properties');
      }
    } catch (e) {
      throw Exception('Failed to get saved properties: $e');
    }
  }
  
  // Create booking
  static Future<Map<String, dynamic>> createBooking(Map<String, dynamic> bookingData) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/bookings'),
        headers: _headers,
        body: jsonEncode(bookingData),
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 201) {
        return data;
      } else {
        throw Exception(data['message'] ?? 'Failed to create booking');
      }
    } catch (e) {
      throw Exception('Failed to create booking: $e');
    }
  }
  
  // Get bookings
  static Future<List<dynamic>> getBookings() async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/bookings'),
        headers: _headers,
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 200) {
        return data['bookings'];
      } else {
        throw Exception(data['message'] ?? 'Failed to get bookings');
      }
    } catch (e) {
      throw Exception('Failed to get bookings: $e');
    }
  }
  
  // Get booking details
  static Future<Map<String, dynamic>> getBookingDetails(String bookingId) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/bookings/$bookingId'),
        headers: _headers,
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 200) {
        return data;
      } else {
        throw Exception(data['message'] ?? 'Failed to get booking details');
      }
    } catch (e) {
      throw Exception('Failed to get booking details: $e');
    }
  }
  
  // Cancel booking
  static Future<void> cancelBooking(String bookingId) async {
    try {
      final response = await http.put(
        Uri.parse('$baseUrl/bookings/$bookingId/cancel'),
        headers: _headers,
      );
      
      if (response.statusCode != 200) {
        final data = jsonDecode(response.body);
        throw Exception(data['message'] ?? 'Failed to cancel booking');
      }
    } catch (e) {
      throw Exception('Failed to cancel booking: $e');
    }
  }
  
  // Process payment
  static Future<Map<String, dynamic>> processPayment(Map<String, dynamic> paymentData) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/payments'),
        headers: _headers,
        body: jsonEncode(paymentData),
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 200 || response.statusCode == 201) {
        return data;
      } else {
        throw Exception(data['message'] ?? 'Payment failed');
      }
    } catch (e) {
      throw Exception('Payment failed: $e');
    }
  }
  
  // Get payment methods
  static Future<List<dynamic>> getPaymentMethods() async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/payment-methods'),
        headers: _headers,
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 200) {
        return data['payment_methods'];
      } else {
        throw Exception(data['message'] ?? 'Failed to get payment methods');
      }
    } catch (e) {
      throw Exception('Failed to get payment methods: $e');
    }
  }
  
  // Add payment method
  static Future<Map<String, dynamic>> addPaymentMethod(Map<String, dynamic> paymentMethodData) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/payment-methods'),
        headers: _headers,
        body: jsonEncode(paymentMethodData),
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 201) {
        return data;
      } else {
        throw Exception(data['message'] ?? 'Failed to add payment method');
      }
    } catch (e) {
      throw Exception('Failed to add payment method: $e');
    }
  }
  
  // Delete payment method
  static Future<void> deletePaymentMethod(String paymentMethodId) async {
    try {
      final response = await http.delete(
        Uri.parse('$baseUrl/payment-methods/$paymentMethodId'),
        headers: _headers,
      );
      
      if (response.statusCode != 200) {
        final data = jsonDecode(response.body);
        throw Exception(data['message'] ?? 'Failed to delete payment method');
      }
    } catch (e) {
      throw Exception('Failed to delete payment method: $e');
    }
  }
  
  // Submit review
  static Future<Map<String, dynamic>> submitReview(Map<String, dynamic> reviewData) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/reviews'),
        headers: _headers,
        body: jsonEncode(reviewData),
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 201) {
        return data;
      } else {
        throw Exception(data['message'] ?? 'Failed to submit review');
      }
    } catch (e) {
      throw Exception('Failed to submit review: $e');
    }
  }
  
  // Get reviews for property
  static Future<List<dynamic>> getPropertyReviews(String propertyId) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/properties/$propertyId/reviews'),
        headers: _headers,
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 200) {
        return data['reviews'];
      } else {
        throw Exception(data['message'] ?? 'Failed to get reviews');
      }
    } catch (e) {
      throw Exception('Failed to get reviews: $e');
    }
  }
  
  // Get user reviews
  static Future<List<dynamic>> getUserReviews() async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/user/reviews'),
        headers: _headers,
      );
      
      final data = jsonDecode(response.body);
      
      if (response.statusCode == 200) {
        return data['reviews'];
      } else {
        throw Exception(data['message'] ?? 'Failed to get user reviews');
      }
    } catch (e) {
      throw Exception('Failed to get user reviews: $e');
    }
  }
}
