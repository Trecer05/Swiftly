import 'package:flutter/material.dart';
import 'chat_menu.dart';
import 'chat_popup_menu.dart';

class ChatContentPanel extends StatelessWidget {
  final ChatItem? selectedChat;
  final VoidCallback onInfoPressed;
  final ValueChanged<String> onMenuPressed;

  const ChatContentPanel({
    super.key,
    this.selectedChat,
    required this.onInfoPressed,
    required this.onMenuPressed,
  });

  @override
  Widget build(BuildContext context) {
    if (selectedChat == null) {
      return Expanded(
        child: Container(
          decoration: _blueGradient(),
          child: Center(
            child: Text(
              'Выберите чат слева',
              style: TextStyle(color: Colors.white.withOpacity(0.85), fontSize: 18),
            ),
          ),
        ),
      );
    }

    return Expanded(
      child: Column(
        children: [
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            decoration: _blueGradient(opacity: 0.3),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  selectedChat!.name,
                  style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Colors.white),
                ),
                Row(
                  children: [
                    IconButton(
                      icon: const Icon(Icons.search, color: Colors.white, size: 16,),
                      onPressed: () {},
                    ),
                    IconButton(
                      icon: Icon(Icons.call, color: Colors.white, size: 16,),
                      onPressed: () {},
                    ),
                    IconButton(
                      icon: Image.asset('assets/chat_bar_icon.png', color: Colors.white, width: 16,),
                      onPressed: onInfoPressed
                    ),
                    ChatPopupMenu(
                      onEdit: () => onMenuPressed('edit'),
                      onBlock: () => onMenuPressed('block'),
                      onDelete: () => onMenuPressed('delete'),
                    ),
                  ],
                )
              ],
            ),
          ),
          Expanded(
            child: Container(
              decoration: _blueGradient(opacity: 0.15),
              child: Center(
                child: Text(
                  'Напишите первое сообщение',
                  style: TextStyle(color: Colors.white.withOpacity(0.85), fontSize: 16),
                ),
              ),
            ),
          ),
          Container(
            padding: const EdgeInsets.all(8),
            decoration: _blueGradient(opacity: 0.25),
            child: Row(
              children: [
                IconButton(icon: const Icon(Icons.attach_file, color: Colors.white), onPressed: () {}),
                const Expanded(
                  child: TextField(
                    style: TextStyle(color: Colors.white),
                    decoration: InputDecoration(
                      hintText: 'Сообщение...',
                      hintStyle: TextStyle(color: Colors.white70),
                      border: InputBorder.none,
                    ),
                  ),
                ),
                IconButton(icon: const Icon(Icons.send, color: Colors.white), onPressed: () {}),
              ],
            ),
          ),
        ],
      ),
    );
  }

  BoxDecoration _blueGradient({double opacity = 0.2}) {
    return BoxDecoration(
      gradient: LinearGradient(
        colors: [
          Colors.blue.shade900.withOpacity(opacity),
          Colors.blue.shade700.withOpacity(opacity - 0.05),
          Colors.blue.shade600.withOpacity(opacity - 0.07),
        ],
        begin: Alignment.topLeft,
        end: Alignment.bottomRight,
      ),
    );
  }
}

