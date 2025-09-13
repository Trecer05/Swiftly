import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/domain/kanban/models/card_item.dart';
import 'package:swiftly_mobile/providers/card_notifier_provider.dart';
import 'package:swiftly_mobile/providers/label_notifier_provider.dart';

import '../../../../domain/kanban/models/priority.dart';
import '../../../../providers/user_notifier_provider.dart';

part 'filter_state.freezed.dart';

@freezed
class FilterState with _$FilterState {
  const factory FilterState({
    String? labelTitle,
    String? priorityTitle,
    String? userName,
    // DateTimeRange? dateRange,
  }) = _FilterState;

  factory FilterState.initial() => const FilterState();
}



class FilterNotifier extends StateNotifier<FilterState> {
  FilterNotifier() : super(FilterState.initial());

  void setLabel(String labelTitle) {
    state = state.copyWith(labelTitle: labelTitle);
  }

  void clearLabel() {
    state = state.copyWith(labelTitle: null);
  }

  void setPriority(String priorityTitle) {
    state = state.copyWith(priorityTitle: priorityTitle);
  }

  void clearPriority() {
    state = state.copyWith(priorityTitle: null);
  }

  void setUser(String userId) {
    state = state.copyWith(userName: userId);
  }

  void clearUser() {
    state = state.copyWith(userName: null);
  }

  void clearAll() {
    state = FilterState.initial();
  }
}

final filterNotifierProvider =
    StateNotifierProvider<FilterNotifier, FilterState>((ref) {
  return FilterNotifier();
});

final filteredCardsProvider = Provider<List<CardItem>>((ref) {
  final cardsState = ref.watch(cardNotifierProvider);
  final filter = ref.watch(filterNotifierProvider);
  final labels = ref.watch(labelNotifierProvider).labels;
  final users = ref.watch(userNotifierProvider).users;

  var cards = cardsState.cards;

  if (filter.labelTitle != null) {
    final filteredCardIds = labels
        .where((l) => l.title == filter.labelTitle)
        .map((l) => l.cardId)
        .whereType<String>()
        .toSet();

    cards = cards.where((c) => filteredCardIds.contains(c.id)).toList();
  }

  if (filter.priorityTitle != null) {
    cards = cards.where((c) => c.priority.title == filter.priorityTitle).toList();
  }

  if (filter.userName != null) {
    final selectedUser = users.firstWhere(
      (u) => '${u.name} ${u.lastName ?? ''}' == filter.userName,
    );

    cards = cards.where((c) => c.userId == selectedUser.id).toList();
    }

  return cards;
});
