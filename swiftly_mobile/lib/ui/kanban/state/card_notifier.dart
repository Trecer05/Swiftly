import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/domain/models/label_item.dart';

import '../../../domain/kanban/models/card_item.dart';
import '../../../domain/kanban/models/priority.dart';
import '../../core/themes/colors.dart';
import 'card_state.dart';

class CardNotifier extends StateNotifier<CardState> {
  CardNotifier() : super(CardState.initial());

  void addCard(CardItem card) {
    state = state.copyWith(cards: [...state.cards, card]);
  }

  void removeCard(String id) {
    state = state.copyWith(
      cards: state.cards.where((card) => card.id != id).toList(),
    );
  }

  Future<void> loadCards() async {
    state = state.copyWith(isLoading: true);
    await Future.delayed(const Duration(seconds: 1));
    state = state.copyWith(
      cards: mockCards,
      isLoading: false,
    );
  }
}

final mockCards = [
  CardItem.create(userId: 'aaa', title: 'Задача 1', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', priority: Priority.high, columnId: 'todo'),
  CardItem.create(userId: 'bbb', columnId: 'progress'),
  CardItem.create(userId: 'bbb', title: 'Задача 2', description: 'Сделать что-то', priority: Priority.medium, category: LabelItem(title: 'программирование', color: AppColors.amaranthMagenta), columnId: 'done'),
  CardItem.create(userId: 'ddd', title: 'Задача 2', description: 'Сделать что-то', priority: Priority.medium, category: LabelItem(title: 'программирование', color: AppColors.amaranthMagenta), columnId: 'progress'),
  CardItem.create(userId: 'eee', title: 'Задача 1', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', priority: Priority.high, columnId: 'progress'),
];