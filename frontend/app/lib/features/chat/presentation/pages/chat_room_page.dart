import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import '../widgets/ai_wingman_bottom_sheet.dart';
import '../../data/datasources/chat_remote_data_source.dart';

class ChatRoomPage extends StatefulWidget {
  final String matchId;

  const ChatRoomPage({Key? key, required this.matchId}) : super(key: key);

  @override
  State<ChatRoomPage> createState() => _ChatRoomPageState();
}

class _ChatRoomPageState extends State<ChatRoomPage> with SingleTickerProviderStateMixin {
  final TextEditingController _messageController = TextEditingController();
  late AnimationController _glowController;
  late Animation<double> _glowAnimation;

  // Mock property to simulate dead air
  bool isChatFrozen = true; // Set to true to show Wingman icon for demo

  @override
  void initState() {
    super.initState();
    _glowController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 1),
    )..repeat(reverse: true);

    _glowAnimation = Tween<double>(begin: 4.0, end: 12.0).animate(
      CurvedAnimation(parent: _glowController, curve: Curves.easeInOut),
    );
  }

  @override
  void dispose() {
    _glowController.dispose();
    _messageController.dispose();
    super.dispose();
  }

  void _showWingmanBottomSheet() {
    final dataSource = ChatRemoteDataSourceImpl(client: http.Client());
    
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      isScrollControlled: true,
      builder: (context) {
        return AiWingmanBottomSheet(
          matchId: widget.matchId,
          dataSource: dataSource,
          onSuggestionSelected: (text) {
            _messageController.text = text;
            Navigator.pop(context);
          },
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0D0D12), // Dark first aesthetic
      appBar: AppBar(
        backgroundColor: const Color(0xFF15151A),
        title: const Text('Match #123', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
        elevation: 0,
        iconTheme: const IconThemeData(color: Colors.white),
      ),
      body: Column(
        children: [
          Expanded(
            child: ListView(
              padding: const EdgeInsets.all(16),
              children: const [
                Align(
                  alignment: Alignment.centerLeft,
                  child: Container(
                    margin: EdgeInsets.only(bottom: 12),
                    padding: EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: Color(0xFF1E1E24),
                      borderRadius: BorderRadius.only(
                        topLeft: Radius.circular(16),
                        topRight: Radius.circular(16),
                        bottomRight: Radius.circular(16),
                      ),
                    ),
                    child: Text("Hello!", style: TextStyle(color: Colors.white)),
                  ),
                ),
                Center(
                  child: Padding(
                    padding: EdgeInsets.symmetric(vertical: 20),
                    child: Text(
                      "2 hours ago",
                      style: TextStyle(color: Colors.grey, fontSize: 12),
                    ),
                  ),
                )
              ],
            ),
          ),
          
          // AI Wingman Icon (Mỏ Hỗn) with Glow Effect
          if (isChatFrozen)
            Align(
              alignment: Alignment.bottomRight,
              child: Padding(
                padding: const EdgeInsets.only(right: 16.0, bottom: 8.0),
                child: GestureDetector(
                  onTap: _showWingmanBottomSheet,
                  child: AnimatedBuilder(
                    animation: _glowAnimation,
                    builder: (context, child) {
                      return Container(
                        decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          boxShadow: [
                            BoxShadow(
                              color: Colors.pinkAccent.withOpacity(0.6),
                              blurRadius: _glowAnimation.value,
                              spreadRadius: _glowAnimation.value / 2,
                            ),
                          ],
                        ),
                        child: child,
                      );
                    },
                    child: const CircleAvatar(
                      radius: 24,
                      backgroundColor: Colors.pinkAccent,
                      child: Icon(Icons.auto_awesome, color: Colors.white),
                    ),
                  ),
                ),
              ),
            ),
          
          // Chat Input
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            decoration: const BoxDecoration(
              color: Color(0xFF15151A),
              border: Border(top: BorderSide(color: Color(0xFF2A2A35))),
            ),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _messageController,
                    style: const TextStyle(color: Colors.white),
                    decoration: InputDecoration(
                      hintText: "Type a message...",
                      hintStyle: const TextStyle(color: Colors.grey),
                      filled: true,
                      fillColor: const Color(0xFF1E1E24),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(24),
                        borderSide: BorderSide.none,
                      ),
                      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Container(
                  decoration: const BoxDecoration(
                    shape: BoxShape.circle,
                    gradient: LinearGradient(
                      colors: [Colors.pinkAccent, Colors.deepPurpleAccent],
                    ),
                  ),
                  child: IconButton(
                    icon: const Icon(Icons.send, color: Colors.white, size: 20),
                    onPressed: () {
                      _messageController.clear();
                      setState(() {
                        isChatFrozen = false;
                      });
                    },
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
