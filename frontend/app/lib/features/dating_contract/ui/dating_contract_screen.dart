// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter/material.dart';
import 'widgets/propose_date_modal.dart';
import 'widgets/qr_display_widget.dart';
import 'qr_scanner_screen.dart';

class DatingContractScreen extends StatelessWidget {
  const DatingContractScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 2,
      child: Scaffold(
        appBar: AppBar(
          title: const Text('Khế Ước Hẹn Hò'),
          bottom: const TabBar(
            tabs: [
              Tab(text: 'Đang chờ (Pending)'),
              Tab(text: 'Đang hiệu lực (Active)'),
            ],
          ),
        ),
        body: TabBarView(
          children: [
            _buildPendingTab(context),
            _buildActiveTab(context),
          ],
        ),
        floatingActionButton: FloatingActionButton.extended(
          onPressed: () async {
            final xu = await showModalBottomSheet<String>(
              context: context,
              isScrollControlled: true,
              backgroundColor: Theme.of(context).colorScheme.surface,
              shape: const RoundedRectangleBorder(
                borderRadius: BorderRadius.vertical(top: Radius.circular(30)),
              ),
              builder: (ctx) => const ProposeDateModal(matchId: 'mock'),
            );
            if (xu != null && xu.isNotEmpty) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text('Đã gửi yêu cầu hẹn hò cọc $xu Xu!')),
              );
            }
          },
          label: const Text('Tạo Khế Ước'),
          icon: const Icon(Icons.handshake),
          backgroundColor: Theme.of(context).colorScheme.primary,
        ),
      ),
    );
  }

  Widget _buildPendingTab(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        Card(
          child: ListTile(
            leading: const CircleAvatar(child: Icon(Icons.hourglass_empty)),
            title: const Text('Hẹn hò với Mai Anh'),
            subtitle: const Text('Cọc: 200 Xu - Đợi phản hồi'),
            trailing: TextButton(
              onPressed: () {},
              child: Text('Hủy', style: TextStyle(color: Theme.of(context).colorScheme.error)),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildActiveTab(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        Card(
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              children: [
                const ListTile(
                  leading: CircleAvatar(backgroundImage: NetworkImage('https://i.pravatar.cc/150?img=5')),
                  title: Text('Hẹn hò với Tuấn'),
                  subtitle: Text('Đã cọc: 500 Xu\nThời gian: 19:00 Hôm nay'),
                ),
                const SizedBox(height: 16),
                const QrDisplayWidget(data: 'TOTP_MOCK_DATA_XYZ_123', size: 150),
                const SizedBox(height: 16),
                ElevatedButton.icon(
                  icon: const Icon(Icons.qr_code_scanner),
                  label: const Text('Quét mã đối phương để xác nhận'),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Theme.of(context).colorScheme.secondary,
                    foregroundColor: Colors.black,
                  ),
                  onPressed: () async {
                    final code = await Navigator.push(
                      context,
                      MaterialPageRoute(builder: (_) => const QrScannerScreen()),
                    );
                    if (code != null) {
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(content: Text('Đã quét thành công mã: $code')),
                      );
                    }
                  },
                )
              ],
            ),
          ),
        ),
      ],
    );
  }
}
