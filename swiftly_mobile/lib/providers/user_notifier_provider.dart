import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/ui/settings/state/user_notifier.dart';

import '../ui/settings/state/user_state.dart';

final userNotifierProvider = StateNotifierProvider<UserNotifier, UserState>((ref) {
  return UserNotifier();
});