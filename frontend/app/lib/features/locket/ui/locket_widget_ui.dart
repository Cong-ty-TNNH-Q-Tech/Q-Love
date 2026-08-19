// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:home_widget/home_widget.dart';
import 'package:qlove/features/locket/data/widget_repository.dart';

class LocketWidgetUI extends StatelessWidget {
  final String senderName;
  final String? imageUrl;
  final int streak;

  const LocketWidgetUI({
    super.key,
    required this.senderName,
    this.imageUrl,
    required this.streak,
  });

  @override
  Widget build(BuildContext context) {
    final isBlurred = streak < 10;

    return Container(
      width: 200,
      height: 200,
      decoration: BoxDecoration(
        color: const Color(0xFF1E1E2E),
        borderRadius: BorderRadius.circular(24),
      ),
      child: Stack(
        fit: StackFit.expand,
        children: [
          // Background Image
          if (imageUrl != null)
            ClipRRect(
              borderRadius: BorderRadius.circular(24),
              child: Image.network(
                imageUrl!,
                fit: BoxFit.cover,
              ),
            )
          else
            const Center(
              child: Icon(Icons.favorite_border, color: Colors.white24, size: 50),
            ),

          // Blur Overlay
          if (isBlurred)
            ClipRRect(
              borderRadius: BorderRadius.circular(24),
              child: BackdropFilter(
                filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10), // Gaussian Blur 90% equivalent
                child: Container(
                  color: Colors.black.withOpacity(0.1),
                ),
              ),
            ),

          // Bottom Bar (Q-Love Logo + Streak)
          Positioned(
            bottom: 12,
            left: 12,
            right: 12,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text(
                  'Q-Love',
                  style: TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                    fontSize: 14,
                  ),
                ),
                Row(
                  children: [
                    const Icon(Icons.local_fire_department, color: Color(0xFFFF4757), size: 16),
                    const SizedBox(width: 4),
                    Text(
                      '$streak',
                      style: const TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                        fontSize: 14,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  /// Render this UI into an image file and pass to HomeWidget
  static Future<void> updateWidgetRender(String senderName, String? imageUrl, int streak) async {
    final repo = WidgetRepository();
    await repo.updateWidgetData(senderName: senderName, imageUrl: imageUrl, streak: streak);

    try {
      final path = await HomeWidget.renderFlutterWidget(
        LocketWidgetUI(
          senderName: senderName,
          imageUrl: imageUrl,
          streak: streak,
        ),
        logicalSize: const Size(200, 200),
        key: 'widget_image',
      );
      if (path != null) {
        await repo.updateWidgetImage('widget_image', path as String);
      }
    } catch (e) {
      debugPrint('Failed to render widget: $e');
    }
  }
}
