import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../domain/kanban/models/card_item.dart';
import '../../../domain/kanban/models/priority.dart';
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

  void updateTitle(String cardId, String newTitle) {
    final processedTitle = newTitle.trim().isEmpty ? 'Новая задача' : newTitle;
    state = state.copyWith(
      cards:
          state.cards.map((card) {
            return card.id == cardId
                ? card.copyWith(title: processedTitle)
                : card;
          }).toList(),
    );
  }

  void updateDescription(String cardId, String newDescription) {
    final processedDescription =
        newDescription.trim().isEmpty ? 'Пустое описание' : newDescription;
    state = state.copyWith(
      cards:
          state.cards.map((card) {
            return card.id == cardId
                ? card.copyWith(description: processedDescription)
                : card;
          }).toList(),
    );
  }

  void updateColumn(String cardId, String newColumn) {
    state = state.copyWith(
      cards:
          state.cards.map((card) {
            return card.id == cardId
                ? card.copyWith(columnId: newColumn)
                : card;
          }).toList(),
    );
  }

  void updatePriority(String cardId, Priority newPriority) {
    state = state.copyWith(
      cards:
          state.cards.map((card) {
            return card.id == cardId
                ? card.copyWith(priority: newPriority)
                : card;
          }).toList(),
    );
  }

  Future<void> loadCards(WidgetRef ref) async {
    state = state.copyWith(isLoading: true);
    await Future.delayed(const Duration(seconds: 1));

    final cards = [
      CardItem.create(
        userId: '1',
        title: 'Задача 1',
        description: 'Lorem ipsum',
        priority: Priority.high,
        columnId: 'todo',
      ),
      CardItem.create(
        userId: '2',
        priority: Priority.low,
        columnId: 'progress',
      ),
      CardItem.create(
        userId: 'aaa',
        title: 'Задача 2',
        description: 'Сделать что-то',
        priority: Priority.medium,
        columnId: 'progress',
      ),
      CardItem.create(
        userId: '3',
        title: 'Задача 1',
        description: 'Lorem ipsum',
        priority: Priority.high,
        columnId: 'progress',
      ),
    ];

    state = state.copyWith(cards: cards, isLoading: false);
  }
}
