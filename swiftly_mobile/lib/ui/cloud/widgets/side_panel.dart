import 'package:flutter/material.dart';

/// Боковая панель навигации
class SidePanel extends StatelessWidget {
  const SidePanel({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Colors.transparent,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SidePanelItem(
            icon: Icons.folder,
            title: 'Все файлы',
          ),
          SidePanelItem(
            icon: Icons.folder,
            title: 'Последние',
          ),
          SidePanelItem(
            icon: Icons.folder,
            title: 'По проектам',
          ),
          SidePanelItem(
            icon: Icons.folder,
            title: 'Избранное',
          ),
          SidePanelItem(
            icon: Icons.folder,
            title: 'По задачам',
          ),
        ],
      ),
    );
  }
}

/// Элемент боковой панели с иконкой
class SidePanelItem extends StatelessWidget {
  final IconData icon;
  final String title;

  const SidePanelItem({
    required this.icon,
    required this.title,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(left: 18, top: 5, bottom: 5),
      child: Row(
        children: [
          Icon(
            icon,
            size: 18,
            color: Colors.grey[400],
          ),
          SizedBox(width: 10),
          Text(
            title,
            style: TextStyle(
              color: Colors.grey[400],
              fontSize: 14,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }
}
