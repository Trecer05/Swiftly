import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../ui/core/state/label_notifier.dart';
import '../ui/core/state/label_state.dart';


final labelNotifierProvider = StateNotifierProvider<LabelNotifier, LabelState>((ref) {
  return LabelNotifier();
});