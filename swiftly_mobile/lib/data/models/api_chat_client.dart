import 'dart:async';
import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:web_socket_channel/web_socket_channel.dart';

class ChatApiConfig {
  static const String baseUrl = 'https://api-dev.example.com'; // TODO: из env
  static const String wsHost = 'wss://ws-dev.example.com';     // TODO: из env

  // TODO: подменить реальным токеном после авторизации
  static String accessToken = 'DEBUG_TOKEN';
}

class ChatApiClient {
  final http.Client _client;

  ChatApiClient(this._client);

  Map<String, String> _headers() => {
        'Authorization': 'Bearer ${ChatApiConfig.accessToken}',
        'Content-Type': 'application/json',
      };

  Future<List<MessageDto>> getMessages(String chatId) async {
    final uri = Uri.parse(
      '${ChatApiConfig.baseUrl}/api/v1/chat/$chatId/messages?limit=50&offset=0',
    );
    final res = await _client.get(uri, headers: _headers());

    if (res.statusCode != 200) {
      throw Exception('Failed to load messages: ${res.statusCode}');
    }

    final data = jsonDecode(res.body) as List<dynamic>;
    return data.map((e) => MessageDto.fromJson(e as Map<String, dynamic>)).toList();
  }
}

class ChatWsClient {
  WebSocketChannel? _channel;

  Stream<dynamic>? get stream => _channel?.stream;

  Future<void> connect(String chatId) async {
    final uri = Uri.parse(
      '${ChatApiConfig.wsHost}/api/v1/chat/$chatId?limit=50&offset=0',
    );
    _channel = WebSocketChannel.connect(uri);
    // если бек требует токен в query/header — добавишь здесь по доке
  }

  void sendText(String text) {
    final data = {
      'type': 'message',
      'text': text,
    };
    _channel?.sink.add(jsonEncode(data));
  }

  void sendTyping(bool typing) {
    _channel?.sink.add(jsonEncode({
      'type': typing ? 'typing' : 'stoptyping',
    }));
  }

  void dispose() {
    _channel?.sink.close();
  }
}

// пример DTO, чтобы все компилилось; потом подгони под свою схему
class MessageDto {
  final String id;
  final String text;
  final DateTime createdAt;
  final String fromId;

  MessageDto({
    required this.id,
    required this.text,
    required this.createdAt,
    required this.fromId,
  });

  factory MessageDto.fromJson(Map<String, dynamic> json) {
    return MessageDto(
      id: json['id'].toString(),
      text: json['text'] as String? ?? '',
      createdAt: DateTime.parse(json['created_at'] as String),
      fromId: json['from_id'].toString(),
    );
  }
}
