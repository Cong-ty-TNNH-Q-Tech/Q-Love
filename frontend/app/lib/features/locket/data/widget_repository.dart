// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:home_widget/home_widget.dart';
import 'package:flutter/material.dart';

class WidgetRepository {
  static const String appGroupId = 'group.com.qtech.qlove';
  static const String iOSWidgetName = 'LocketWidget';
  static const String androidWidgetName = 'LocketWidgetProvider';

  Future<void> initialize() async {
    await HomeWidget.setAppGroupId(appGroupId);
  }

  Future<void> updateWidgetData({
    required String senderName,
    required String? imageUrl,
    required int streak,
  }) async {
    try {
      await HomeWidget.saveWidgetData<String>('sender_name', senderName);
      if (imageUrl != null) {
        await HomeWidget.saveWidgetData<String>('image_url', imageUrl);
      }
      await HomeWidget.saveWidgetData<int>('streak', streak);

      // Save a flag whether it should be blurred or not
      await HomeWidget.saveWidgetData<bool>('is_blurred', streak < 10);

      await HomeWidget.updateWidget(
        iOSName: iOSWidgetName,
        androidName: androidWidgetName,
      );
    } catch (e) {
      debugPrint('Error updating widget data: $e');
      throw Exception('Error updating widget: $e');
    }
  }

  Future<void> updateWidgetImage(String key, String path) async {
    try {
      await HomeWidget.saveWidgetData<String>(key, path);
      await HomeWidget.updateWidget(
        iOSName: iOSWidgetName,
        androidName: androidWidgetName,
      );
    } catch (e) {
      debugPrint('Error updating widget image: $e');
    }
  }
}
