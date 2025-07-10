class County {
  final String id;
  final String name;
  final String code;

  County({
    required this.id,
    required this.name,
    required this.code,
  });

  factory County.fromJson(Map<String, dynamic> json) {
    return County(
      id: json['id'] as String,
      name: json['name'] as String,
      code: json['code'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'code': code,
    };
  }
}

class SubCounty {
  final String id;
  final String name;
  final String countyId;
  final String? county;

  SubCounty({
    required this.id,
    required this.name,
    required this.countyId,
    this.county,
  });

  factory SubCounty.fromJson(Map<String, dynamic> json) {
    return SubCounty(
      id: json['id'] as String,
      name: json['name'] as String,
      countyId: json['county_id'] as String,
      county: json['county'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'county_id': countyId,
      'county': county,
    };
  }
}
