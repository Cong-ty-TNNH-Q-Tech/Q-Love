import 'package:flutter/material.dart';
import 'package:mapbox_maps_flutter/mapbox_maps_flutter.dart';
import 'package:qlove/features/clan/services/landmark_service.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';

class ClanMapScreen extends StatefulWidget {
  const ClanMapScreen({Key? key}) : super(key: key);

  @override
  _ClanMapScreenState createState() => _ClanMapScreenState();
}

class _ClanMapScreenState extends State<ClanMapScreen> {
  MapboxMap? mapboxMap;
  PointAnnotationManager? pointAnnotationManager;
  final LandmarkService _landmarkService = LandmarkService();
  List<Landmark> _landmarks = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    // In a real app, token should be fetched from .env or securely injected
    MapboxOptions.setAccessToken("pk.mock_access_token_for_qlove_dev");
    _fetchLandmarks();
  }

  Future<void> _fetchLandmarks() async {
    final data = await _landmarkService.getLandmarks();
    setState(() {
      _landmarks = data;
      _isLoading = false;
    });
    _addMarkers();
  }

  _onMapCreated(MapboxMap mapboxMap) async {
    this.mapboxMap = mapboxMap;
    // Set map style to Dark for Gen-Z aesthetic
    mapboxMap.loadStyleURI(MapboxStyles.DARK);
    
    pointAnnotationManager = await mapboxMap.annotations.createPointAnnotationManager();
    _addMarkers();
  }

  void _addMarkers() {
    if (pointAnnotationManager == null || _landmarks.isEmpty) return;

    pointAnnotationManager?.deleteAll();

    final options = _landmarks.map((l) => PointAnnotationOptions(
      geometry: Point(coordinates: Position(l.longitude, l.latitude)),
      textField: l.name,
      textColor: Colors.white.value,
      textHaloColor: Colors.pinkAccent.value,
      textHaloWidth: 1.5,
      textOffset: [0.0, -2.0],
      // Normally would use image for pin, using a default setup here
    )).toList();

    pointAnnotationManager?.createMulti(options);
  }

  void _handleCheckIn() async {
    // For demo purposes, assume we check in to the first landmark
    if (_landmarks.isEmpty) return;
    final success = await _landmarkService.checkIn(_landmarks.first.id);
    
    if (success && mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(AppLocalizations.of(context)!.landmarkCaptured)),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0F0F1A),
      appBar: AppBar(
        title: Text(AppLocalizations.of(context)!.clanMapTitle, style: const TextStyle(fontWeight: FontWeight.bold)),
        backgroundColor: Colors.transparent,
        elevation: 0,
        centerTitle: true,
      ),
      body: Stack(
        children: [
          MapWidget(
            key: const ValueKey("mapWidget"),
            onMapCreated: _onMapCreated,
            cameraOptions: CameraOptions(
              center: Point(coordinates: Position(105.8542, 21.0285)), // Hanoi center
              zoom: 13.0,
            ),
          ),
          
          if (_isLoading)
            const Center(child: CircularProgressIndicator(color: Colors.pinkAccent)),

          // Overlay Actions
          Positioned(
            bottom: 40,
            left: 24,
            right: 24,
            child: SizedBox(
              height: 60,
              child: ElevatedButton.icon(
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.pinkAccent,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(30),
                  ),
                  shadowColor: Colors.pinkAccent.withOpacity(0.5),
                  elevation: 10,
                ),
                icon: const Icon(Icons.flag, color: Colors.white, size: 28),
                label: Text(
                  AppLocalizations.of(context)!.checkInLandmark,
                  style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Colors.white),
                ),
                onPressed: _handleCheckIn,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
