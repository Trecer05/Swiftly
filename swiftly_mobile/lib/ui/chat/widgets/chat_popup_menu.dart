import 'package:flutter/material.dart';

class ChatPopupMenu extends StatelessWidget {
  final VoidCallback onEdit;
  final VoidCallback onBlock;
  final VoidCallback onDelete;

  const ChatPopupMenu({
    super.key,
    required this.onEdit,
    required this.onBlock,
    required this.onDelete,
  });

  @override
  Widget build(BuildContext context) {
    return PopupMenuButton<int>(
      color: const Color(0xFF1E3A8A),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
      icon: const Icon(Icons.more_vert, color: Colors.white),
      itemBuilder: (context) => [
        _customMenuItem(0, Icons.edit_outlined, "Редактировать"),
        _customMenuItem(1, Icons.block, "Заблокировать пользователя"),
        _customMenuItem(2, Icons.delete, "Удалить чат"),
      ],
      onSelected: (value) {
        switch (value) {
          case 0:
            onEdit();
            break;
          case 1:
            onBlock();
            break;
          case 2:
            onDelete();
            break;
        }
      },
    );
  }

  PopupMenuItem<int> _customMenuItem(int value, IconData icon, String text) {
    return PopupMenuItem<int>(
      value: value,
      child: Row(
        children: [
          Icon(icon, color: Colors.red, size: 20),
          const SizedBox(width: 8),
          Text(
            text,
            style: const TextStyle(
              color: Colors.red,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }
}
