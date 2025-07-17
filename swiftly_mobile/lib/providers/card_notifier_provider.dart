import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/ui/board/view_models/card_notifier.dart';

import '../ui/board/widgets/card_state.dart';

final cardNotifierProvider = StateNotifierProvider<CardNotifier, CardState>((ref) {
  return CardNotifier();
});