// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:qlove/core/models/court_case_model.dart';

class CourtCaseCard extends StatelessWidget {
  final CourtCaseModel courtCase;

  const CourtCaseCard({super.key, required this.courtCase});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: const Color(0xFF1E1E2E), // Dark theme
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: const Color(0xFF2A2A3D), width: 2),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.5),
            blurRadius: 15,
            spreadRadius: 2,
          ),
        ],
      ),
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header: Defendant info
          Row(
            children: [
              const Icon(Icons.gavel, color: Color(0xFFFF4757), size: 32),
              const SizedBox(width: 12),
              Expanded(
                child: Text(
                  'Bị cáo: ${courtCase.defendantNameMasked}',
                  style: GoogleFonts.outfit(
                    color: Colors.white,
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),
          
          // Reason / Crime
          Text(
            'Lý do khởi kiện:',
            style: GoogleFonts.inter(
              color: Colors.white54,
              fontSize: 14,
              fontWeight: FontWeight.w500,
            ),
          ),
          const SizedBox(height: 8),
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: const Color(0xFF0D0D14),
              borderRadius: BorderRadius.circular(16),
              border: Border.all(color: const Color(0xFFFF4757).withOpacity(0.3)),
            ),
            child: Text(
              courtCase.reason,
              style: GoogleFonts.inter(
                color: Colors.white,
                fontSize: 18,
                height: 1.5,
              ),
            ),
          ),
          
          const Spacer(),
          
          // Footer: Vote count & Date
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(
                children: [
                  const Icon(Icons.people, color: Colors.white54, size: 16),
                  const SizedBox(width: 4),
                  Text(
                    '${courtCase.voteCount} Jury Votes',
                    style: GoogleFonts.inter(
                      color: Colors.white54,
                      fontSize: 14,
                    ),
                  ),
                ],
              ),
              Text(
                _formatDate(courtCase.createdAt),
                style: GoogleFonts.inter(
                  color: Colors.white38,
                  fontSize: 14,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  String _formatDate(DateTime date) {
    return '${date.day.toString().padLeft(2, '0')}/${date.month.toString().padLeft(2, '0')}/${date.year}';
  }
}
