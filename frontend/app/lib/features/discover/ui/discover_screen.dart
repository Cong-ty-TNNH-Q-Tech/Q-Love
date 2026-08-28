// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_card_swiper/flutter_card_swiper.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import 'package:qlove/features/discover/bloc/discover_bloc.dart';
import 'package:qlove/features/discover/bloc/discover_event.dart';
import 'package:qlove/features/discover/bloc/discover_state.dart';
import 'package:qlove/features/discover/ui/widgets/profile_swipe_card.dart';
import 'package:qlove/features/discover/ui/widgets/swipe_buttons.dart';

class DiscoverScreen extends StatefulWidget {
  const DiscoverScreen({super.key});

  @override
  State<DiscoverScreen> createState() => _DiscoverScreenState();
}

class _DiscoverScreenState extends State<DiscoverScreen> {
  final CardSwiperController _swiperController = CardSwiperController();
  String _currentFilter = 'default';

  @override
  void initState() {
    super.initState();
    context.read<DiscoverBloc>().add(FetchFeedRequested(filter: _currentFilter));
  }

  @override
  void dispose() {
    _swiperController.dispose();
    super.dispose();
  }

  void _onSwipe(int previousIndex, int? currentIndex, CardSwiperDirection direction, DiscoverLoaded state) {
    if (direction == CardSwiperDirection.right || direction == CardSwiperDirection.left) {
      final user = state.profiles[previousIndex];
      final action = direction == CardSwiperDirection.right ? 'like' : 'pass';
      
      context.read<DiscoverBloc>().add(SwipeActionRequested(
        targetId: user.id,
        action: action,
      ));
    }
  }

  void _showFilterSheet() {
    showModalBottomSheet(
      context: context,
      backgroundColor: const Color(0xFF1E1E2E),
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (ctx) {
        return Container(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                AppLocalizations.of(context)!.filter,
                style: GoogleFonts.inter(
                  color: Colors.white,
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 24),
              ListTile(
                title: Text(AppLocalizations.of(context)!.defaultFilter, style: GoogleFonts.inter(color: Colors.white)),
                trailing: _currentFilter == 'default' ? const Icon(Icons.check, color: Color(0xFFFF3D6B)) : null,
                onTap: () {
                  setState(() => _currentFilter = 'default');
                  context.read<DiscoverBloc>().add(const FetchFeedRequested(filter: 'default', isRefresh: true));
                  Navigator.pop(ctx);
                },
              ),
              ListTile(
                title: Text(AppLocalizations.of(context)!.spiritualFilter, style: GoogleFonts.inter(color: Colors.white)),
                trailing: _currentFilter == 'spiritual' ? const Icon(Icons.check, color: Color(0xFFFF3D6B)) : null,
                onTap: () {
                  setState(() => _currentFilter = 'spiritual');
                  context.read<DiscoverBloc>().add(const FetchFeedRequested(filter: 'spiritual', isRefresh: true));
                  Navigator.pop(ctx);
                },
              ),
            ],
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0D0D14),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        title: Text(
          'Q-Love',
          style: GoogleFonts.outfit(
            color: const Color(0xFFFF3D6B), // Primary Pink
            fontSize: 28,
            fontWeight: FontWeight.bold,
          ),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.tune, color: Colors.white),
            onPressed: _showFilterSheet,
          ),
        ],
      ),
      body: BlocConsumer<DiscoverBloc, DiscoverState>(
        listener: (context, state) {
          if (state is DiscoverError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text(state.message)),
            );
          } else if (state is DiscoverMatch) {
            // Show match dialog
            showDialog(
              context: context,
              barrierDismissible: false,
              builder: (ctx) => AlertDialog(
                backgroundColor: const Color(0xFF1E1E2E),
                title: Text(AppLocalizations.of(context)!.itsAMatch, style: GoogleFonts.inter(color: Colors.white)),
                content: Text(
                  AppLocalizations.of(context)!.matchDescription(state.matchedUser.name ?? 'người ấy'),
                  style: GoogleFonts.inter(color: Colors.white70),
                ),
                actions: [
                  TextButton(
                    onPressed: () {
                      Navigator.pop(ctx);
                      context.read<DiscoverBloc>().add(ResumeDiscoverRequested());
                    },
                    child: Text(AppLocalizations.of(context)!.keepSwiping, style: GoogleFonts.inter(color: const Color(0xFFFF3D6B))),
                  ),
                ],
              ),
            );
          }
        },
        builder: (context, state) {
          if (state is DiscoverLoading || state is DiscoverInitial) {
            return const Center(child: CircularProgressIndicator(color: Color(0xFFFF3D6B)));
          }

          if (state is DiscoverLoaded) {
            if (state.profiles.isEmpty) {
              return _buildEmptyState();
            }

            return Column(
              children: [
                Expanded(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
                    child: CardSwiper(
                      controller: _swiperController,
                      cardsCount: state.profiles.length,
                      allowedSwipeDirection: const AllowedSwipeDirection.only(right: true, left: true),
                      onSwipe: (previousIndex, currentIndex, direction) {
                        _onSwipe(previousIndex, currentIndex, direction, state);
                        return true;
                      },
                      cardBuilder: (context, index, percentThresholdX, percentThresholdY) {
                        return ProfileSwipeCard(profile: state.profiles[index]);
                      },
                    ),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.only(bottom: 32.0, top: 16.0),
                  child: SwipeButtons(
                    onPass: () => _swiperController.swipe(CardSwiperDirection.left),
                    onLike: () => _swiperController.swipe(CardSwiperDirection.right),
                  ),
                ),
              ],
            );
          }

          return const SizedBox.shrink();
        },
      ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.person_search, size: 100, color: Colors.white24),
          const SizedBox(height: 24),
          Text(
            AppLocalizations.of(context)!.outOfSwipes(AppLocalizations.of(context)!.expandRadius),
            textAlign: TextAlign.center,
            style: GoogleFonts.inter(
              color: Colors.white,
              fontSize: 18,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 24),
          ElevatedButton(
            onPressed: () {
              context.read<DiscoverBloc>().add(const FetchFeedRequested(isRefresh: true));
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFFFF3D6B),
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(24)),
              padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 16),
            ),
            child: Text(
              AppLocalizations.of(context)!.expandRadius,
              style: GoogleFonts.inter(fontSize: 16, fontWeight: FontWeight.bold),
            ),
          ),
        ],
      ),
    );
  }
}
