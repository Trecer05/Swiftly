import 'package:freezed_annotation/freezed_annotation.dart';

import '../widgets/card_item.dart';

part 'card_state.freezed.dart';

@freezed
class CardState with _$CardState {
  const factory CardState({
    required List<CardItem> cards,
    @Default(false) bool isLoading,
  }) = _CardState;

  factory CardState.initial() => const CardState(cards: []);
}
