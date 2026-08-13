// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/wingman_bloc.dart';

class MatchmakerPopup extends StatelessWidget {
  final String targetId;

  const MatchmakerPopup({super.key, required this.targetId});

  static void show(BuildContext context, String targetId) {
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      isScrollControlled: true,
      builder: (context) => BlocProvider(
        create: (context) => WingmanBloc(),
        child: MatchmakerPopup(targetId: targetId),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFF1E1E2C).withOpacity(0.9),
        borderRadius: BorderRadius.circular(30),
        border: Border.all(color: Colors.white.withOpacity(0.1)),
        boxShadow: [
          BoxShadow(
            color: const Color(0xFFFF2D55).withOpacity(0.2),
            blurRadius: 30,
            spreadRadius: 5,
          )
        ],
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(30),
        child: BackdropFilter(
          filter: ImageFilter.blur(sigmaX: 20, sigmaY: 20),
          child: Padding(
            padding: const EdgeInsets.all(24.0),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const SizedBox(height: 8),
                Text(
                  'Mai Mối Cò Mối',
                  style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  'Bạn muốn ghép đôi ai với người này?',
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: Colors.white70,
                  ),
                ),
                const SizedBox(height: 24),
                // Placeholder for Friend Selection
                Container(
                  height: 120,
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.05),
                    borderRadius: BorderRadius.circular(20),
                    border: Border.all(color: Colors.white.withOpacity(0.1)),
                  ),
                  child: Center(
                    child: Text(
                      'Danh sách bạn bè (Mock)',
                      style: Theme.of(context).textTheme.bodyMedium?.copyWith(color: Colors.white54),
                    ),
                  ),
                ),
                const SizedBox(height: 24),
                BlocConsumer<WingmanBloc, WingmanState>(
                  listener: (context, state) {
                    if (state is MatchmakerSuccess) {
                      Navigator.of(context).pop();
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('Ghép đôi thành công! Chờ họ đồng ý nhé.'),
                          backgroundColor: Color(0xFFFF2D55),
                        ),
                      );
                    } else if (state is WingmanError) {
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(content: Text(state.message), backgroundColor: Colors.red),
                      );
                    }
                  },
                  builder: (context, state) {
                    if (state is WingmanLoading) {
                      return const CircularProgressIndicator(color: Color(0xFFFF2D55));
                    }
                    return ElevatedButton(
                      onPressed: () {
                        context.read<WingmanBloc>().add(
                          MatchFriendEvent(targetId: targetId, friendId: 'mock_friend_123'),
                        );
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFFFF2D55),
                        foregroundColor: Colors.white,
                        minimumSize: const Size(double.infinity, 56),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(16),
                        ),
                        elevation: 8,
                        shadowColor: const Color(0xFFFF2D55).withOpacity(0.5),
                      ),
                      child: const Text(
                        'Ghép Đôi Ngay',
                        style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                      ),
                    );
                  },
                ),
                const SizedBox(height: 16),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
