import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class MobileCallScreen extends StatefulWidget {
  final String username;
  /// false — исходящий, true — входящий
  final bool isIncoming;
  final ImageProvider? remoteImage;
  final ImageProvider? localImage;

  const MobileCallScreen({
    super.key,
    required this.username,
    this.isIncoming = false,
    this.remoteImage,
    this.localImage,
  });

  @override
  State<MobileCallScreen> createState() => _MobileCallScreenState();
}

class _MobileCallScreenState extends State<MobileCallScreen> {
  bool isMuted = false;
  bool isSpeakerOn = false;
  bool isVideoOn = true;

  /// Пока true — экран во входящем режиме (только принять/отклонить)
  bool _isRinging = false;

  bool _controlsVisible = true;

  Duration _callDuration = Duration.zero;
  Timer? _timer;

  @override
  void initState() {
    super.initState();

    _isRinging = widget.isIncoming;

    // Фуллскрин без системной навигации.
    SystemChrome.setEnabledSystemUIMode(SystemUiMode.immersiveSticky); // [1][web:84]

    // --- Заготовки под API/RTC ---

    // late final CallApi _callApi;
    // late final RemoteVideoStream _remoteStream;
    // late final LocalVideoStream _localStream;
    //
    // Future<void> _initVideoCall() async {
    //   // TODO: поднять WebRTC/SDK, подписаться на remote/local стримы
    //   // _remoteStream = await _callApi.subscribeRemote(...);
    //   // _localStream  = await _callApi.getLocalStream();
    // }
    //
    // if (!widget.isIncoming) {
    //   // TODO: исходящий вызов
    //   // await _callApi.startCall(to: widget.username);
    // } else {
    //   // TODO: показать входящий (пуш, CallKit и т.д.)
    // }
    //
    // _initVideoCall();

    // Таймер сразу запускаем только для исходящих/активных
    if (!widget.isIncoming) {
      _startTimer();
    }
  }

  @override
  void dispose() {
    _timer?.cancel();

    // Вернуть системные бары.
    SystemChrome.setEnabledSystemUIMode(
      SystemUiMode.edgeToEdge,
      overlays: SystemUiOverlay.values,
    ); // [2][web:84]

    super.dispose();
  }

  void _startTimer() {
    _timer ??= Timer.periodic(const Duration(seconds: 1), (_) {
      setState(() => _callDuration += const Duration(seconds: 1));
    });
  }

  String _formatDuration(Duration d) {
    final m = d.inMinutes.remainder(60).toString().padLeft(2, '0');
    final s = d.inSeconds.remainder(60).toString().padLeft(2, '0');
    return '$m:$s';
  }

  String _buildInitials(String name) {
    final parts = name.trim().split(RegExp(r'\s+'));
    if (parts.length == 1 && parts.first.isNotEmpty) {
      final n = parts.first;
      return (n.length >= 2 ? n.substring(0, 2) : n[0]).toUpperCase();
    }
    final buf = StringBuffer();
    for (final p in parts.take(2)) {
      if (p.isNotEmpty) buf.write(p[0].toUpperCase());
    }
    return buf.toString();
  }

  void _toggleMute() {
    setState(() => isMuted = !isMuted);
    // _callApi.setMuted(isMuted);
  }

  void _toggleSpeaker() {
    setState(() => isSpeakerOn = !isSpeakerOn);
    // _callApi.setSpeaker(isSpeakerOn);
  }

  void _toggleVideo() {
    setState(() => isVideoOn = !isVideoOn);
    // _callApi.setVideoEnabled(isVideoOn);
  }

  void _endCall() {
    // _callApi.endCall();
    Navigator.of(context).pop();
  }

  void _acceptIncoming() {
    setState(() {
      _isRinging = false;
      // можно пометить звонок как «принят» и т.п.
    });
    _startTimer();
    // _callApi.acceptCall(); // заготовка
  }

  void _rejectIncoming() {
    // _callApi.rejectCall();
    Navigator.of(context).pop();
  }

  void _toggleControls() {
    setState(() => _controlsVisible = !_controlsVisible);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF000000),
      body: GestureDetector(
        behavior: HitTestBehavior.opaque,
        onTap: _toggleControls,
        child: SafeArea(
          child: Stack(
            children: [
              // фон / фото / инициалы
              Positioned.fill(
                child: _RemoteBackground(
                  image: widget.remoteImage,
                  initials: _buildInitials(widget.username),
                ),
              ),

              // верхняя панель
              Positioned(
                left: 0,
                right: 0,
                top: 0,
                child: AnimatedOpacity(
                  duration: const Duration(milliseconds: 200),
                  opacity: _controlsVisible ? 1 : 0,
                  child: IgnorePointer(
                    ignoring: !_controlsVisible,
                    child: Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 12,
                        vertical: 8,
                      ),
                      color: Colors.black.withOpacity(0.35),
                      child: Row(
                        children: [
                          IconButton(
                            icon: const Icon(
                              Icons.arrow_back_ios,
                              color: Colors.white,
                            ),
                            onPressed: () => Navigator.of(context).pop(),
                          ),
                          const SizedBox(width: 4),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.center,
                              children: [
                                Text(
                                  widget.username,
                                  maxLines: 1,
                                  overflow: TextOverflow.ellipsis,
                                  style: const TextStyle(
                                    color: Colors.white,
                                    fontSize: 18,
                                    fontWeight: FontWeight.w600,
                                  ),
                                ),
                                const SizedBox(height: 2),
                                Text(
                                  _isRinging
                                      ? 'Входящий звонок'
                                      : _formatDuration(_callDuration),
                                  style: const TextStyle(
                                    color: Colors.white70,
                                    fontSize: 14,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          const SizedBox(width: 40),
                        ],
                      ),
                    ),
                  ),
                ),
              ),

              // локальное превью (можно тоже скрывать в режиме звонка, если нужно)
              Positioned(
                right: 16,
                bottom: 140,
                child: AnimatedOpacity(
                  duration: const Duration(milliseconds: 200),
                  opacity: _controlsVisible && !_isRinging ? 1 : 0,
                  child: IgnorePointer(
                    ignoring: !_controlsVisible || _isRinging,
                    child: _LocalPreview(
                      image: widget.localImage,
                    ),
                  ),
                ),
              ),

              // нижний блок: либо принять/отклонить, либо полный набор
              Positioned(
                left: 0,
                right: 0,
                bottom: 0,
                child: AnimatedOpacity(
                  duration: const Duration(milliseconds: 200),
                  opacity: _controlsVisible ? 1 : 0,
                  child: IgnorePointer(
                    ignoring: !_controlsVisible,
                    child: Padding(
                      padding:
                          const EdgeInsets.fromLTRB(16, 16, 16, 24),
                      child: _isRinging
                          ? _buildIncomingButtons()
                          : _buildInCallButtons(),
                    ),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  /// Режим активного звонка: полный набор контролов.
  Widget _buildInCallButtons() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        MobileCallButton(
          icon: isMuted ? Icons.mic_off : Icons.mic,
          label: 'Mute',
          onTap: _toggleMute,
        ),
        MobileCallButton(
          icon: isVideoOn ? Icons.videocam : Icons.videocam_off,
          label: 'Video',
          onTap: _toggleVideo,
        ),
        MobileCallButton(
          icon:
              isSpeakerOn ? Icons.volume_up : Icons.volume_mute,
          label: 'Speaker',
          onTap: _toggleSpeaker,
        ),
        MobileCallButton(
          icon: Icons.person_add,
          label: 'Добавить',
          backgroundColor: Colors.green,
          onTap: () {
            // TODO: add participant
          },
        ),
        MobileCallButton(
          icon: Icons.call_end,
          label: 'Завершить',
          backgroundColor: Colors.red,
          onTap: _endCall,
        ),
      ],
    );
  }

  /// Режим входящего звонка: только принять/отклонить.
  Widget _buildIncomingButtons() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        MobileCallButton(
          icon: Icons.call_end,
          label: 'Отклонить',
          backgroundColor: Colors.red,
          onTap: _rejectIncoming,
        ),
        MobileCallButton(
          icon: Icons.call,
          label: 'Принять',
          backgroundColor: Colors.green,
          onTap: _acceptIncoming,
        ),
      ],
    );
  }
}

class _RemoteBackground extends StatelessWidget {
  final ImageProvider? image;
  final String initials;

  const _RemoteBackground({
    required this.image,
    required this.initials,
  });

  @override
  Widget build(BuildContext context) {
    if (image != null) {
      return Container(
        decoration: BoxDecoration(
          image: DecorationImage(
            image: image!,
            fit: BoxFit.cover,
          ),
        ),
      );
    }

    return Container(
      decoration: const BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [
            Color(0xFF111121),
            Color(0xFF0D0D1D),
          ],
        ),
      ),
      child: Center(
        child: CircleAvatar(
          radius: 52,
          backgroundColor: const Color(0xFF23283B),
          child: Text(
            initials,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 24,
              fontWeight: FontWeight.w700,
            ),
          ),
        ),
      ),
    );
  }
}

class _LocalPreview extends StatelessWidget {
  final ImageProvider? image;

  const _LocalPreview({this.image});

  @override
  Widget build(BuildContext context) {
    return ClipRRect(
      borderRadius: BorderRadius.circular(12),
      child: SizedBox(
        width: 92,
        height: 140,
        child: image != null
            ? Image(
                image: image!,
                fit: BoxFit.cover,
              )
            : Container(
                color: const Color(0xFF23283B),
                child: const Icon(
                  Icons.person,
                  color: Colors.white70,
                  size: 40,
                ),
              ),
      ),
    );
  }
}

class MobileCallButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final Color backgroundColor;
  final VoidCallback onTap;

  const MobileCallButton({
    super.key,
    required this.icon,
    required this.label,
    required this.onTap,
    this.backgroundColor = const Color(0x22FFFFFF),
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        InkWell(
          borderRadius: BorderRadius.circular(25),
          onTap: onTap,
          child: Container(
            width: 50,
            height: 50,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: backgroundColor,
            ),
            child: Icon(icon, color: Colors.white, size: 22),
          ),
        ),
        const SizedBox(height: 6),
        Text(
          label,
          style: const TextStyle(color: Colors.white70, fontSize: 12),
        ),
      ],
    );
  }
}
