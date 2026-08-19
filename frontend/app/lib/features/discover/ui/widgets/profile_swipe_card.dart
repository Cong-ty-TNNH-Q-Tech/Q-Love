// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter/material.dart';
import 'package:qlove/core/models/user_model.dart';
import 'package:google_fonts/google_fonts.dart';

class ProfileSwipeCard extends StatelessWidget {
  final UserModel profile;

  const ProfileSwipeCard({super.key, required this.profile});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(24),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.5),
            blurRadius: 15,
            spreadRadius: 2,
          ),
        ],
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(24),
        child: Stack(
          fit: StackFit.expand,
          children: [
            // Profile Image (Mock if null)
            profile.avatarUrl != null && profile.avatarUrl!.isNotEmpty
                ? Image.network(
                    profile.avatarUrl!,
                    fit: BoxFit.cover,
                    errorBuilder: (context, error, stackTrace) => _buildPlaceholder(),
                  )
                : _buildPlaceholder(),
                
            // Gradient Overlay
            Container(
              decoration: const BoxDecoration(
                gradient: LinearGradient(
                  colors: [Colors.transparent, Color(0xFF0D0D14)],
                  begin: Alignment.topCenter,
                  end: Alignment.bottomCenter,
                  stops: [0.4, 1.0],
                ),
              ),
            ),
            
            // Profile Info
            Positioned(
              bottom: 24,
              left: 20,
              right: 20,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Text(
                        '${profile.name ?? 'Ẩn danh'}, ${profile.age ?? '?'}',
                        style: GoogleFonts.inter(
                          color: Colors.white,
                          fontSize: 28,
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                      const SizedBox(width: 8),
                      if (profile.spiritualCompatibility != null)
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                          decoration: BoxDecoration(
                            color: const Color(0xFF8B5CF6).withOpacity(0.3), // spiritual purple
                            borderRadius: BorderRadius.circular(12),
                            border: Border.all(color: const Color(0xFF8B5CF6)),
                          ),
                          child: Text(
                            '🔮 ${profile.spiritualCompatibility}%',
                            style: GoogleFonts.inter(
                              color: Colors.white,
                              fontSize: 12,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  if (profile.distanceKm != null)
                    Row(
                      children: [
                        const Icon(Icons.location_on, color: Color(0xFFFFA5B4), size: 16),
                        const SizedBox(width: 4),
                        Text(
                          'Cách đây ${profile.distanceKm}km',
                          style: GoogleFonts.inter(
                            color: const Color(0xFFFFA5B4),
                            fontSize: 14,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ],
                    ),
                  const SizedBox(height: 12),
                  if (profile.quote != null && profile.quote!.isNotEmpty)
                    Text(
                      '"${profile.quote}"',
                      style: GoogleFonts.inter(
                        color: Colors.white70,
                        fontSize: 16,
                        fontStyle: FontStyle.italic,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildPlaceholder() {
    return Container(
      color: const Color(0xFF2A2A3D),
      child: const Center(
        child: Icon(Icons.person, size: 100, color: Colors.white24),
      ),
    );
  }
}
