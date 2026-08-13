import 'dart:ui';
import 'package:flutter/material.dart';
import '../../data/datasources/chat_remote_data_source.dart';
import '../../domain/entities/wingman_suggestion.dart';

class AiWingmanBottomSheet extends StatefulWidget {
  final String matchId;
  final ChatRemoteDataSource dataSource;
  final Function(String) onSuggestionSelected;

  const AiWingmanBottomSheet({
    Key? key,
    required this.matchId,
    required this.dataSource,
    required this.onSuggestionSelected,
  }) : super(key: key);

  @override
  State<AiWingmanBottomSheet> createState() => _AiWingmanBottomSheetState();
}

class _AiWingmanBottomSheetState extends State<AiWingmanBottomSheet> {
  bool isLoading = true;
  List<WingmanSuggestion> suggestions = [];
  String? errorMessage;

  @override
  void initState() {
    super.initState();
    _fetchSuggestions();
  }

  Future<void> _fetchSuggestions() async {
    try {
      final data = await widget.dataSource.getWingmanSuggestions(widget.matchId);
      setState(() {
        suggestions = data;
        isLoading = false;
      });
    } catch (e) {
      setState(() {
        errorMessage = "Không thể kết nối đến Trợ lý Mỏ Hỗn 😢";
        isLoading = false;
      });
    }
  }

  IconData _getToneIcon(String tone) {
    if (tone.toLowerCase().contains("hài hước")) return Icons.emoji_emotions;
    if (tone.toLowerCase().contains("thả thính")) return Icons.favorite;
    return Icons.flash_on; // thẳng thắn
  }

  Color _getToneColor(String tone) {
    if (tone.toLowerCase().contains("hài hước")) return Colors.orangeAccent;
    if (tone.toLowerCase().contains("thả thính")) return Colors.pinkAccent;
    return Colors.cyanAccent;
  }

  @override
  Widget build(BuildContext context) {
    return ClipRRect(
      borderRadius: const BorderRadius.only(
        topLeft: Radius.circular(24),
        topRight: Radius.circular(24),
      ),
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: 15, sigmaY: 15),
        child: Container(
          padding: const EdgeInsets.all(20),
          decoration: BoxDecoration(
            color: const Color(0xFF1E1E24).withOpacity(0.6), // Glassmorphism background
            border: Border.all(
              color: Colors.white.withOpacity(0.2),
              width: 1.5,
            ),
            borderRadius: const BorderRadius.only(
              topLeft: Radius.circular(24),
              topRight: Radius.circular(24),
            ),
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              // Handle bar
              Center(
                child: Container(
                  width: 40,
                  height: 4,
                  margin: const EdgeInsets.only(bottom: 20),
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.3),
                    borderRadius: BorderRadius.circular(2),
                  ),
                ),
              ),
              
              // Title
              Row(
                children: [
                  const Icon(Icons.auto_awesome, color: Colors.pinkAccent),
                  const SizedBox(width: 8),
                  const Text(
                    "Trợ lý Mỏ Hỗn",
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const Spacer(),
                  IconButton(
                    icon: const Icon(Icons.refresh, color: Colors.white70),
                    onPressed: () {
                      setState(() {
                        isLoading = true;
                        errorMessage = null;
                      });
                      _fetchSuggestions();
                    },
                  ),
                ],
              ),
              const SizedBox(height: 16),
              
              // Content
              if (isLoading)
                const Padding(
                  padding: EdgeInsets.symmetric(vertical: 40),
                  child: Center(
                    child: CircularProgressIndicator(color: Colors.pinkAccent),
                  ),
                )
              else if (errorMessage != null)
                Padding(
                  padding: const EdgeInsets.symmetric(vertical: 20),
                  child: Center(
                    child: Text(
                      errorMessage!,
                      style: const TextStyle(color: Colors.redAccent, fontSize: 16),
                    ),
                  ),
                )
              else
                ...suggestions.map((suggestion) {
                  return Container(
                    margin: const EdgeInsets.only(bottom: 12),
                    decoration: BoxDecoration(
                      color: Colors.white.withOpacity(0.05),
                      borderRadius: BorderRadius.circular(16),
                      border: Border.all(color: Colors.white.withOpacity(0.1)),
                    ),
                    child: Material(
                      color: Colors.transparent,
                      child: InkWell(
                        borderRadius: BorderRadius.circular(16),
                        onTap: () => widget.onSuggestionSelected(suggestion.text),
                        child: Padding(
                          padding: const EdgeInsets.all(16),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                children: [
                                  Icon(
                                    _getToneIcon(suggestion.tone),
                                    color: _getToneColor(suggestion.tone),
                                    size: 16,
                                  ),
                                  const SizedBox(width: 8),
                                  Text(
                                    suggestion.tone,
                                    style: TextStyle(
                                      color: _getToneColor(suggestion.tone),
                                      fontWeight: FontWeight.bold,
                                      fontSize: 12,
                                    ),
                                  ),
                                ],
                              ),
                              const SizedBox(height: 8),
                              Text(
                                suggestion.text,
                                style: const TextStyle(
                                  color: Colors.white,
                                  fontSize: 15,
                                  height: 1.4,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                  );
                }).toList(),
                
              const SizedBox(height: 20),
            ],
          ),
        ),
      ),
    );
  }
}
