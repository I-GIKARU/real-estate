import 'package:flutter/material.dart';

class PropertyDetailsScreen extends StatefulWidget {
  const PropertyDetailsScreen({super.key});

  @override
  State<PropertyDetailsScreen> createState() => _PropertyDetailsScreenState();
}

class _PropertyDetailsScreenState extends State<PropertyDetailsScreen> {
  bool _isFavorite = false;
  
  @override
  Widget build(BuildContext context) {
    // Get property data from route arguments
    final property = ModalRoute.of(context)!.settings.arguments as Map<String, dynamic>;
    
    return Scaffold(
      body: CustomScrollView(
        slivers: [
          // App bar with property image
          SliverAppBar(
            expandedHeight: 250,
            pinned: true,
            flexibleSpace: FlexibleSpaceBar(
              background: Container(
                color: Colors.grey.shade300,
                child: const Center(
                  child: Icon(
                    Icons.home,
                    size: 80,
                    color: Colors.white,
                  ),
                ),
              ),
            ),
            actions: [
              // Favorite button
              IconButton(
                icon: Icon(
                  _isFavorite ? Icons.favorite : Icons.favorite_border,
                  color: _isFavorite ? Colors.red : null,
                ),
                onPressed: () {
                  setState(() {
                    _isFavorite = !_isFavorite;
                  });
                  
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text(
                        _isFavorite
                            ? 'Added to favorites'
                            : 'Removed from favorites',
                      ),
                      duration: const Duration(seconds: 1),
                    ),
                  );
                },
              ),
              // Share button
              IconButton(
                icon: const Icon(Icons.share),
                onPressed: () {
                  // Share property
                },
              ),
            ],
          ),
          
          // Property details
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Property name
                  Text(
                    property['name'],
                    style: const TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  
                  // Property address
                  Row(
                    children: [
                      const Icon(
                        Icons.location_on,
                        size: 16,
                        color: Colors.grey,
                      ),
                      const SizedBox(width: 4),
                      Expanded(
                        child: Text(
                          property['address'],
                          style: TextStyle(
                            color: Colors.grey.shade600,
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  
                  // Property price
                  Text(
                    property['price'],
                    style: TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                      color: Theme.of(context).colorScheme.primary,
                    ),
                  ),
                  const SizedBox(height: 24),
                  
                  // Property features
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceAround,
                    children: [
                      _buildFeatureItem(
                        context,
                        Icons.bed,
                        '${property['bedrooms']}',
                        'Bedrooms',
                      ),
                      _buildFeatureItem(
                        context,
                        Icons.bathtub,
                        '${property['bathrooms']}',
                        'Bathrooms',
                      ),
                      _buildFeatureItem(
                        context,
                        Icons.square_foot,
                        property['area'],
                        'Area',
                      ),
                    ],
                  ),
                  const SizedBox(height: 24),
                  
                  // Property description
                  const Text(
                    'Description',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    'This beautiful property offers modern living in a convenient location. '
                    'Featuring spacious rooms, updated appliances, and plenty of natural light. '
                    'The neighborhood provides easy access to shopping, dining, and public transportation. '
                    'Perfect for individuals or families looking for comfort and convenience.',
                    style: TextStyle(
                      color: Colors.grey.shade700,
                      height: 1.5,
                    ),
                  ),
                  const SizedBox(height: 24),
                  
                  // Property amenities
                  const Text(
                    'Amenities',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 16),
                  Wrap(
                    spacing: 16,
                    runSpacing: 16,
                    children: [
                      _buildAmenityItem(Icons.wifi, 'Wi-Fi'),
                      _buildAmenityItem(Icons.local_parking, 'Parking'),
                      _buildAmenityItem(Icons.pool, 'Swimming Pool'),
                      _buildAmenityItem(Icons.fitness_center, 'Gym'),
                      _buildAmenityItem(Icons.ac_unit, 'Air Conditioning'),
                      _buildAmenityItem(Icons.tv, 'TV'),
                      _buildAmenityItem(Icons.local_laundry_service, 'Laundry'),
                      _buildAmenityItem(Icons.security, 'Security'),
                    ],
                  ),
                  const SizedBox(height: 24),
                  
                  // Location map
                  const Text(
                    'Location',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 16),
                  Container(
                    height: 200,
                    decoration: BoxDecoration(
                      color: Colors.grey.shade300,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: const Center(
                      child: Icon(
                        Icons.map,
                        size: 50,
                        color: Colors.white,
                      ),
                    ),
                  ),
                  const SizedBox(height: 24),
                  
                  // Reviews
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text(
                        'Reviews',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      Row(
                        children: [
                          const Icon(
                            Icons.star,
                            color: Colors.amber,
                            size: 20,
                          ),
                          const SizedBox(width: 4),
                          Text(
                            '${property['rating']} (42 reviews)',
                            style: const TextStyle(
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  
                  // Review list (sample)
                  _buildReviewItem(
                    'John Doe',
                    4.5,
                    'Great property! Very clean and comfortable. The location is perfect for my needs.',
                    '2 weeks ago',
                  ),
                  const Divider(),
                  _buildReviewItem(
                    'Jane Smith',
                    5.0,
                    'Absolutely loved staying here. The amenities are top-notch and the host was very responsive.',
                    '1 month ago',
                  ),
                  const SizedBox(height: 16),
                  
                  // View all reviews button
                  Center(
                    child: TextButton(
                      onPressed: () {
                        // Navigate to all reviews
                      },
                      child: const Text('View All Reviews'),
                    ),
                  ),
                  const SizedBox(height: 100), // Space for bottom buttons
                ],
              ),
            ),
          ),
        ],
      ),
      
      // Bottom buttons
      bottomSheet: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        decoration: BoxDecoration(
          color: Theme.of(context).scaffoldBackgroundColor,
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.05),
              blurRadius: 10,
              offset: const Offset(0, -5),
            ),
          ],
        ),
        child: Row(
          children: [
            // Contact button
            Expanded(
              child: OutlinedButton.icon(
                onPressed: () {
                  // Contact owner/agent
                },
                icon: const Icon(Icons.phone),
                label: const Text('Contact'),
              ),
            ),
            const SizedBox(width: 16),
            // Book viewing button
            Expanded(
              child: ElevatedButton.icon(
                onPressed: () {
                  // Navigate to booking screen
                  Navigator.pushNamed(context, '/booking', arguments: property);
                },
                icon: const Icon(Icons.calendar_today),
                label: const Text('Book Viewing'),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFeatureItem(
    BuildContext context,
    IconData icon,
    String value,
    String label,
  ) {
    return Column(
      children: [
        Container(
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: Theme.of(context).colorScheme.primary.withOpacity(0.1),
            shape: BoxShape.circle,
          ),
          child: Icon(
            icon,
            color: Theme.of(context).colorScheme.primary,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          value,
          style: const TextStyle(
            fontWeight: FontWeight.bold,
            fontSize: 16,
          ),
        ),
        Text(
          label,
          style: TextStyle(
            color: Colors.grey.shade600,
            fontSize: 12,
          ),
        ),
      ],
    );
  }

  Widget _buildAmenityItem(IconData icon, String label) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(
        color: Colors.grey.shade100,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            icon,
            size: 16,
            color: Colors.grey.shade700,
          ),
          const SizedBox(width: 8),
          Text(
            label,
            style: TextStyle(
              color: Colors.grey.shade700,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildReviewItem(
    String name,
    double rating,
    String comment,
    String date,
  ) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                name,
                style: const TextStyle(
                  fontWeight: FontWeight.bold,
                ),
              ),
              Text(
                date,
                style: TextStyle(
                  color: Colors.grey.shade600,
                  fontSize: 12,
                ),
              ),
            ],
          ),
          const SizedBox(height: 4),
          Row(
            children: List.generate(5, (index) {
              return Icon(
                index < rating.floor() ? Icons.star : 
                (index == rating.floor() && rating % 1 > 0) ? Icons.star_half : Icons.star_border,
                color: Colors.amber,
                size: 16,
              );
            }),
          ),
          const SizedBox(height: 8),
          Text(
            comment,
            style: TextStyle(
              color: Colors.grey.shade700,
            ),
          ),
        ],
      ),
    );
  }
}
