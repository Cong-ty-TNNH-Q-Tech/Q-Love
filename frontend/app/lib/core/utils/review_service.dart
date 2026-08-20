import 'package:in_app_review/in_app_review.dart';
import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';

class ReviewService {
  static final ReviewService _instance = ReviewService._internal();
  factory ReviewService() => _instance;
  ReviewService._internal();

  final InAppReview _inAppReview = InAppReview.instance;
  static const String _keyMatchCount = 'match_count';
  static const String _keyHasReviewed = 'has_reviewed';

  Future<void> onUserMatched() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      
      final hasReviewed = prefs.getBool(_keyHasReviewed) ?? false;
      if (hasReviewed) return;

      int matchCount = (prefs.getInt(_keyMatchCount) ?? 0) + 1;
      await prefs.setInt(_keyMatchCount, matchCount);

      // Trigger review after 3 successful matches
      if (matchCount >= 3) {
        if (await _inAppReview.isAvailable()) {
          await _inAppReview.requestReview();
          await prefs.setBool(_keyHasReviewed, true);
        }
      }
    } catch (e) {
      debugPrint('Error triggering in-app review: $e');
    }
  }

  Future<void> openStoreListing() async {
    try {
      if (await _inAppReview.isAvailable()) {
        await _inAppReview.openStoreListing(appStoreId: '1234567890');
      }
    } catch (e) {
      debugPrint('Error opening store listing: $e');
    }
  }
}
