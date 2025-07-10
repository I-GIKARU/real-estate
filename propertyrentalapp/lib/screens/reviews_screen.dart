import 'package:flutter/material.dart';

class ReviewsScreen extends StatefulWidget {
  const ReviewsScreen({super.key});

  @override
  State<ReviewsScreen> createState() => _ReviewsScreenState();
}

class _ReviewsScreenState extends State<ReviewsScreen> {
  // Sample reviews data
  final List<Map<String, dynamic>> _reviews = [
    {
      'id': '1',
      'userName': 'John Doe',
      'userImage': '',
      'rating': 4.5,
      'comment': 'Great property! Very clean and comfortable. The location is perfect for my needs.',
      'date': '2 weeks ago',
      'propertyName': 'Modern Apartment in Downtown',
    },
    {
      'id': '2',
      'userName': 'Jane Smith',
      'userImage': '',
      'rating': 5.0,
      'comment': 'Absolutely loved staying here. The amenities are top-notch and the host was very responsive.',
      'date': '1 month ago',
      'propertyName': 'Luxury Villa with Pool',
    },
    {
      'id': '3',
      'userName': 'Robert Johnson',
      'userImage': '',
      'rating': 3.5,
      'comment': 'The property was nice but there were some issues with the plumbing. The host was quick to address them though.',
      'date': '2 months ago',
      'propertyName': 'Cozy Studio near University',
    },
    {
      'id': '4',
      'userName': 'Emily Wilson',
      'userImage': '',
      'rating': 4.0,
      'comment': 'Good location and comfortable space. Would recommend for short stays.',
      'date': '3 months ago',
      'propertyName': 'Family Home with Garden',
    },
    {
      'id': '5',
      'userName': 'Michael Brown',
      'userImage': '',
      'rating': 5.0,
      'comment': 'Exceptional property with amazing views. Everything was perfect from check-in to check-out.',
      'date': '3 months ago',
      'propertyName': 'Penthouse with City View',
    },
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('My Reviews'),
      ),
      body: _reviews.isEmpty
          ? _buildEmptyState()
          : ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: _reviews.length,
              itemBuilder: (context, index) {
                final review = _reviews[index];
                return _buildReviewCard(review);
              },
            ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          // Navigate to write review screen
          _showWriteReviewDialog();
        },
        child: const Icon(Icons.rate_review),
      ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.rate_review_outlined,
            size: 80,
            color: Colors.grey.shade400,
          ),
          const SizedBox(height: 16),
          Text(
            'No Reviews Yet',
            style: Theme.of(context).textTheme.titleLarge,
          ),
          const SizedBox(height: 8),
          Text(
            'Your reviews will appear here',
            style: TextStyle(
              color: Colors.grey.shade600,
            ),
          ),
          const SizedBox(height: 24),
          ElevatedButton.icon(
            onPressed: () {
              // Navigate to write review screen
              _showWriteReviewDialog();
            },
            icon: const Icon(Icons.rate_review),
            label: const Text('Write a Review'),
          ),
        ],
      ),
    );
  }

  Widget _buildReviewCard(Map<String, dynamic> review) {
    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
      ),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Property name
            Text(
              review['propertyName'],
              style: const TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 16,
              ),
            ),
            const SizedBox(height: 12),
            
            // Review header
            Row(
              children: [
                // User avatar
                CircleAvatar(
                  radius: 20,
                  backgroundColor: Theme.of(context).colorScheme.primary,
                  child: const Icon(
                    Icons.person,
                    color: Colors.white,
                  ),
                ),
                const SizedBox(width: 12),
                // User name and date
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        review['userName'],
                        style: const TextStyle(
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      Text(
                        review['date'],
                        style: TextStyle(
                          color: Colors.grey.shade600,
                          fontSize: 12,
                        ),
                      ),
                    ],
                  ),
                ),
                // Rating
                Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 8,
                    vertical: 4,
                  ),
                  decoration: BoxDecoration(
                    color: Theme.of(context).colorScheme.primary,
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Row(
                    children: [
                      const Icon(
                        Icons.star,
                        color: Colors.white,
                        size: 16,
                      ),
                      const SizedBox(width: 4),
                      Text(
                        review['rating'].toString(),
                        style: const TextStyle(
                          color: Colors.white,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            
            // Review comment
            Text(
              review['comment'],
              style: TextStyle(
                color: Colors.grey.shade700,
              ),
            ),
            const SizedBox(height: 16),
            
            // Actions
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                TextButton.icon(
                  onPressed: () {
                    // Edit review
                    _showWriteReviewDialog(review: review);
                  },
                  icon: const Icon(Icons.edit, size: 16),
                  label: const Text('Edit'),
                ),
                const SizedBox(width: 8),
                TextButton.icon(
                  onPressed: () {
                    // Delete review
                    _showDeleteConfirmation(review);
                  },
                  icon: const Icon(Icons.delete, size: 16),
                  label: const Text('Delete'),
                  style: TextButton.styleFrom(
                    foregroundColor: Colors.red,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  void _showWriteReviewDialog({Map<String, dynamic>? review}) {
    final _ratingController = TextEditingController(
      text: review != null ? review['rating'].toString() : '5.0',
    );
    final _commentController = TextEditingController(
      text: review != null ? review['comment'] : '',
    );
    double _rating = review != null ? review['rating'] : 5.0;

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(review != null ? 'Edit Review' : 'Write a Review'),
        content: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              if (review == null) ...[
                // Property selection (only for new reviews)
                const Text(
                  'Select Property',
                  style: TextStyle(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),
                DropdownButtonFormField<String>(
                  decoration: const InputDecoration(
                    border: OutlineInputBorder(),
                  ),
                  value: 'Modern Apartment in Downtown',
                  items: const [
                    DropdownMenuItem(
                      value: 'Modern Apartment in Downtown',
                      child: Text('Modern Apartment in Downtown'),
                    ),
                    DropdownMenuItem(
                      value: 'Luxury Villa with Pool',
                      child: Text('Luxury Villa with Pool'),
                    ),
                    DropdownMenuItem(
                      value: 'Cozy Studio near University',
                      child: Text('Cozy Studio near University'),
                    ),
                  ],
                  onChanged: (value) {
                    // Handle property selection
                  },
                ),
                const SizedBox(height: 16),
              ],
              
              // Rating
              const Text(
                'Rating',
                style: TextStyle(
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 8),
              StatefulBuilder(
                builder: (context, setState) {
                  return Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: List.generate(5, (index) {
                      return IconButton(
                        icon: Icon(
                          index < _rating.floor() ? Icons.star : 
                          (index == _rating.floor() && _rating % 1 > 0) ? Icons.star_half : Icons.star_border,
                          color: Colors.amber,
                        ),
                        onPressed: () {
                          setState(() {
                            _rating = index + 1.0;
                            _ratingController.text = _rating.toString();
                          });
                        },
                      );
                    }),
                  );
                },
              ),
              const SizedBox(height: 16),
              
              // Comment
              const Text(
                'Comment',
                style: TextStyle(
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 8),
              TextField(
                controller: _commentController,
                maxLines: 5,
                decoration: const InputDecoration(
                  hintText: 'Share your experience...',
                  border: OutlineInputBorder(),
                ),
              ),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () {
              Navigator.pop(context);
            },
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              // Save review
              Navigator.pop(context);
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text(
                    review != null ? 'Review updated successfully' : 'Review submitted successfully',
                  ),
                  backgroundColor: Colors.green,
                ),
              );
            },
            child: Text(review != null ? 'Update' : 'Submit'),
          ),
        ],
      ),
    );
  }

  void _showDeleteConfirmation(Map<String, dynamic> review) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Review'),
        content: const Text(
          'Are you sure you want to delete this review? This action cannot be undone.',
        ),
        actions: [
          TextButton(
            onPressed: () {
              Navigator.pop(context);
            },
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () {
              // Delete review
              Navigator.pop(context);
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Review deleted successfully'),
                  backgroundColor: Colors.red,
                ),
              );
            },
            style: TextButton.styleFrom(
              foregroundColor: Colors.red,
            ),
            child: const Text('Delete'),
          ),
        ],
      ),
    );
  }
}
