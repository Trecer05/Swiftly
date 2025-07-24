import 'dart:ui';
import 'package:flutter/material.dart';

class ChatMenuPanel extends StatefulWidget {
  final ValueChanged<int>? onChatSelected;

  const ChatMenuPanel({super.key, this.onChatSelected});

  @override
  State<ChatMenuPanel> createState() => _ChatMenuPanelState();
}

class _ChatMenuPanelState extends State<ChatMenuPanel> {
  int selectedChatIndex = -1;

  final List<ChatItem> pinnedChats = [
    ChatItem(
      name: 'Иван Дернов',
      message: 'Say my name',
      time: '23:02',
      unread: 6,
      tags: ['work', 'swifty'],
      avatarColor: Colors.white,
    ),
  ];

  final List<ChatItem> allChats = [
    ChatItem(
      name: 'Ярослав Хохлов',
      message: 'Heisenberg',
      time: '13:37',
      unread: 0,
    ),
    ChatItem(
      name: 'Иван Дорн',
      message: 'you\'re goddamn right',
      time: 'Tu',
      unread: 13,
    ),
  ];

  @override
  Widget build(BuildContext context) {
    final double width = MediaQuery.of(context).size.width;
    final bool isCollapsed = width < 900; 

    return ClipRRect(
      borderRadius: BorderRadius.zero,
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: 20, sigmaY: 20),
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 300),
          width: isCollapsed ? 80 : 300,
          padding: const EdgeInsets.only(top: 20),
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              colors: [
                Color(0xAA0D1B3B),
                Color(0xAA0B1730),
                Color(0x66081824),
              ],
              begin: Alignment.topCenter,
              end: Alignment.bottomCenter,
            ),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              if (!isCollapsed)
                const Padding(
                  padding: EdgeInsets.symmetric(horizontal: 16.0),
                  child: Text(
                    'Главная',
                    style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold, color: Colors.white),
                  ),
                ),
              if (!isCollapsed) const SizedBox(height: 12),

              if (!isCollapsed)
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  child: TextField(
                    style: const TextStyle(color: Colors.white),
                    decoration: InputDecoration(
                      hintText: 'Поиск',
                      hintStyle: TextStyle(color: Colors.grey.shade400),
                      filled: true,
                      fillColor: Colors.black26,
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(8),
                        borderSide: BorderSide.none,
                      ),
                      prefixIcon: const Icon(Icons.search, color: Colors.white),
                    ),
                  ),
                ),
              if (!isCollapsed) const SizedBox(height: 20),

              if (!isCollapsed && pinnedChats.isNotEmpty)
                const Padding(
                  padding: EdgeInsets.symmetric(horizontal: 16, vertical: 4),
                  child: Text(
                    'ЗАКРЕПЛЕННЫЕ',
                    style: TextStyle(fontSize: 12, color: Colors.grey),
                  ),
                ),
              ...pinnedChats.asMap().entries.map((entry) {
                final i = entry.key;
                final chat = entry.value;
                return _buildChatItem(chat, i, isCollapsed, isPinned: true);
              }),

              if (!isCollapsed && allChats.isNotEmpty)
                const Padding(
                  padding: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                  child: Text(
                    'ВСЕ ЧАТЫ',
                    style: TextStyle(fontSize: 12, color: Colors.grey),
                  ),
                ),
              Expanded(
                child: ListView.builder(
                  itemCount: allChats.length,
                  itemBuilder: (context, index) {
                    return _buildChatItem(
                        allChats[index], index + pinnedChats.length, isCollapsed);
                  },
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildChatItem(ChatItem chat, int index, bool isCollapsed, {bool isPinned = false}) {
    final isSelected = selectedChatIndex == index;
    return InkWell(
      onTap: () {
        setState(() {
          selectedChatIndex = index;
        });
        widget.onChatSelected?.call(index);
      },
      child: Container(
        margin: EdgeInsets.symmetric(horizontal: isCollapsed ? 8 : 12, vertical: 4),
        padding: EdgeInsets.all(isCollapsed ? 8 : 12),
        decoration: BoxDecoration(
          color: isSelected ? Colors.white.withOpacity(0.1) : Colors.transparent,
          borderRadius: BorderRadius.circular(8),
        ),
        child: Row(
          children: [
            CircleAvatar(
              radius: isCollapsed ? 16 : 20,
              backgroundColor: chat.avatarColor ?? Colors.grey.shade700,
              child: chat.avatarColor == null
                  ? Text(chat.name[0], style: const TextStyle(color: Colors.white))
                  : null,
            ),
            if (!isCollapsed) ...[
              const SizedBox(width: 10),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(chat.name,
                        style: const TextStyle(
                            fontWeight: FontWeight.bold, color: Colors.white)),
                    if (chat.message.isNotEmpty) ...[
                      const SizedBox(height: 4),
                      Text(
                        chat.message,
                        style: const TextStyle(color: Colors.white, fontSize: 12),
                        overflow: TextOverflow.ellipsis,
                      ),
                    ],
                    if (isPinned && chat.tags.isNotEmpty)
                      Padding(
                        padding: const EdgeInsets.only(top: 4),
                        child: Wrap(
                          spacing: 4,
                          children: chat.tags
                              .map((tag) => Container(
                                    padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                                    decoration: BoxDecoration(
                                      color: tag == 'work' ? Colors.green : Colors.orange,
                                      borderRadius: BorderRadius.circular(6),
                                    ),
                                    child: Text(
                                      tag,
                                      style: const TextStyle(fontSize: 10, color: Colors.white),
                                    ),
                                  ))
                              .toList(),
                        ),
                      )
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(chat.time, style: const TextStyle(fontSize: 12, color: Colors.grey)),
                  if (chat.unread > 0)
                    Container(
                      margin: const EdgeInsets.only(top: 4),
                      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                      decoration: BoxDecoration(
                        color: Colors.blue,
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: Text(
                        '${chat.unread}',
                        style: const TextStyle(fontSize: 10, color: Colors.white),
                      ),
                    ),
                ],
              ),
            ]
          ],
        ),
      ),
    );
  }
}

class ChatItem {
  final String name;
  final String message;
  final String time;
  final int unread;
  final List<String> tags;
  final Color? avatarColor;

  ChatItem({
    required this.name,
    required this.message,
    required this.time,
    required this.unread,
    this.tags = const [],
    this.avatarColor,
  });
}
