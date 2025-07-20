import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../domain/kanban/models/card_item.dart';
import '../../../domain/kanban/models/priority.dart';
import 'card_state.dart';

class CardNotifier extends StateNotifier<CardState> {
  CardNotifier() : super(CardState.initial());

  void addCart(CardItem card) {
    state = state.copyWith(cards: [...state.cards, card]);
  }

  void removeCart(String id) {
    state = state.copyWith(
      cards: state.cards.where((card) => card.id != id).toList(),
    );
  }

  Future<void> loadCarts() async {
    state = state.copyWith(isLoading: true);
    await Future.delayed(const Duration(seconds: 1));
    state = state.copyWith(
      cards: mockCards,
      isLoading: false,
    );
  }
}

final mockCards = [
  CardItem.create(title: 'Задача 1', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', priority: Priority.high, columnId: 'todo'),
  CardItem.create(columnId: 'progress'),
  CardItem.create(title: 'Задача 2', description: 'Сделать что-то', priority: Priority.medium, category: 'программирование', columnId: 'done'),
  CardItem.create(title: 'Задача 2', description: 'Сделать что-то', priority: Priority.medium, category: 'программирование', columnId: 'progress'),
  CardItem.create(title: 'Задача 1', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', priority: Priority.high, columnId: 'progress'),
];