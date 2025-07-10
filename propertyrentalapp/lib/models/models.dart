import 'package:flutter/material.dart';

// User model
class User {
  final String id;
  final String name;
  final String email;
  final String phone;
  final String? profileImage;

  User({
    required this.id,
    required this.name,
    required this.email,
    required this.phone,
    this.profileImage,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'],
      name: json['name'],
      email: json['email'],
      phone: json['phone'],
      profileImage: json['profile_image'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'email': email,
      'phone': phone,
      'profile_image': profileImage,
    };
  }
}

// Property model
class Property {
  final String id;
  final String name;
  final String address;
  final String description;
  final double price;
  final String priceUnit; // e.g., "month", "week", "day"
  final String type; // e.g., "apartment", "house", "villa"
  final int bedrooms;
  final int bathrooms;
  final String area;
  final List<String> images;
  final List<String> amenities;
  final double rating;
  final int reviewCount;
  final bool isFeatured;
  final Map<String, dynamic>? location; // lat, lng

  Property({
    required this.id,
    required this.name,
    required this.address,
    required this.description,
    required this.price,
    required this.priceUnit,
    required this.type,
    required this.bedrooms,
    required this.bathrooms,
    required this.area,
    required this.images,
    required this.amenities,
    required this.rating,
    required this.reviewCount,
    required this.isFeatured,
    this.location,
  });

  factory Property.fromJson(Map<String, dynamic> json) {
    return Property(
      id: json['id'],
      name: json['name'],
      address: json['address'],
      description: json['description'],
      price: json['price'].toDouble(),
      priceUnit: json['price_unit'],
      type: json['type'],
      bedrooms: json['bedrooms'],
      bathrooms: json['bathrooms'],
      area: json['area'],
      images: List<String>.from(json['images']),
      amenities: List<String>.from(json['amenities']),
      rating: json['rating'].toDouble(),
      reviewCount: json['review_count'],
      isFeatured: json['is_featured'],
      location: json['location'],
    );
  }
}

// Booking model
class Booking {
  final String id;
  final Property property;
  final DateTime date;
  final String time;
  final String duration;
  final String status; // e.g., "pending", "confirmed", "cancelled"
  final DateTime createdAt;

  Booking({
    required this.id,
    required this.property,
    required this.date,
    required this.time,
    required this.duration,
    required this.status,
    required this.createdAt,
  });

  factory Booking.fromJson(Map<String, dynamic> json) {
    return Booking(
      id: json['id'],
      property: Property.fromJson(json['property']),
      date: DateTime.parse(json['date']),
      time: json['time'],
      duration: json['duration'],
      status: json['status'],
      createdAt: DateTime.parse(json['created_at']),
    );
  }
}

// Payment model
class Payment {
  final String id;
  final String bookingId;
  final double amount;
  final String status; // e.g., "pending", "completed", "failed"
  final String method; // e.g., "credit_card", "paypal"
  final DateTime date;

  Payment({
    required this.id,
    required this.bookingId,
    required this.amount,
    required this.status,
    required this.method,
    required this.date,
  });

  factory Payment.fromJson(Map<String, dynamic> json) {
    return Payment(
      id: json['id'],
      bookingId: json['booking_id'],
      amount: json['amount'].toDouble(),
      status: json['status'],
      method: json['method'],
      date: DateTime.parse(json['date']),
    );
  }
}

// Review model
class Review {
  final String id;
  final String propertyId;
  final String userId;
  final String userName;
  final String? userImage;
  final double rating;
  final String comment;
  final DateTime date;

  Review({
    required this.id,
    required this.propertyId,
    required this.userId,
    required this.userName,
    this.userImage,
    required this.rating,
    required this.comment,
    required this.date,
  });

  factory Review.fromJson(Map<String, dynamic> json) {
    return Review(
      id: json['id'],
      propertyId: json['property_id'],
      userId: json['user_id'],
      userName: json['user_name'],
      userImage: json['user_image'],
      rating: json['rating'].toDouble(),
      comment: json['comment'],
      date: DateTime.parse(json['date']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'property_id': propertyId,
      'user_id': userId,
      'user_name': userName,
      'user_image': userImage,
      'rating': rating,
      'comment': comment,
      'date': date.toIso8601String(),
    };
  }
}

// Payment Method model
class PaymentMethod {
  final String id;
  final String type; // e.g., "credit_card", "paypal"
  final String? cardNumber;
  final String? cardHolderName;
  final String? expiryDate;
  final String? email; // For PayPal
  final bool isDefault;

  PaymentMethod({
    required this.id,
    required this.type,
    this.cardNumber,
    this.cardHolderName,
    this.expiryDate,
    this.email,
    required this.isDefault,
  });

  factory PaymentMethod.fromJson(Map<String, dynamic> json) {
    return PaymentMethod(
      id: json['id'],
      type: json['type'],
      cardNumber: json['card_number'],
      cardHolderName: json['card_holder_name'],
      expiryDate: json['expiry_date'],
      email: json['email'],
      isDefault: json['is_default'],
    );
  }
}