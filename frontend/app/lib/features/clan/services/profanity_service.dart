class ProfanityService {
  // A basic hardcoded list of Vietnamese profanities
  // In a real production app, this should be fetched from backend or a comprehensive library.
  final List<String> _blacklist = [
    'dm', 'đm', 'vkl', 'vl', 'địt', 'loz', 'lồn', 'cac', 'cặc', 'chó', 'đĩ', 'phò'
  ];

  bool containsProfanity(String text) {
    if (text.isEmpty) return false;
    
    final lowerText = text.toLowerCase();
    
    // Check if any word in the text matches the blacklist exactly
    final words = lowerText.split(RegExp(r'\s+'));
    for (final word in words) {
      if (_blacklist.contains(word)) {
        return true;
      }
    }
    
    // Also check for substrings for exact severe bad words (simplified)
    for (final badword in _blacklist) {
      if (lowerText.contains(' $badword ') || lowerText.startsWith('$badword ') || lowerText.endsWith(' $badword') || lowerText == badword) {
        return true;
      }
    }
    
    return false;
  }
}
