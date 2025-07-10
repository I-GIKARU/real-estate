import 'dart:convert';

class Property {
  final String id;
  final String title;
  final String description;
  final String propertyType;
  final double price;
  final String currency;
  final String paymentFrequency;
  final int bedrooms;
  final int bathrooms;
  final double? squareFootage;
  final String address;
  final String county;
  final String subCounty;
  final double latitude;
  final double longitude;
  final List<String> amenities;
  final List<String> utilities;
  final List<PropertyImage> images;
  final bool isAvailable;
  final DateTime? availableFrom;
  final String landlordId;
  final DateTime createdAt;
  final DateTime updatedAt;

  Property({
    required this.id,
    required this.title,
    required this.description,
    required this.propertyType,
    required this.price,
    required this.currency,
    required this.paymentFrequency,
    required this.bedrooms,
    required this.bathrooms,
    this.squareFootage,
    required this.address,
    required this.county,
    required this.subCounty,
    required this.latitude,
    required this.longitude,
    required this.amenities,
    required this.utilities,
    required this.images,
    required this.isAvailable,
    this.availableFrom,
    required this.landlordId,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Property.fromJson(Map<String, dynamic> json) {
    return Property(
      id: json['id'] as String,
      title: json['title'] as String,
      description: json['description'] as String,
      propertyType: json['property_type'] as String,
      price: (json['price'] as num).toDouble(),
      currency: json['currency'] as String,
      paymentFrequency: json['payment_frequency'] as String,
      bedrooms: json['bedrooms'] as int,
      bathrooms: json['bathrooms'] as int,
      squareFootage: json['square_footage'] != null ? (json['square_footage'] as num).toDouble() : null,
      address: json['address'] as String,
      county: json['county'] as String,
      subCounty: json['sub_county'] as String,
      latitude: (json['latitude'] as num).toDouble(),
      longitude: (json['longitude'] as num).toDouble(),
      amenities: List<String>.from(json['amenities'] ?? []),
      utilities: List<String>.from(json['utilities'] ?? []),
      images: (json['images'] as List?)?.map((e) => PropertyImage.fromJson(e)).toList() ?? [],
      isAvailable: json['is_available'] as bool,
      availableFrom: json['available_from'] != null ? DateTime.parse(json['available_from']) : null,
      landlordId: json['landlord_id'] as String,
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'description': description,
      'property_type': propertyType,
      'price': price,
      'currency': currency,
      'payment_frequency': paymentFrequency,
      'bedrooms': bedrooms,
      'bathrooms': bathrooms,
      'square_footage': squareFootage,
      'address': address,
      'county': county,
      'sub_county': subCounty,
      'latitude': latitude,
      'longitude': longitude,
      'amenities': amenities,
      'utilities': utilities,
      'images': images.map((e) => e.toJson()).toList(),
      'is_available': isAvailable,
      'available_from': availableFrom?.toIso8601String(),
      'landlord_id': landlordId,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}

class PropertyImage {
  final String id;
  final String url;
  final String publicId;
  final bool isPrimary;
  final int displayOrder;

  PropertyImage({
    required this.id,
    required this.url,
    required this.publicId,
    required this.isPrimary,
    required this.displayOrder,
  });

  factory PropertyImage.fromJson(Map<String, dynamic> json) {
    return PropertyImage(
      id: json['id'] as String,
      url: json['url'] as String,
      publicId: json['public_id'] as String,
      isPrimary: json['is_primary'] as bool,
      displayOrder: json['display_order'] as int,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'url': url,
      'public_id': publicId,
      'is_primary': isPrimary,
      'display_order': displayOrder,
    };
  }
}
