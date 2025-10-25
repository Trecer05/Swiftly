import 'package:flutter/material.dart';

class ChatPopupMenu extends StatelessWidget {
  final VoidCallback onEdit;
  final VoidCallback onBlock;
  final VoidCallback onDelete;

  // Пер-элементные параметры (по умолчанию подобраны под UX)
  final Color editColor;
  final Color blockColor;
  final Color deleteColor;

  const ChatPopupMenu({
    super.key,
    required this.onEdit,
    required this.onBlock,
    required this.onDelete,
    this.editColor = const Color(0xFFEAEAEA),
    this.blockColor = const Color(0xFFEAEAEA), 
    this.deleteColor = const Color(0xFFE53935),
  });

  @override
  Widget build(BuildContext context) {
    return PopupMenuButton<int>(
      color: const Color(0x0A535353),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
      icon: const Icon(Icons.more_vert, color: Colors.white),
      itemBuilder: (context) => [
        _customMenuItem(0, Icons.edit_outlined, "Редактировать", editColor),
        _customMenuItem(1, Icons.block, "Заблокировать пользователя", blockColor),
        _customMenuItem(2, Icons.delete, "Удалить чат", deleteColor),
      ],
      onSelected: (value) {
        switch (value) {
          case 0:
            _handleEdit(context);
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

  void _handleEdit(BuildContext context) {
    // TODO: вставь сюда свою реализацию редактирования
    onEdit();
  }

  PopupMenuItem<int> _customMenuItem(int value, IconData icon, String text, Color color) {
    return PopupMenuItem<int>(
      value: value,
      child: Row(
        children: [
          Icon(icon, color: color, size: 20),
          const SizedBox(width: 8),
          Text(
            text,
            style: TextStyle(
              color: color,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }
}
