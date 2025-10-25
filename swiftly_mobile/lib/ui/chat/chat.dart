import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart' show kIsWeb, defaultTargetPlatform;
import 'widgets/chat_menu.dart';
import 'widgets/chat_content_panel.dart';
import 'widgets/chat_right_panel.dart';
import 'mob/mob_chat_page.dart';

bool _isDesktopLike(BuildContext context) {
  const double kDesktopBreakpoint = 900;
  final w = MediaQuery.of(context).size.width;
  if (kIsWeb) return w >= kDesktopBreakpoint;
  switch (defaultTargetPlatform) {
    case TargetPlatform.macOS:
    case TargetPlatform.windows:
    case TargetPlatform.linux:
      return w >= kDesktopBreakpoint;
    default:
      return w >= kDesktopBreakpoint;
  }
}

class ChatScreenDesktop extends StatefulWidget {
  const ChatScreenDesktop({super.key});

  @override
  State<ChatScreenDesktop> createState() => _ChatScreenDesktopState();
}

class _ChatScreenDesktopState extends State<ChatScreenDesktop> {
  ChatItem? selectedChat;
  bool showRightPanel = false;

  void _toggleRightPanel() {
    setState(() => showRightPanel = !showRightPanel);
  }

  void _handleMenuAction(String value) {
    switch (value) {
      case 'edit':
        debugPrint("Редактировать нажато");
        break;
      case 'block':
        debugPrint("Заблокировать нажато");
        break;
      case 'delete':
        debugPrint("Удалить чат нажато");
        break;
    }
  }

  @override
  Widget build(BuildContext context) {
    const double rightPanelWidth = 300;

    return Scaffold(
      backgroundColor: Colors.transparent,
      body: Row(
        children: [
          SizedBox(
            child: ChatMenuPanel(
              onChatSelected: (index) {
                final all = [..._pinnedChats, ..._allChats];
                setState(() {
                  selectedChat = all[index];
                });
              },
            ),
          ),
          Expanded(
            child: ChatContentPanel(
              selectedChat: selectedChat,
              onInfoPressed: _toggleRightPanel,
              onMenuPressed: (value) => _handleMenuAction(value),
            ),
          ),

          AnimatedContainer(
            duration: const Duration(milliseconds: 300),
            width: showRightPanel ? rightPanelWidth : 0,
            curve: Curves.easeInOut,
            child: showRightPanel
                ? ChatRightPanel(
                    username: selectedChat?.name ?? 'Username',
                    onClose: _toggleRightPanel,
                  )
                : const SizedBox.shrink(),
          ),
        ],
      ),
    );
  }

  final List<ChatItem> _pinnedChats = [
    ChatItem(
      name: 'Иван Дернов',
      message: '',
      time: '23:02',
      unread: 6,
      tags: ['work', 'swifty'],
    ),
  ];

  final List<ChatItem> _allChats = [
    ChatItem(name: 'Ярослав Хохлов', message: 'До связи!', time: '13:37', unread: 0),
    ChatItem(name: 'Иван Дорн', message: 'Почему дизайнер ничего...', time: 'Tu', unread: 13),
  ];
}

class ChatScreen extends StatelessWidget {
  const ChatScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return _isDesktopLike(context)
        ? const ChatScreenDesktop()
        : const MobileChatPage();
  }
}
