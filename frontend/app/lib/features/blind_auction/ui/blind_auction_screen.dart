import 'dart:async';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';

import 'widgets/auction_countdown.dart';
import 'widgets/bid_input_area.dart';
import 'widgets/top_hunters_list.dart';

class BlindAuctionScreen extends StatefulWidget {
  const BlindAuctionScreen({super.key});

  @override
  State<BlindAuctionScreen> createState() => _BlindAuctionScreenState();
}

class _BlindAuctionScreenState extends State<BlindAuctionScreen> {
  late Timer _timer;
  Duration _timeRemaining = const Duration(hours: 23, minutes: 59, seconds: 59);

  @override
  void initState() {
    super.initState();
    _startCountdown();
  }

  void _startCountdown() {
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
      if (_timeRemaining.inSeconds > 0) {
        setState(() {
          _timeRemaining -= const Duration(seconds: 1);
        });
      } else {
        _timer.cancel();
      }
    });
  }

  @override
  void dispose() {
    _timer.cancel();
    super.dispose();
  }

  void _handlePlaceBid() {
    // ScaffoldMessenger.of(context).showSnackBar(
    //   const SnackBar(content: Text('Bid placed successfully!')),
    // );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0D0D14), // Deep dark background
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.white),
          onPressed: () => Navigator.pop(context),
        ),
        title: Text(
          AppLocalizations.of(context)?.blindAuctionTitle ?? 'Đấu Giá Mù',
          style: GoogleFonts.outfit(
            color: Colors.white,
            fontWeight: FontWeight.bold,
          ),
        ),
        centerTitle: true,
      ),
      body: SafeArea(
        child: Column(
          children: [
            const SizedBox(height: 24),
            // TOP: Countdown Timer
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24.0),
              child: AuctionCountdown(timeRemaining: _timeRemaining),
            ),
            
            const SizedBox(height: 32),
            
            // MIDDLE: Bid Input
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24.0),
              child: BidInputArea(
                onBidChanged: (val) {
                  // handle bid change locally if needed
                },
                onPlaceBid: _handlePlaceBid,
              ),
            ),
            
            const SizedBox(height: 32),
            
            // BOTTOM: Leaderboard
            const Expanded(
              child: TopHuntersList(),
            ),
          ],
        ),
      ),
    );
  }
}
