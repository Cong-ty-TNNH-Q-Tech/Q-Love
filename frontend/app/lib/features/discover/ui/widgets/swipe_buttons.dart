// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter/material.dart';

class SwipeButtons extends StatelessWidget {
  final VoidCallback onPass;
  final VoidCallback onLike;

  const SwipeButtons({
    super.key,
    required this.onPass,
    required this.onLike,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        _buildButton(
          icon: Icons.close,
          color: const Color(0xFFFF4757), // Pass Red
          size: 64,
          onPressed: onPass,
        ),
        const SizedBox(width: 32),
        _buildButton(
          icon: Icons.favorite,
          color: const Color(0xFFFF3D6B), // Like Pink
          size: 72, // Slightly larger
          onPressed: onLike,
        ),
      ],
    );
  }

  Widget _buildButton({
    required IconData icon,
    required Color color,
    required double size,
    required VoidCallback onPressed,
  }) {
    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        color: Colors.transparent,
        border: Border.all(color: color, width: 2),
        boxShadow: [
          BoxShadow(
            color: color.withOpacity(0.2),
            blurRadius: 20,
            spreadRadius: 2,
          ),
        ],
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onPressed,
          customBorder: const CircleBorder(),
          splashColor: color.withOpacity(0.3),
          child: Icon(
            icon,
            color: color,
            size: size * 0.5,
          ),
        ),
      ),
    );
  }
}
