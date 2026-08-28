class Landmark {
  final String id;
  final String name;
  final double latitude;
  final double longitude;
  final String? currentOwnerClan;

  Landmark({
    required this.id,
    required this.name,
    required this.latitude,
    required this.longitude,
    this.currentOwnerClan,
  });
}

class LandmarkService {
  Future<List<Landmark>> getLandmarks() async {
    // Mock API call to get landmarks
    await Future.delayed(const Duration(milliseconds: 500));
    return [
      Landmark(
        id: 'lm-1',
        name: 'Phố Đi Bộ Hoàn Kiếm',
        latitude: 21.0285,
        longitude: 105.8542,
        currentOwnerClan: 'Hội Mỏ Neo',
      ),
      Landmark(
        id: 'lm-2',
        name: 'Nhà Thờ Lớn',
        latitude: 21.0289,
        longitude: 105.8491,
      ),
      Landmark(
        id: 'lm-3',
        name: 'Hồ Tây',
        latitude: 21.0538,
        longitude: 105.8239,
        currentOwnerClan: 'Dân Chơi GenZ',
      ),
    ];
  }

  Future<bool> checkIn(String landmarkId) async {
    // Mock API check-in
    await Future.delayed(const Duration(seconds: 1));
    return true;
  }
}
