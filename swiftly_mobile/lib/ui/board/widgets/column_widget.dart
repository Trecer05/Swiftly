import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/providers/card_notifier_provider.dart';

import '../../core/themes/colors.dart';
import '../../core/themes/theme.dart';
import 'card_item.dart';
import 'card_item_widget.dart';

class ColumnWidget extends ConsumerWidget {
  final String columnId;
  final String title;
  const ColumnWidget({super.key, required this.columnId, required this.title});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final cards =
        ref
            .watch(cardNotifierProvider)
            .cards
            .where((card) => card.columnId == columnId)
            .toList();
    return IntrinsicWidth(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(title, style: AppTextStyles.style10),
          const SizedBox(height: 10),
          cards.isEmpty
              ? AddCardWidget(onCreate: () => _handleCreate(ref))
              : Column(
                children: [
                  ...cards
                      .expand(
                        (card) => [
                          const SizedBox(height: 10),
                          CardItemWidget(
                            card: card,
                            onDelete: () => _handleDelete(ref, card.id),
                          ),
                        ],
                      )
                      .skip(1),
                ],
              ),
        ],
      ),
    );
  }

  void _handleDelete(WidgetRef ref, String id) {
    ref.read(cardNotifierProvider.notifier).removeCart(id);
  }

  void _handleCreate(WidgetRef ref) {
    ref
        .read(cardNotifierProvider.notifier)
        .addCart(
          CardItem.create(
            title: 'Задача 1',
            description:
                'Lorem ipsum dolor sit amet, consectetur adipiscing elit',
            priority: Priority.high,
            columnId: 'todo',
          ),
        );
  }
}

class AddCardWidget extends StatelessWidget {
  final VoidCallback onCreate;
  const AddCardWidget({super.key, required this.onCreate});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onCreate,
      child: Container(
        width: 300,
        height: 100,
        padding: const EdgeInsets.all(10),
        decoration: BoxDecoration(
          color: AppColors.white15,
          borderRadius: BorderRadius.circular(15),
        ),
        child: const Center(
          child: Icon(Icons.add, color: AppColors.white, size: 18),
        ),
      ),
    );
  }
}
