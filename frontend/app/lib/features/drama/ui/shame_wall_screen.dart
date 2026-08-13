// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/shame_wall_bloc.dart';

class ShameWallScreen extends StatelessWidget {
  const ShameWallScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) => ShameWallBloc()..add(LoadShameWall()),
      child: Scaffold(
        extendBodyBehindAppBar: true,
        appBar: AppBar(
          title: Text('Tường Thành Phong Sát', style: Theme.of(context).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.w800, color: Colors.white)),
          backgroundColor: Colors.transparent,
          elevation: 0,
        ),
        body: Container(
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              colors: [Color(0xFF0A0A0A), Color(0xFF1A0005)],
              begin: Alignment.topCenter,
              end: Alignment.bottomCenter,
            ),
          ),
          child: BlocBuilder<ShameWallBloc, ShameWallState>(
            builder: (context, state) {
              if (state is ShameWallLoading) {
                return const Center(child: CircularProgressIndicator(color: Color(0xFFFF2D55)));
              } else if (state is ShameWallLoaded) {
                return ListView.builder(
                  padding: const EdgeInsets.fromLTRB(16, 100, 16, 20),
                  itemCount: state.shames.length,
                  itemBuilder: (context, index) {
                    final shame = state.shames[index];
                    return _buildShameCard(context, shame);
                  },
                );
              } else if (state is ShameWallError) {
                return Center(child: Text(state.message, style: Theme.of(context).textTheme.bodyMedium?.copyWith(color: Colors.red)));
              }
              return const SizedBox.shrink();
            },
          ),
        ),
      ),
    );
  }

  Widget _buildShameCard(BuildContext context, ShameModel shame) {
    return Container(
      margin: const EdgeInsets.only(bottom: 20),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(20),
        child: BackdropFilter(
          filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
          child: Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.05),
              borderRadius: BorderRadius.circular(20),
              border: Border.all(color: Colors.white.withOpacity(0.1)),
              boxShadow: [
                BoxShadow(
                  color: const Color(0xFFFF2D55).withOpacity(0.1),
                  blurRadius: 20,
                  spreadRadius: -5,
                )
              ],
            ),
            child: Row(
              children: [
                // Avatar placeholder
                Container(
                  width: 60,
                  height: 60,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    color: Colors.white.withOpacity(0.1),
                    border: Border.all(color: const Color(0xFFFF2D55), width: 2),
                  ),
                  child: const Icon(Icons.person, color: Colors.white54, size: 30),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        shame.userName,
                        style: Theme.of(context).textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        shame.reason,
                        style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: Colors.white.withOpacity(0.7),
                        ),
                      ),
                      const SizedBox(height: 8),
                      Row(
                        children: [
                          Text('🍅 ', style: Theme.of(context).textTheme.bodyMedium),
                          Text(
                            '${shame.tomatoes} cà chua',
                            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                              fontWeight: FontWeight.w600,
                              color: const Color(0xFFFF2D55),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
                // Throw Tomato Button
                GestureDetector(
                  onTap: () {
                    HapticFeedback.heavyImpact();
                    context.read<ShameWallBloc>().add(ThrowTomato(shameId: shame.id));
                    // Show a quick snackbar or animation here
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(
                        content: Text('🍅 Đã ném cà chua! (-1 Xu)'),
                        duration: Duration(milliseconds: 1000),
                        backgroundColor: Color(0xFFFF2D55),
                      )
                    );
                  },
                  child: Container(
                    width: 50,
                    height: 50,
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      gradient: const LinearGradient(
                        colors: [Color(0xFFFF2D55), Color(0xFFFF7A00)],
                      ),
                      boxShadow: [
                        BoxShadow(
                          color: const Color(0xFFFF2D55).withOpacity(0.4),
                          blurRadius: 10,
                          spreadRadius: 2,
                        ),
                      ],
                    ),
                    child: Center(
                      child: Text('🍅', style: Theme.of(context).textTheme.headlineSmall),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
