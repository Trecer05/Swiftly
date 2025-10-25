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

  static const double expandedWidth = 300;
  static const double collapsedWidth = 80;
  bool isCollapsed = false;

  void _onDragUpdate(DragUpdateDetails details) {
    if (details.delta.dx < -6) {
      if (!isCollapsed) setState(() => isCollapsed = true);
    } else if (details.delta.dx > 6) {
      if (isCollapsed) setState(() => isCollapsed = false);
    }
  }

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
    final double screenWidth = MediaQuery.of(context).size.width;

    return ClipRRect(
      borderRadius: BorderRadius.zero,
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: 20, sigmaY: 20),
        child: Stack(
          children: [
            AnimatedContainer(
              duration: const Duration(milliseconds: 160),
              curve: Curves.easeOut,
              width: isCollapsed ? collapsedWidth : expandedWidth,
              padding: const EdgeInsets.only(top: 20),
              decoration: const BoxDecoration(
                color: Color(0x8C080808),
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
                          hintStyle: TextStyle(color: Colors.white70),
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
                  if (!isCollapsed) const SizedBox(height: 12),

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
                      padding: isCollapsed ? EdgeInsets.zero : const EdgeInsets.symmetric(vertical: 8),
                      itemExtent: isCollapsed ? (collapsedWidth * 0.6 + 16) : null,
                      itemCount: allChats.length,
                      itemBuilder: (context, index) {
                        return _buildChatItem(
                          allChats[index],
                          index + pinnedChats.length,
                          isCollapsed,
                        );
                      },
                    ),
                  ),
                ],
              ),
            ),

            // Хэндл для ресайза справа
            Positioned(
              right: 0,
              top: 0,
              bottom: 0,
              child: MouseRegion(
                cursor: SystemMouseCursors.resizeLeftRight,
                child: GestureDetector(
                  behavior: HitTestBehavior.translucent,
                  onHorizontalDragUpdate: _onDragUpdate,
                  onDoubleTap: () {
                    setState(() {
                      isCollapsed = !isCollapsed;
                    });
                  },
                  child: const SizedBox(width: 8),
                ),
              ),
            ),
          ],
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
        width: double.infinity,
        decoration: BoxDecoration(
          color: isSelected ? Colors.white.withOpacity(0.1) : Colors.transparent,
          borderRadius: BorderRadius.circular(8),
        ),
        child: Padding(
          padding: EdgeInsets.symmetric(
            horizontal: isCollapsed ? collapsedWidth * 0.2 : 12,
            vertical: isCollapsed ? 8 : 12,
          ),
          child: Row(
            children: [
            if (isCollapsed)
              SizedBox(
                width: collapsedWidth * 0.6,
                height: collapsedWidth * 0.6,
                child: Center(
                  child: CircleAvatar(
                    radius: (collapsedWidth * 0.6) / 2,
                    backgroundColor: chat.avatarColor ?? Colors.grey.shade700,
                    child: chat.avatarColor == null
                        ? Text(chat.name[0], style: const TextStyle(color: Colors.white))
                        : null,
                  ),
                ),
              )
            else
              CircleAvatar(
                radius: 20,
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
