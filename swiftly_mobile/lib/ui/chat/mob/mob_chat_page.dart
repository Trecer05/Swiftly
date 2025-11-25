import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:go_router/go_router.dart';
import 'mob_call_screen.dart';
import '../chat.dart';

import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:swiftly_mobile/data/models/api_chat_client.dart';


const BoxDecoration kChatBackground = BoxDecoration(
  gradient: LinearGradient(
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
    colors: [
      Color(0xFF111121),
      Color(0xFF0D0D1D),
    ],
  ),
);

enum _ChatMenuAction { edit, block, delete }



class MobileChatPage extends StatelessWidget {
  const MobileChatPage({super.key});

  @override
  Widget build(BuildContext context) {
    final chats = const [
      _Chat('1', 'Username', 'вот, теперь нормально', '23:02', 6, true, 'today', true),
      _Chat('2', 'Ярослав Хохлов', 'До связи!', '13:37', 0, false, 'yesterday', false),
      _Chat('3', 'Иван Дорн', 'Почему дизайнер ничего н...', 'Tu', 13, false, 'Mon', false),
    ];

    return Scaffold(
      extendBodyBehindAppBar: true,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        toolbarHeight: 0,
        systemOverlayStyle: const SystemUiOverlayStyle(
          statusBarColor: Colors.transparent,
          statusBarIconBrightness: Brightness.light,
          statusBarBrightness: Brightness.dark,
        ),
      ),
      body: Container(
        padding: EdgeInsets.only(top: MediaQuery.of(context).padding.top),
        decoration: kChatBackground,
        child: Column(
          children: [
            const SizedBox(height: 8),
            const Padding(
              padding: EdgeInsets.symmetric(horizontal: 16),
              child: Align(
                alignment: Alignment.centerLeft,
                child: Text('Чаты', style: TextStyle(fontSize: 32, fontWeight: FontWeight.w800, color: Colors.white)),
              ),
            ),
            const SizedBox(height: 10),
            const Padding(
              padding: EdgeInsets.symmetric(horizontal: 16),
              child: _SearchField(),
            ),
            const SizedBox(height: 6),
            Expanded(
              child: ListView.builder(
                padding: EdgeInsets.zero,
                primary: false,
                itemCount: chats.length,
                itemBuilder: (context, i) {
                  final c = chats[i];
                  return ListTile(
                    onTap: () => Navigator.push(
                      context,
                      MaterialPageRoute(builder: (_) => MobileChatThreadScreen(title: c.title, chatId: c.id)),
                    ),
                    leading: Stack(
                      children: [
                        CircleAvatar(
                          radius: 22,
                          backgroundColor: const Color(0xFF23283B),
                          child: Text(c.title[0].toUpperCase(),
                              style: const TextStyle(fontWeight: FontWeight.w700)),
                        ),
                        if (c.online)
                          Positioned(
                            right: 0,
                            bottom: 0,
                            child: Container(
                              width: 10,
                              height: 10,
                              decoration: BoxDecoration(
                                color: const Color(0xFF4ADE80),
                                border: Border.all(color: const Color(0xFF0F1320), width: 2),
                                shape: BoxShape.circle,
                              ),
                            ),
                          ),
                      ],
                    ),
                    title: Row(
                      children: [
                        Expanded(
                          child: Text(c.title,
                              style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 16, color: Colors.white)),
                        ),
                        Text(c.time,
                            style: TextStyle(color: Colors.white, fontSize: 12)),
                      ],
                    ),
                    subtitle: Padding(
                      padding: const EdgeInsets.only(top: 6),
                      child: Text(
                        c.last,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        style: TextStyle(color: Colors.white),
                      ),
                    ),
                    trailing: c.unread > 0
                        ? Container(
                            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                            decoration: BoxDecoration(
                              color: const Color(0xFF2E63F6),
                              borderRadius: BorderRadius.circular(12),
                            ),
                            child: Text('${c.unread}',
                                style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 12, color: Colors.white )),
                          )
                        : null,
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}


class MobileChatThreadScreen extends StatefulWidget {
  final String title;
  final String chatId;

  const MobileChatThreadScreen({super.key, required this.title, required this.chatId});

  @override
  State<MobileChatThreadScreen> createState() => _MobileChatThreadScreenState();
}

class _MobileChatThreadScreenState extends State<MobileChatThreadScreen> {
  late final ChatApiClient _api;
  late final ChatWsClient _ws;

  final List<MessageDto> _messages = [];
  bool _isTyping = false;
  StreamSubscription? _wsSub;

  @override
  void initState() {
    super.initState();
    _api = ChatApiClient(http.Client());
    _ws = ChatWsClient();

    _loadHistory();
    _connectWs();
  }

  Future<void> _loadHistory() async {
    try {
      // пока чат один, можно захардкодить id или пробросить через конструктор
      final msgs = await _api.getMessages(widget.chatId);
      setState(() {
        _messages
          ..clear()
          ..addAll(msgs);
      });
    } catch (e) {
      debugPrint('Failed to load history: $e');
    }
  }

  Future<void> _connectWs() async {
    try {
      await _ws.connect(widget.chatId);
      _wsSub = _ws.stream?.listen((event) {
        try {
          final data = jsonDecode(event as String) as Map<String, dynamic>;
          final type = data['type'] as String?;
          if (type == 'message') {
            final msg = MessageDto.fromJson(data);
            setState(() => _messages.add(msg));
          } else if (type == 'typing') {
            setState(() => _isTyping = true);
          } else if (type == 'stoptyping') {
            setState(() => _isTyping = false);
          }
        } catch (e) {
          debugPrint('WS parse error: $e');
        }
      });
    } catch (e) {
      debugPrint('WS connect error: $e');
    }
  }

  @override
  void dispose() {
    _wsSub?.cancel();
    _ws.dispose();
    super.dispose();
  }

  void _handleSend(String text) {
    _ws.sendText(text);
    final local = MessageDto(
      id: DateTime.now().millisecondsSinceEpoch.toString(),
      text: text,
      createdAt: DateTime.now(),
      fromId: 'me', // TODO: заменить на currentUserId
    );
    setState(() => _messages.add(local));
  }

  void _handleTypingChanged(bool isTyping) {
    _ws.sendTyping(isTyping);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        systemOverlayStyle: const SystemUiOverlayStyle(
          statusBarColor: Colors.transparent,
          statusBarIconBrightness: Brightness.light,
          statusBarBrightness: Brightness.dark,
        ),
        leading: IconButton(
          icon: const Icon(CupertinoIcons.back, color: Color(0x80FFFFFF)),
          onPressed: () {
            Navigator.pop(context);
          },
        ),
        centerTitle: true,
        title: GestureDetector(
          behavior: HitTestBehavior.opaque,
          onTap: () => Navigator.push(
            context,
            MaterialPageRoute(
              builder: (_) => MobileChatProfileScreen(title: widget.title),
            ),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              const CircleAvatar(
                radius: 14,
                backgroundColor: Color(0xFF23283B),
              ),
              const SizedBox(width: 8),
              Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    widget.title,
                    style: const TextStyle(
                      fontWeight: FontWeight.w700,
                      color: Colors.white,
                      fontSize: 18,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    _isTyping ? 'печатает…' : 'в сети',
                    style: TextStyle(
                      fontSize: 12,
                      color: Colors.white.withOpacity(0.6),
                      height: 1.2,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
        backgroundColor: Colors.transparent,
        flexibleSpace: Container(decoration: kChatBackground),
        elevation: 0,
        scrolledUnderElevation: 0,
        surfaceTintColor: Colors.transparent,
        actions: [
          IconButton(
            icon: const Icon(Icons.call, color: Colors.white),
            onPressed: () {
              context.push('/call', extra: widget.title);
            },
          ),
          PopupMenuButton<_ChatMenuAction>(
            icon: const Icon(
              CupertinoIcons.ellipsis_vertical,
              color: Color(0x80FFFFFF),
            ),
            color: const Color(0xFF171B2A),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
            ),
            onSelected: (value) {
              switch (value) {
                case _ChatMenuAction.edit:
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Редактировать')),
                  );
                  break;
                case _ChatMenuAction.block:
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Заблокировать')),
                  );
                  break;
                case _ChatMenuAction.delete:
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Удалить чат')),
                  );
                  break;
              }
            },
            itemBuilder: (context) => const [
              PopupMenuItem(
                value: _ChatMenuAction.edit,
                child: Row(
                  children: [
                    Icon(Icons.edit, size: 18, color: Colors.white),
                    SizedBox(width: 8),
                    Text('Редактировать',
                        style: TextStyle(color: Colors.white)),
                  ],
                ),
              ),
              PopupMenuItem(
                value: _ChatMenuAction.block,
                child: Row(
                  children: [
                    Icon(Icons.block, size: 18, color: Colors.white),
                    SizedBox(width: 8),
                    Text('Заблокировать чат',
                        style: TextStyle(color: Colors.white)),
                  ],
                ),
              ),
              PopupMenuItem(
                value: _ChatMenuAction.delete,
                child: Row(
                  children: [
                    Icon(Icons.delete, size: 18, color: Colors.red),
                    SizedBox(width: 8),
                    Text('Удалить чат',
                        style: TextStyle(color: Colors.red)),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(width: 8),
        ],
      ),
      body: Container(
        decoration: kChatBackground,
        child: Column(
          children: [
            const Divider(height: 1, color: Color(0x22FFFFFF)),
            Expanded(
              child: ListView.builder(
                padding: const EdgeInsets.fromLTRB(12, 12, 12, 12),
                reverse: true,
                itemCount: _messages.length,
                itemBuilder: (context, index) {
                  final m = _messages[_messages.length - 1 - index];
                  final isMine = m.fromId == 'me'; // TODO: сравнить с currentUserId

                  const bg = Color(0x0FFFFFFF);
                  final radius = BorderRadius.only(
                    topLeft: const Radius.circular(16),
                    topRight: const Radius.circular(16),
                    bottomLeft: Radius.circular(isMine ? 16 : 4),
                    bottomRight: Radius.circular(isMine ? 4 : 16),
                  );

                  return Padding(
                    padding: const EdgeInsets.symmetric(vertical: 6),
                    child: Align(
                      alignment: isMine
                          ? Alignment.centerRight
                          : Alignment.centerLeft,
                      child: ConstrainedBox(
                        constraints:
                            const BoxConstraints(maxWidth: 320),
                        child: DecoratedBox(
                          decoration: BoxDecoration(
                            color: bg,
                            borderRadius: radius,
                          ),
                          child: Padding(
                            padding: const EdgeInsets.fromLTRB(
                                12, 10, 12, 8),
                            child: Column(
                              crossAxisAlignment: isMine
                                  ? CrossAxisAlignment.end
                                  : CrossAxisAlignment.start,
                              children: [
                                Text(
                                  m.text,
                                  style: const TextStyle(
                                    fontSize: 15,
                                    height: 1.25,
                                    color: Colors.white,
                                  ),
                                ),
                                const SizedBox(height: 4),
                                Text(
                                  // можно красиво форматнуть время
                                  '${m.createdAt.hour.toString().padLeft(2, '0')}:${m.createdAt.minute.toString().padLeft(2, '0')}',
                                  style: TextStyle(
                                    color: Colors.white.withOpacity(.75),
                                    fontSize: 11,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ),
                  );
                },
              ),
            ),
            _InputBar(
              onSend: _handleSend,
              onTypingChanged: _handleTypingChanged,
            ),
          ],
        ),
      ),
    );
  }
}

class _SearchField extends StatelessWidget {
  const _SearchField();

  @override
  Widget build(BuildContext context) {
    return TextField(
      style: const TextStyle(color: Colors.white),
      decoration: InputDecoration(
        filled: true,
        fillColor: const Color(0xFF171B2A),
        hintText: 'Поиск',
        hintStyle: TextStyle(color: Colors.white.withOpacity(.6)),
        contentPadding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: BorderSide.none,
        ),
        prefixIcon: const Icon(CupertinoIcons.search),
      ),
    );
  }
}

class _InputBar extends StatefulWidget {
  final ValueChanged<String> onSend;
  final ValueChanged<bool>? onTypingChanged;

  const _InputBar({
    required this.onSend,
    this.onTypingChanged,
  });

  @override
  State<_InputBar> createState() => _InputBarState();
}

class _InputBarState extends State<_InputBar> {
  final _controller = TextEditingController();
  Timer? _typingTimer;

  void _handleChanged(String value) {
    widget.onTypingChanged?.call(true);
    _typingTimer?.cancel();
    _typingTimer = Timer(const Duration(seconds: 2), () {
      widget.onTypingChanged?.call(false);
    });
  }

  void _handleSend() {
    final text = _controller.text.trim();
    if (text.isEmpty) return;
    widget.onSend(text);
    _controller.clear();
    widget.onTypingChanged?.call(false);
  }

  @override
  void dispose() {
    _typingTimer?.cancel();
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      top: false,
      child: Container(
        decoration: const BoxDecoration(
          color: Color(0xFF111121),
          border: Border(
            top: BorderSide(
              color: Color(0x22FFFFFF),
              width: 1,
            ),
          ),
        ),
        padding: const EdgeInsets.fromLTRB(12, 8, 12, 8),
        child: Row(
          children: [
            Container(
              width: 36,
              height: 36,
              decoration: BoxDecoration(
                color: const Color(0xFF23283B),
                borderRadius: BorderRadius.circular(12),
              ),
              child: IconButton(
                padding: EdgeInsets.zero,
                icon: const Icon(
                  CupertinoIcons.add,
                  size: 20,
                  color: Colors.white,
                ),
                onPressed: () {
                  // сюда позже добавим вложения
                },
              ),
            ),
            const SizedBox(width: 8),
            Expanded(
              child: Container(
                decoration: BoxDecoration(
                  color: const Color(0xFF171B2A),
                  borderRadius: BorderRadius.circular(16),
                ),
                padding: const EdgeInsets.symmetric(horizontal: 12),
                child: Row(
                  children: [
                    Expanded(
                      child: TextField(
                        controller: _controller,
                        onChanged: _handleChanged,
                        maxLines: null,
                        keyboardType: TextInputType.multiline,
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 15,
                        ),
                        decoration: const InputDecoration(
                          hintText: 'Сообщение',
                          hintStyle: TextStyle(
                            color: Color(0x80FFFFFF),
                            fontSize: 15,
                          ),
                          border: InputBorder.none,
                          isCollapsed: true,
                        ),
                      ),
                    ),
                    IconButton(
                      padding: EdgeInsets.zero,
                      icon: const Icon(
                        CupertinoIcons.smiley,
                        size: 22,
                        color: Color(0x80FFFFFF),
                      ),
                      onPressed: () {
                        // сюда потом эмодзи
                      },
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(width: 8),
            Container(
              width: 36,
              height: 36,
              decoration: BoxDecoration(
                gradient: const LinearGradient(
                  colors: [
                    Color(0xFF6C5CE7),
                    Color(0xFF9B59B6),
                  ],
                ),
                borderRadius: BorderRadius.circular(12),
              ),
              child: IconButton(
                padding: EdgeInsets.zero,
                icon: const Icon(
                  CupertinoIcons.paperplane_fill,
                  size: 18,
                  color: Colors.white,
                ),
                onPressed: _handleSend,
              ),
            ),
          ],
        ),
      ),
    );
  }
}


class _Chat {
  final String id;
  final String title;
  final String subtitle;
  final String time;
  final int unread;
  final bool pinned;
  final String last;
  final bool online;

  const _Chat(
    this.id,
    this.title,
    this.subtitle,
    this.time,
    this.unread,
    this.pinned,
    this.last,
    this.online,
  );
}


class MobileChatProfileScreen extends StatefulWidget {
  final String title;
  const MobileChatProfileScreen({super.key, required this.title});

  @override
  State<MobileChatProfileScreen> createState() => _MobileChatProfileScreenState();
}

class _MobileChatProfileScreenState extends State<MobileChatProfileScreen> {
  bool _notificationsEnabled = true;
  int _tabIndex = 0;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(CupertinoIcons.back, color: Colors.white,),
          onPressed: () => Navigator.pop(context),
        ),
        systemOverlayStyle: const SystemUiOverlayStyle(
          statusBarColor: Colors.transparent,
          statusBarIconBrightness: Brightness.light,
          statusBarBrightness: Brightness.dark,
        ),
        backgroundColor: Colors.transparent,
        flexibleSpace: Container(decoration: kChatBackground),
        elevation: 0,
        scrolledUnderElevation: 0,
        surfaceTintColor: Colors.transparent,
      ),
      body: Container(
        decoration: kChatBackground,
        child: ListView(
          padding: const EdgeInsets.fromLTRB(16, 8, 16, 24),
          children: [
            const SizedBox(height: 8),
            Center(
              child: CircleAvatar(
                radius: 40,
                backgroundColor: const Color(0xFF23283B),
                child: Text(
                  widget.title.isNotEmpty ? widget.title[0].toUpperCase() : '?',
                  style: const TextStyle(fontSize: 28, fontWeight: FontWeight.w700),
                ),
              ),
            ),
            const SizedBox(height: 12),
            Center(
              child: Text(
                widget.title,
                style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w800, color: Colors.white),
              ),
            ),
            const SizedBox(height: 4),
            Center(
              child: Text(
                'в сети',
                style: TextStyle(color: Colors.white.withOpacity(.6), fontSize: 13),
              ),
            ),
            const SizedBox(height: 14),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: const [
                _SquareActionButton(
                  icon: CupertinoIcons.chat_bubble_text,
                  label: 'Чаты',
                  onTap: _noop,
                ),
                SizedBox(width: 12),
                _SquareActionButton(
                  icon: CupertinoIcons.phone,
                  label: 'Звонок',
                  onTap: _noop,
                ),
              ],
            ),
            const SizedBox(height: 16),
            Container(
              decoration: BoxDecoration(
                color: const Color(0xFF171B2A),
                borderRadius: BorderRadius.circular(16),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: const [
                  _InfoBlock(
                    title: 'О себе',
                    value: 'UX/UI Designer в Swiftly',
                  ),
                  _CardDivider(),
                  _InfoBlock(
                    title: 'Грейд/Уровень',
                    value: 'Junior/Junior +',
                  ),
                  _CardDivider(),
                  _InfoBlock(
                    title: 'Имя пользователя',
                    value: '@cutecrysta1',
                  ),
                ],
              ),
            ),
            Container(
              decoration: BoxDecoration(
                color: const Color(0xFF171B2A),
                borderRadius: BorderRadius.circular(16),
              ),
              margin: const EdgeInsets.only(top: 8),
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: Row(
                  children: [
                    const Icon(CupertinoIcons.bell, size: 18, color: Colors.white,),
                    const SizedBox(width: 10),
                    const Expanded(
                      child: Text(
                        'Уведомления',
                        style: TextStyle(fontWeight: FontWeight.w600),
                      ),
                    ),
                    Switch(
                      value: _notificationsEnabled,
                      onChanged: (v) => setState(() => _notificationsEnabled = v),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 18),
            Container(
              padding: const EdgeInsets.all(6),
              decoration: BoxDecoration(
                color: const Color(0xFF171B2A),
                borderRadius: BorderRadius.circular(14),
              ),
              child: Row(
                children: [
                  _Segment(
                    text: 'Медиа',
                    selected: _tabIndex == 0,
                    onTap: () => setState(() => _tabIndex = 0),
                  ),
                  _Segment(
                    text: 'Файлы',
                    selected: _tabIndex == 1,
                    onTap: () => setState(() => _tabIndex = 1),
                  ),
                  _Segment(
                    text: 'Аудио',
                    selected: _tabIndex == 2,
                    onTap: () => setState(() => _tabIndex = 2),
                  ),
                  _Segment(
                    text: 'Ссылки',
                    selected: _tabIndex == 3,
                    onTap: () => setState(() => _tabIndex = 3),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 12),
            Container(
              height: 180,
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color: const Color(0xFF0F1320),
                borderRadius: BorderRadius.circular(14),
                border: Border.all(color: const Color(0x22FFFFFF)),
              ),
              child: Text(
                _tabIndex == 0
                    ? 'Медиа пока нет'
                    : _tabIndex == 1
                        ? 'Файлы пока нет'
                        : _tabIndex == 2
                            ? 'Аудио пока нет'
                            : 'Ссылок пока нет',
                style: TextStyle(color: Colors.white.withOpacity(.7)),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

void _noop() {}

class _SquareActionButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;
  const _SquareActionButton({required this.icon, required this.label, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
        decoration: BoxDecoration(
          color: const Color(0xFF171B2A),
          borderRadius: BorderRadius.circular(14),
        ),
        child: Row(
          children: [
            Icon(icon, size: 18, color: Colors.white,),
            const SizedBox(width: 8),
            Text(label, style: const TextStyle(fontWeight: FontWeight.w600, color: Colors.white)),
          ],
        ),
      ),
    );
  }
}

class _InfoBlock extends StatelessWidget {
  final String title;
  final String value;
  const _InfoBlock({required this.title, required this.value});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(title, style: TextStyle(color: Colors.white.withOpacity(.6), fontSize: 12)),
          const SizedBox(height: 6),
          Text(value, style: const TextStyle(fontWeight: FontWeight.w600, color: Colors.white)),
        ],
      ),
    );
  }
}

class _CardDivider extends StatelessWidget {
  const _CardDivider();
  @override
  Widget build(BuildContext context) {
    return const Divider(height: 1, color: Color(0x22FFFFFF), indent: 16, endIndent: 16);
  }
}

class _Segment extends StatelessWidget {
  final String text;
  final bool selected;
  final VoidCallback onTap;

  const _Segment({required this.text, required this.selected, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: GestureDetector(
        onTap: onTap,
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 8),
          decoration: BoxDecoration(
            color: selected ? const Color(0xFF2E63F6) : Colors.transparent,
            borderRadius: BorderRadius.circular(12),
          ),
          alignment: Alignment.center,
          child: Text(
            text,
            style: TextStyle(
              color: selected ? Colors.white : Colors.white.withOpacity(0.6),
              fontWeight: FontWeight.w600,
            ),
          ),
        ),
      ),
    );
  }
}