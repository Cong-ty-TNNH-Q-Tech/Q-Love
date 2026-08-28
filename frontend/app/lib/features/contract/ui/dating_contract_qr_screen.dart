import 'dart:async';
import 'package:flutter/material.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:qlove/features/contract/services/totp_service.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';

class DatingContractQrScreen extends StatefulWidget {
  final String secretKey;

  const DatingContractQrScreen({Key? key, required this.secretKey}) : super(key: key);

  @override
  _DatingContractQrScreenState createState() => _DatingContractQrScreenState();
}

class _DatingContractQrScreenState extends State<DatingContractQrScreen> {
  final TOTPService _totpService = TOTPService();
  late String _currentCode;
  late int _remainingSeconds;
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _updateTOTP();
    _startTimer();
  }

  void _updateTOTP() {
    setState(() {
      _currentCode = _totpService.generateTOTP(widget.secretKey);
      _remainingSeconds = _totpService.getRemainingSeconds();
    });
  }

  void _startTimer() {
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
      final remaining = _totpService.getRemainingSeconds();
      if (remaining == 30) {
        // A new 30s window just started
        _updateTOTP();
      } else {
        setState(() {
          _remainingSeconds = remaining;
        });
      }
    });
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0F0F1A),
      appBar: AppBar(
        title: Text(AppLocalizations.of(context)!.datingContractTitle, style: const TextStyle(fontWeight: FontWeight.bold)),
        backgroundColor: Colors.transparent,
        elevation: 0,
        centerTitle: true,
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(32.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(
                AppLocalizations.of(context)!.scanQrToSign,
                style: const TextStyle(color: Colors.white70, fontSize: 16),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 40),
              
              // QR Code Card
              Container(
                padding: const EdgeInsets.all(20),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(24),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.pinkAccent.withOpacity(0.3),
                      blurRadius: 30,
                      spreadRadius: 5,
                    )
                  ],
                ),
                child: QrImageView(
                  data: _currentCode,
                  version: QrVersions.auto,
                  size: 240.0,
                  foregroundColor: const Color(0xFF0F0F1A),
                ),
              ),
              
              const SizedBox(height: 40),
              
              // Code Display
              Text(
                "${_currentCode.substring(0, 3)} ${_currentCode.substring(3, 6)}",
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 48,
                  fontWeight: FontWeight.w900,
                  letterSpacing: 8,
                ),
              ),
              
              const SizedBox(height: 16),
              
              // Timer indicator
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.timer_outlined, color: Colors.pinkAccent, size: 20),
                  const SizedBox(width: 8),
                  Text(
                    AppLocalizations.of(context)!.qrRefreshIn(_remainingSeconds.toString()),
                    style: const TextStyle(color: Colors.pinkAccent, fontSize: 16, fontWeight: FontWeight.bold),
                  ),
                ],
              ),
              
              const SizedBox(height: 16),
              
              SizedBox(
                width: 200,
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(10),
                  child: LinearProgressIndicator(
                    value: _remainingSeconds / 30.0,
                    backgroundColor: Colors.white.withOpacity(0.1),
                    valueColor: const AlwaysStoppedAnimation<Color>(Colors.pinkAccent),
                    minHeight: 8,
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
