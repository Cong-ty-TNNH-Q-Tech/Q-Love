// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:qlove/core/theme/app_theme.dart';
import 'package:qlove/features/auth/bloc/auth_bloc.dart';
import 'package:qlove/features/auth/bloc/auth_event.dart';
import 'package:qlove/features/auth/bloc/auth_state.dart';
import 'package:qlove/features/auth/ui/otp_verification_screen.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import 'widgets/glass_text_field.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final TextEditingController _phoneController = TextEditingController();

  void _onLoginPressed() {
    final phone = _phoneController.text.trim();
    if (phone.isNotEmpty) {
      context.read<AuthBloc>().add(SendOtpRequested(phone));
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(AppLocalizations.of(context)!.enterPhonePrompt)),
      );
    }
  }

  @override
  void dispose() {
    _phoneController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return BlocListener<AuthBloc, AuthState>(
      listener: (context, state) {
        if (state is AuthOtpSent) {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (_) => OtpVerificationScreen(phone: state.phone),
            ),
          );
        } else if (state is AuthError) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text(state.message)),
          );
        }
      },
      child: Scaffold(
        body: Container(
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              colors: [
                AppTheme.scaffoldBackground,
                Color(0xFF2A0010), // Deep pinkish dark
              ],
              begin: Alignment.topCenter,
              end: Alignment.bottomCenter,
            ),
          ),
          child: SafeArea(
            child: SingleChildScrollView(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 24.0, vertical: 48.0),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    const Icon(
                      Icons.favorite,
                      color: AppTheme.primaryColor,
                      size: 80,
                    ),
                    const SizedBox(height: 24),
                    Text(
                      AppLocalizations.of(context)!.appTitle,
                      textAlign: TextAlign.center,
                      style: Theme.of(context).textTheme.displayMedium?.copyWith(
                            color: AppTheme.primaryColor,
                          ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      AppLocalizations.of(context)!.loveJourneySubtitle,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        color: Colors.white.withOpacity(0.7),
                        fontSize: 16,
                      ),
                    ),
                    const SizedBox(height: 48),
                    GlassTextField(
                      controller: _phoneController,
                      hintText: AppLocalizations.of(context)!.enterPhoneHint,
                      keyboardType: TextInputType.phone,
                      prefixIcon: Icons.phone,
                    ),
                    const SizedBox(height: 24),
                    BlocBuilder<AuthBloc, AuthState>(
                      builder: (context, state) {
                        final isLoading = state is AuthLoading;
                        return ElevatedButton(
                          onPressed: isLoading ? null : _onLoginPressed,
                          style: ElevatedButton.styleFrom(
                            backgroundColor: AppTheme.primaryColor,
                            foregroundColor: Colors.white,
                            padding: const EdgeInsets.symmetric(vertical: 16),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(16),
                            ),
                            elevation: 8,
                            shadowColor: AppTheme.primaryColor.withOpacity(0.5),
                          ),
                          child: isLoading
                              ? const SizedBox(
                                  height: 24,
                                  width: 24,
                                  child: CircularProgressIndicator(
                                    color: Colors.white,
                                    strokeWidth: 2,
                                  ),
                                )
                              : Text(
                                  AppLocalizations.of(context)!.continueBtn,
                                  style: const TextStyle(
                                    fontSize: 18,
                                    fontWeight: FontWeight.bold,
                                    letterSpacing: 1.2,
                                  ),
                                ),
                        );
                      },
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
