import 'dart:ui';
import 'package:flutter/material.dart';

class ChatRightPanel extends StatelessWidget {
  final String username;
  final VoidCallback onClose;

  const ChatRightPanel({
    super.key,
    required this.username,
    required this.onClose,
  });

  Widget _buildInfoButton(IconData icon, String title, {Widget? trailing}) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6.0, horizontal: 12.0),
      child: Row(
        children: [
          Icon(icon, color: Colors.white, size: 22),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              title,
              style: const TextStyle(color: Colors.white, fontSize: 16),
            ),
          ),
          if (trailing != null) trailing,
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final double w = constraints.maxWidth;
        final bool showProfile = w >= 140; // прячем аватар/username на очень узких ширинах во время анимации

        return ClipRRect(
          borderRadius: BorderRadius.zero,
          child: BackdropFilter(
            filter: ImageFilter.blur(sigmaX: 20, sigmaY: 20),
            child: Container(
              clipBehavior: Clip.hardEdge, // не даём детям рисовать за пределами при узкой ширине
              decoration: const BoxDecoration(
                color: Color(0x8C080808),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Заголовок
                  Padding(
                    padding: const EdgeInsets.all(16),
                    child: Row(
                      children: [
                        // Текст тянется и обрезается, чтобы не было overflow при узкой ширине
                        const Expanded(
                          child: Text(
                            'Информация',
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            softWrap: false,
                            style: TextStyle(
                              color: Colors.white,
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                        IconButton(
                          icon: const Icon(Icons.close, color: Colors.white),
                          onPressed: onClose,
                          padding: EdgeInsets.zero,
                          constraints: const BoxConstraints(minWidth: 36, minHeight: 36),
                        ),
                      ],
                    ),
                  ),

                  if (showProfile)
                    Center(
                      child: Column(
                        children: [
                          const CircleAvatar(
                            radius: 40,
                            backgroundColor: Colors.grey,
                          ),
                          const SizedBox(height: 12),
                          Text(
                            username,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            softWrap: false,
                            style: const TextStyle(
                              color: Colors.white,
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const Text(
                            '@username',
                            style: TextStyle(color: Colors.grey, fontSize: 14),
                          ),
                          const SizedBox(height: 16),
                        ],
                      ),
                    ),

                  // Контент
                  Expanded(
                    child: SingleChildScrollView(
                      padding: const EdgeInsets.symmetric(vertical: 12),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          _buildInfoButton(
                            Icons.notifications,
                            'Уведомления',
                            trailing: Switch(
                              value: true,
                              onChanged: (v) {},
                              activeColor: Colors.blue,
                            ),
                          ),
                          _buildInfoButton(Icons.image, 'Изображения'),
                          _buildInfoButton(Icons.videocam, 'Видео'),
                          _buildInfoButton(Icons.music_note, 'Аудиофайлы'),
                          _buildInfoButton(Icons.link, 'Ссылки'),
                          _buildInfoButton(Icons.insert_drive_file, 'Файлы'),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}
