class WingmanSuggestion {
  final String tone;
  final String text;

  WingmanSuggestion({
    required this.tone,
    required this.text,
  });

  factory WingmanSuggestion.fromJson(Map<String, dynamic> json) {
    return WingmanSuggestion(
      tone: json['tone'] ?? '',
      text: json['text'] ?? '',
    );
  }
}
