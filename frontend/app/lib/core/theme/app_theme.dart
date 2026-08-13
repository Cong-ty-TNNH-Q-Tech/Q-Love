// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class AppTheme {
  // Premium Gen-Z Colors
  static const Color primaryColor = Color(0xFFFF2D55); // Vibrant Pink
  static const Color secondaryColor = Color(0xFF00F0FF); // Neon Cyan
  static const Color scaffoldBackground = Color(0xFF0A0A0A); // Deep dark
  static const Color surfaceColor = Color(0xFF1C1C1E); // Elevated dark

  static ThemeData get darkTheme {
    return ThemeData(
      brightness: Brightness.dark,
      primaryColor: primaryColor,
      scaffoldBackgroundColor: scaffoldBackground,
      colorScheme: const ColorScheme.dark(
        primary: primaryColor,
        secondary: secondaryColor,
        surface: surfaceColor,
      ),
      appBarTheme: const AppBarTheme(
        backgroundColor: Colors.transparent,
        elevation: 0,
        centerTitle: true,
      ),
      textTheme: GoogleFonts.interTextTheme(ThemeData.dark().textTheme).copyWith(
        displayLarge: GoogleFonts.inter(
          fontSize: 56,
          fontWeight: FontWeight.bold,
          letterSpacing: -1.5,
        ),
        displayMedium: GoogleFonts.inter(
          fontSize: 48,
          fontWeight: FontWeight.bold,
          letterSpacing: -0.5,
        ),
      ),
      useMaterial3: true,
    );
  }
}
