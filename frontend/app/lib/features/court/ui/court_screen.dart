// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_card_swiper/flutter_card_swiper.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:qlove/features/court/bloc/court_bloc.dart';
import 'package:qlove/features/court/bloc/court_event.dart';
import 'package:qlove/features/court/bloc/court_state.dart';
import 'package:qlove/features/court/ui/widgets/court_case_card.dart';

class CourtScreen extends StatefulWidget {
  const CourtScreen({super.key});

  @override
  State<CourtScreen> createState() => _CourtScreenState();
}

class _CourtScreenState extends State<CourtScreen> {
  final CardSwiperController _swiperController = CardSwiperController();

  @override
  void initState() {
    super.initState();
    context.read<CourtBloc>().add(const FetchCasesRequested());
  }

  @override
  void dispose() {
    _swiperController.dispose();
    super.dispose();
  }

  void _onSwipe(int previousIndex, int? currentIndex, CardSwiperDirection direction, CourtLoaded state) {
    if (direction == CardSwiperDirection.right || direction == CardSwiperDirection.left) {
      final courtCase = state.cases[previousIndex];
      // Right (Red) -> Guilty
      // Left (Blue/Green) -> Not Guilty
      final voteType = direction == CardSwiperDirection.right ? 'guilty' : 'not_guilty';
      
      context.read<CourtBloc>().add(VoteActionRequested(
        caseId: courtCase.id,
        voteType: voteType,
      ));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0D0D14),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        title: Text(
          'Tòa Án Tình Yêu',
          style: GoogleFonts.outfit(
            color: const Color(0xFFFF4757),
            fontSize: 28,
            fontWeight: FontWeight.bold,
          ),
        ),
      ),
      body: BlocConsumer<CourtBloc, CourtState>(
        listener: (context, state) {
          if (state is CourtError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text(state.message)),
            );
          }
        },
        builder: (context, state) {
          if (state is CourtLoading || state is CourtInitial) {
            return const Center(child: CircularProgressIndicator(color: Color(0xFFFF4757)));
          }

          if (state is CourtLoaded) {
            if (state.cases.isEmpty) {
              return _buildEmptyState();
            }

            return Column(
              children: [
                Expanded(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
                    child: CardSwiper(
                      controller: _swiperController,
                      cardsCount: state.cases.length,
                      onSwipe: (previousIndex, currentIndex, direction) {
                        _onSwipe(previousIndex, currentIndex, direction, state);
                        return true;
                      },
                      cardBuilder: (context, index, percentThresholdX, percentThresholdY) {
                        return CourtCaseCard(courtCase: state.cases[index]);
                      },
                    ),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.only(bottom: 32.0, top: 16.0),
                  child: _buildVoteButtons(),
                ),
              ],
            );
          }

          return const SizedBox.shrink();
        },
      ),
    );
  }

  Widget _buildVoteButtons() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        _buildVoteButton(
          title: 'VÔ TỘI',
          color: const Color(0xFF2ED573), // Green for innocent
          onPressed: () => _swiperController.swipe(CardSwiperDirection.left),
        ),
        const SizedBox(width: 32),
        _buildVoteButton(
          title: 'CÓ TỘI',
          color: const Color(0xFFFF4757), // Red for guilty
          onPressed: () => _swiperController.swipe(CardSwiperDirection.right),
        ),
      ],
    );
  }

  Widget _buildVoteButton({required String title, required Color color, required VoidCallback onPressed}) {
    return ElevatedButton(
      onPressed: onPressed,
      style: ElevatedButton.styleFrom(
        backgroundColor: Colors.transparent,
        shadowColor: Colors.transparent,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(24),
          side: BorderSide(color: color, width: 2),
        ),
        padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 16),
      ),
      child: Text(
        title,
        style: GoogleFonts.outfit(
          color: color,
          fontSize: 18,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.balance, size: 100, color: Colors.white24),
          const SizedBox(height: 24),
          Text(
            'Hôm nay Q-Love yên bình.\nKhông có vụ kiện nào đang chờ.',
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
              context.read<CourtBloc>().add(const FetchCasesRequested(isRefresh: true));
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFFFF4757),
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(24)),
              padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 16),
            ),
            child: Text(
              'Tải lại',
              style: GoogleFonts.inter(fontSize: 16, fontWeight: FontWeight.bold),
            ),
          ),
        ],
      ),
    );
  }
}
