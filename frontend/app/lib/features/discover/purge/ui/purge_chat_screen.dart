// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter/material.dart';
import 'package:qlove/core/theme/app_theme.dart';
import 'dart:async';
import 'dart:math';

class PurgeChatScreen extends StatefulWidget {
  const PurgeChatScreen({super.key});

  @override
  State<PurgeChatScreen> createState() => _PurgeChatScreenState();
}

class _PurgeChatScreenState extends State<PurgeChatScreen> with SingleTickerProviderStateMixin {
  int _secondsLeft = 600; // 10 minutes
  Timer? _timer;
  late AnimationController _shakeController;
  final TextEditingController _msgController = TextEditingController();
  final List<String> _messages = ["Chào người lạ", "Bạn ở đâu thế?"];

  @override
  void initState() {
    super.initState();
    _shakeController = AnimationController(
      duration: const Duration(milliseconds: 50),
      vsync: this,
    );
    _startTimer();
  }

  void _startTimer() {
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
      if (_secondsLeft > 0) {
        setState(() {
          _secondsLeft--;
        });
        if (_secondsLeft <= 60) {
          _shakeController.repeat(reverse: true);
        }
      } else {
        _timer?.cancel();
        _shakeController.stop();
        // Show reveal button logic
      }
    });
  }

  @override
  void dispose() {
    _timer?.cancel();
    _shakeController.dispose();
    _msgController.dispose();
    super.dispose();
  }

  String get _timeString {
    int m = _secondsLeft ~/ 60;
    int s = _secondsLeft % 60;
    return '${m.toString().padLeft(2, '0')}:${s.toString().padLeft(2, '0')}';
  }

  @override
  Widget build(BuildContext context) {
    final bool isCritical = _secondsLeft <= 60;
    final bool isFinished = _secondsLeft == 0;

    return Scaffold(
      backgroundColor: const Color(0xFF080810),
      appBar: AppBar(
        backgroundColor: const Color(0xFF151020),
        title: const Text('Đối tượng ẩn danh'),
        actions: [
          AnimatedBuilder(
            animation: _shakeController,
            builder: (context, child) {
              double offset = 0;
              if (isCritical && !isFinished) {
                offset = (Random().nextDouble() - 0.5) * 4;
              }
              return Transform.translate(
                offset: Offset(offset, 0),
                child: Center(
                  child: Padding(
                    padding: const EdgeInsets.only(right: 16.0),
                    child: Text(
                      _timeString,
                      style: TextStyle(
                        color: isCritical ? AppTheme.neonRed : Colors.white,
                        fontWeight: FontWeight.bold,
                        fontSize: 18,
                      ),
                    ),
                  ),
                ),
              );
            },
          ),
        ],
      ),
      body: Column(
        children: [
          Expanded(
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: _messages.length,
              itemBuilder: (context, index) {
                bool isMe = index % 2 != 0;
                return Align(
                  alignment: isMe ? Alignment.centerRight : Alignment.centerLeft,
                  child: Container(
                    margin: const EdgeInsets.symmetric(vertical: 4),
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                    decoration: BoxDecoration(
                      color: isMe ? AppTheme.primaryColor : const Color(0xFF2C2C3E),
                      borderRadius: BorderRadius.circular(16),
                    ),
                    child: Text(
                      _messages[index],
                      style: const TextStyle(color: Colors.white),
                    ),
                  ),
                );
              },
            ),
          ),
          if (isFinished)
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: () {
                    // Go to reveal
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppTheme.neonBlue,
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  child: const Text(
                    'LỘ DIỆN',
                    style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                  ),
                ),
              ),
            ),
          // Input
          Container(
            padding: const EdgeInsets.all(12),
            color: const Color(0xFF151020),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _msgController,
                    decoration: InputDecoration(
                      hintText: 'Nhắn tin ẩn danh...',
                      hintStyle: const TextStyle(color: Colors.white54),
                      filled: true,
                      fillColor: const Color(0xFF2C2C3E),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(24),
                        borderSide: BorderSide.none,
                      ),
                      contentPadding: const EdgeInsets.symmetric(horizontal: 20),
                    ),
                    style: const TextStyle(color: Colors.white),
                  ),
                ),
                const SizedBox(width: 8),
                IconButton(
                  icon: const Icon(Icons.send, color: AppTheme.primaryColor),
                  onPressed: () {
                    if (_msgController.text.isNotEmpty) {
                      setState(() {
                        _messages.add(_msgController.text);
                        _msgController.clear();
                      });
                    }
                  },
                )
              ],
            ),
          )
        ],
      ),
    );
  }
}

