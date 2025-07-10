import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;
import '../config/api_config.dart';
import '../models/property_model.dart';
import 'auth_service.dart';

class PropertyService {
  final AuthService _authService = AuthService();

  // Get all public properties
  Future<List<Property>> getProperties({
    int page = 1,
    int limit = 20,
    String? county,
    String? subCounty,
    String? propertyType,
    double? minPrice,
    double? maxPrice,
    int? bedrooms,
    int? bathrooms,
  }) async {
    try {
      final queryParams = <String, String>{
        'page': page.toString(),
        'limit': limit.toString(),
      };

      if (county != null) queryParams['county'] = county;
      if (subCounty != null) queryParams['sub_county'] = subCounty;
      if (propertyType != null) queryParams['property_type'] = propertyType;
      if (minPrice != null) queryParams['min_price'] = minPrice.toString();
      if (maxPrice != null) queryParams['max_price'] = maxPrice.toString();
      if (bedrooms != null) queryParams['bedrooms'] = bedrooms.toString();
      if (bathrooms != null) queryParams['bathrooms'] = bathrooms.toString();

      final uri = Uri.parse('${ApiConfig.baseUrl}${ApiConfig.properties}')
          .replace(queryParameters: queryParams);

      final response = await http.get(
        uri,
        headers: ApiConfig.headers,
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        final List<dynamic> propertiesJson = data['data']['properties'];
        return propertiesJson.map((json) => Property.fromJson(json)).toList();
      }
      return [];
    } catch (e) {
      print('Get properties error: $e');
      return [];
    }
  }

  // Get single property by ID
  Future<Property?> getProperty(String id) async {
    try {
      final response = await http.get(
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.properties}/$id'),
        headers: ApiConfig.headers,
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        return Property.fromJson(data['data']);
      }
      return null;
    } catch (e) {
      print('Get property error: $e');
      return null;
    }
  }

  // Get landlord's properties
  Future<List<Property>> getMyProperties() async {
    try {
      final token = await _authService.getToken();
      if (token == null) return [];

      final response = await http.get(
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.myProperties}'),
        headers: ApiConfig.getAuthHeaders(token),
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        final List<dynamic> propertiesJson = data['data'];
        return propertiesJson.map((json) => Property.fromJson(json)).toList();
      }
      return [];
    } catch (e) {
      print('Get my properties error: $e');
      return [];
    }
  }

  // Create new property (landlord only)
  Future<Property?> createProperty({
    required String title,
    required String description,
    required String propertyType,
    required double price,
    required String currency,
    required String paymentFrequency,
    required int bedrooms,
    required int bathrooms,
    double? squareFootage,
    required String address,
    required String county,
    required String subCounty,
    required double latitude,
    required double longitude,
    required List<String> amenities,
    required List<String> utilities,
    DateTime? availableFrom,
  }) async {
    try {
      final token = await _authService.getToken();
      if (token == null) return null;

      final body = {
        'title': title,
        'description': description,
        'property_type': propertyType,
        'price': price,
        'currency': currency,
        'payment_frequency': paymentFrequency,
        'bedrooms': bedrooms,
        'bathrooms': bathrooms,
        'address': address,
        'county': county,
        'sub_county': subCounty,
        'latitude': latitude,
        'longitude': longitude,
        'amenities': amenities,
        'utilities': utilities,
      };

      if (squareFootage != null) body['square_footage'] = squareFootage;
      if (availableFrom != null) body['available_from'] = availableFrom.toIso8601String();

      final response = await http.post(
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.properties}'),
        headers: ApiConfig.getAuthHeaders(token),
        body: jsonEncode(body),
      );

      if (response.statusCode == 201) {
        final data = jsonDecode(response.body);
        return Property.fromJson(data['data']);
      }
      return null;
    } catch (e) {
      print('Create property error: $e');
      return null;
    }
  }

  // Update property (landlord only)
  Future<Property?> updateProperty({
    required String id,
    String? title,
    String? description,
    String? propertyType,
    double? price,
    String? currency,
    String? paymentFrequency,
    int? bedrooms,
    int? bathrooms,
    double? squareFootage,
    String? address,
    String? county,
    String? subCounty,
    double? latitude,
    double? longitude,
    List<String>? amenities,
    List<String>? utilities,
    bool? isAvailable,
    DateTime? availableFrom,
  }) async {
    try {
      final token = await _authService.getToken();
      if (token == null) return null;

      final body = <String, dynamic>{};
      if (title != null) body['title'] = title;
      if (description != null) body['description'] = description;
      if (propertyType != null) body['property_type'] = propertyType;
      if (price != null) body['price'] = price;
      if (currency != null) body['currency'] = currency;
      if (paymentFrequency != null) body['payment_frequency'] = paymentFrequency;
      if (bedrooms != null) body['bedrooms'] = bedrooms;
      if (bathrooms != null) body['bathrooms'] = bathrooms;
      if (squareFootage != null) body['square_footage'] = squareFootage;
      if (address != null) body['address'] = address;
      if (county != null) body['county'] = county;
      if (subCounty != null) body['sub_county'] = subCounty;
      if (latitude != null) body['latitude'] = latitude;
      if (longitude != null) body['longitude'] = longitude;
      if (amenities != null) body['amenities'] = amenities;
      if (utilities != null) body['utilities'] = utilities;
      if (isAvailable != null) body['is_available'] = isAvailable;
      if (availableFrom != null) body['available_from'] = availableFrom.toIso8601String();

      final response = await http.put(
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.properties}/$id'),
        headers: ApiConfig.getAuthHeaders(token),
        body: jsonEncode(body),
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        return Property.fromJson(data['data']);
      }
      return null;
    } catch (e) {
      print('Update property error: $e');
      return null;
    }
  }

  // Delete property (landlord only)
  Future<bool> deleteProperty(String id) async {
    try {
      final token = await _authService.getToken();
      if (token == null) return false;

      final response = await http.delete(
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.properties}/$id'),
        headers: ApiConfig.getAuthHeaders(token),
      );

      return response.statusCode == 200;
    } catch (e) {
      print('Delete property error: $e');
      return false;
    }
  }

  // Add property image
  Future<bool> addPropertyImage(String propertyId, File imageFile) async {
    try {
      final token = await _authService.getToken();
      if (token == null) return false;

      final request = http.MultipartRequest(
        'POST',
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.properties}/$propertyId/images'),
      );

      request.headers.addAll(ApiConfig.getAuthHeaders(token));
      request.files.add(await http.MultipartFile.fromPath('image', imageFile.path));

      final response = await request.send();
      return response.statusCode == 201;
    } catch (e) {
      print('Add property image error: $e');
      return false;
    }
  }

  // Delete property image
  Future<bool> deletePropertyImage(String propertyId, String imageId) async {
    try {
      final token = await _authService.getToken();
      if (token == null) return false;

      final response = await http.delete(
        Uri.parse('${ApiConfig.baseUrl}${ApiConfig.properties}/$propertyId/images/$imageId'),
        headers: ApiConfig.getAuthHeaders(token),
      );

      return response.statusCode == 200;
    } catch (e) {
      print('Delete property image error: $e');
      return false;
    }
  }
}
