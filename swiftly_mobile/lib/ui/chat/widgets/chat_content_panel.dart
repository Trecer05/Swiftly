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
          color: Color(0xBF080808),
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
            color: Color(0x01080808),
            child: Row(
              children: [
                Expanded(
                  child: Text(
                    selectedChat!.name,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    softWrap: false,
                    style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Colors.white),
                  ),
                ),
                Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    IconButton(
                      icon: const Icon(Icons.search, color: Colors.white, size: 16),
                      padding: EdgeInsets.zero,
                      constraints: const BoxConstraints(minWidth: 36, minHeight: 36),
                      onPressed: () {},
                    ),
                    IconButton(
                      icon: const Icon(Icons.call, color: Colors.white, size: 16),
                      padding: EdgeInsets.zero,
                      constraints: const BoxConstraints(minWidth: 36, minHeight: 36),
                      onPressed: () {},
                    ),
                    IconButton(
                      icon: Image.asset('assets/chat_bar_icon.png', color: Colors.white, width: 16),
                      padding: EdgeInsets.zero,
                      constraints: const BoxConstraints(minWidth: 36, minHeight: 36),
                      onPressed: onInfoPressed,
                    ),
                    ChatPopupMenu(
                      onEdit: () => onMenuPressed('edit'),
                      onBlock: () => onMenuPressed('block'),
                      onDelete: () => onMenuPressed('delete'),
                    ),
                  ],
                ),
              ],
            ),
          ),
          Expanded(
            child: Container(
              color: Color(0xBF080808),
              child: Center(
                child: Text(
                  'Напишите первое сообщение',
                  style: TextStyle(color: Colors.white.withOpacity(0.85), fontSize: 16),
                ),
              ),
            ),
          ),
          Container(
            color: Color(0xBF080808), // фон темы за плашкой
            padding: const EdgeInsets.symmetric(vertical: 12),
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 12),
              margin: const EdgeInsets.symmetric(horizontal: 16),
              decoration: BoxDecoration(
                color: Color(0x06FFFFFF),
                borderRadius: BorderRadius.circular(20),
              ),
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
          )
        ],
      ),
    );
  }
}
