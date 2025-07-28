import 'package:freezed_annotation/freezed_annotation.dart';

import '../../../domain/models/label_item.dart';

part 'label_state.freezed.dart';

@freezed
class LabelState with _$LabelState {
  const factory LabelState({
    required List<LabelItem> labels,
    @Default(false) bool isLoading,
  }) = _LabelState;

  factory LabelState.initial() => const LabelState(labels: []);
}
