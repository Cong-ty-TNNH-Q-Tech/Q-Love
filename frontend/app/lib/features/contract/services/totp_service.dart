import 'dart:math';
import 'package:crypto/crypto.dart';
import 'package:base32/base32.dart';

class TOTPService {
  /// Generate a TOTP code given a base32 encoded secret and a timestamp.
  /// Standard TOTP uses 30 seconds time step.
  String generateTOTP(String secret, {int? time, int period = 30}) {
    time ??= DateTime.now().millisecondsSinceEpoch;
    
    // Calculate the number of time steps
    int timeStep = (time / 1000).floor() ~/ period;
    
    // Convert time step to a 8-byte array (big-endian)
    List<int> msg = List<int>.filled(8, 0);
    for (int i = 7; i >= 0; i--) {
      msg[i] = timeStep & 0xff;
      timeStep >>= 8;
    }
    
    // Decode the base32 secret
    List<int> key = base32.decode(secret);
    
    // Generate HMAC-SHA1
    var hmac = Hmac(sha1, key);
    var digest = hmac.convert(msg);
    var hash = digest.bytes;
    
    // Dynamic truncation (RFC 4226)
    int offset = hash[hash.length - 1] & 0xf;
    
    int binary = ((hash[offset] & 0x7f) << 24) |
                 ((hash[offset + 1] & 0xff) << 16) |
                 ((hash[offset + 2] & 0xff) << 8) |
                 (hash[offset + 3] & 0xff);
                 
    int otp = binary % pow(10, 6).toInt();
    
    return otp.toString().padLeft(6, '0');
  }

  /// Calculates how many seconds are remaining in the current 30s period
  int getRemainingSeconds({int period = 30}) {
    final now = DateTime.now().millisecondsSinceEpoch;
    final seconds = (now / 1000).floor();
    return period - (seconds % period);
  }
}
