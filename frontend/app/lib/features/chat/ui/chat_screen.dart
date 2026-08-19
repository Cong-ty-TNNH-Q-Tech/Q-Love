// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/chat_bloc.dart';
import '../data/chat_repository.dart';
import 'widgets/chat_bubble.dart';

class ChatScreen extends StatelessWidget {
  final String matchId;
  final String targetId;

  final String currentUserId;

  const ChatScreen({super.key, required this.matchId, required this.targetId, required this.currentUserId});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) => ChatBloc(
        repository: ChatRepository(currentUserId: currentUserId),
      )..add(LoadChatHistory(matchId: matchId))
       ..add(ConnectWebSocket()),
      child: ChatView(matchId: matchId, targetId: targetId),
    );
  }
}

class ChatView extends StatefulWidget {
  final String matchId;
  final String targetId;

  const ChatView({super.key, required this.matchId, required this.targetId});

  @override
  State<ChatView> createState() => _ChatViewState();
}

class _ChatViewState extends State<ChatView> {
  final TextEditingController _textController = TextEditingController();
  final FocusNode _focusNode = FocusNode();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    
    return Scaffold(
      appBar: AppBar(
        title: Row(
          children: [
            const CircleAvatar(
              backgroundImage: NetworkImage('https://i.pravatar.cc/150?img=11'),
              radius: 18,
            ),
            const SizedBox(width: 12),
            Text('Partner', style: theme.textTheme.titleMedium),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.more_vert),
            onPressed: () {},
          ),
        ],
      ),
      body: Column(
        children: [
          Expanded(
            child: BlocBuilder<ChatBloc, ChatState>(
              builder: (context, state) {
                if (state is ChatLoading) {
                  return const Center(child: CircularProgressIndicator());
                } else if (state is ChatError) {
                  return Center(child: Text(state.message));
                } else if (state is ChatLoaded) {
                  return ListView.builder(
                    reverse: true, // Optimizes large lists by starting from bottom
                    itemCount: state.messages.length,
                    itemBuilder: (context, index) {
                      return ChatBubble(message: state.messages[index]);
                    },
                  );
                }
                return const SizedBox();
              },
            ),
          ),
          _buildInputBar(context),
        ],
      ),
    );
  }

  Widget _buildInputBar(BuildContext context) {
    final theme = Theme.of(context);
    
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border(top: BorderSide(color: Colors.white.withOpacity(0.05))),
      ),
      child: SafeArea(
        child: Row(
          children: [
            IconButton(
              icon: Icon(Icons.camera_alt, color: theme.colorScheme.primary),
              onPressed: () {
                // Send Locket mock
                context.read<ChatBloc>().add(
                  SendMessage(
                    matchId: widget.matchId,
                    targetId: widget.targetId,
                    content: 'https://i.pravatar.cc/300?img=12', // Mock locket URL
                    type: 'locket',
                  )
                );
              },
            ),
            IconButton(
              icon: Icon(Icons.smart_toy, color: theme.colorScheme.secondary),
              tooltip: 'AI Wingman',
              onPressed: () {
                _showWingmanPopup(context);
              },
            ),
            Expanded(
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.white.withOpacity(0.05),
                  borderRadius: BorderRadius.circular(24),
                ),
                child: TextField(
                  controller: _textController,
                  focusNode: _focusNode,
                  decoration: const InputDecoration(
                    hintText: 'Nhắn tin...',
                    hintStyle: TextStyle(color: Colors.white38),
                    border: InputBorder.none,
                    contentPadding: EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                  ),
                  onSubmitted: (_) => _sendMessage(context),
                ),
              ),
            ),
            const SizedBox(width: 8),
            CircleAvatar(
              backgroundColor: theme.colorScheme.primary,
              child: IconButton(
                icon: const Icon(Icons.send, color: Colors.white, size: 20),
                onPressed: () => _sendMessage(context),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _sendMessage(BuildContext context) {
    final text = _textController.text.trim();
    if (text.isNotEmpty) {
      context.read<ChatBloc>().add(
        SendMessage(
          matchId: widget.matchId,
          targetId: widget.targetId,
          content: text,
        )
      );
      _textController.clear();
      _focusNode.requestFocus();
    }
  }

  void _showWingmanPopup(BuildContext context) {
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      builder: (ctx) => Container(
        padding: const EdgeInsets.all(24),
        decoration: BoxDecoration(
          color: Theme.of(context).colorScheme.surface,
          borderRadius: const BorderRadius.vertical(top: Radius.circular(30)),
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Row(
              children: [
                Icon(Icons.smart_toy, color: Theme.of(context).colorScheme.secondary),
                const SizedBox(width: 12),
                Text('AI Wingman Gợi Ý', style: Theme.of(context).textTheme.titleLarge),
              ],
            ),
            const SizedBox(height: 20),
            _buildWingmanSuggestion(context, 'Hài hước', 'Em ăn cơm chưa hay để anh ăn hộ? 🤣'),
            _buildWingmanSuggestion(context, 'Thả thính', 'Bầu trời hôm nay nhiều sao, nhưng sao anh chỉ thấy em? ✨'),
            _buildWingmanSuggestion(context, 'Thẳng thắn', 'Cuối tuần này đi cà phê không em? ☕'),
            const SizedBox(height: 20),
          ],
        ),
      ),
    );
  }

  Widget _buildWingmanSuggestion(BuildContext context, String title, String content) {
    return InkWell(
      onTap: () {
        _textController.text = content;
        Navigator.pop(context);
      },
      child: Container(
        margin: const EdgeInsets.only(bottom: 12),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.05),
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: Colors.white.withOpacity(0.1)),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(title, style: TextStyle(color: Theme.of(context).colorScheme.secondary, fontWeight: FontWeight.bold, fontSize: 12)),
            const SizedBox(height: 4),
            Text(content, style: const TextStyle(color: Colors.white)),
          ],
        ),
      ),
    );
  }
}
