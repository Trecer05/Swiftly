import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/ui/kanban/state/card_notifier.dart';

import '../ui/kanban/state/card_state.dart';

final cardNotifierProvider = StateNotifierProvider<CardNotifier, CardState>((ref) {
  return CardNotifier();
});