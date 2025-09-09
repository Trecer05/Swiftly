import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/providers/card_notifier_provider.dart';

import '../../../../providers/selected_column_status_provider.dart';
import '../../../core/themes/colors.dart';
import '../../../core/themes/theme.dart';
import '../../../core/ui/card_item/card_item_mobile.dart';

class ColumnWidgetMobile extends ConsumerWidget {
  const ColumnWidgetMobile({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final selectedColumnStatus = ref.watch(selectedColumnStatusProvider);
    final cards =
        ref
            .watch(cardNotifierProvider)
            .cards
            .where((card) => card.columnId == selectedColumnStatus)
            .toList();
    return IntrinsicWidth(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          cards.isEmpty
              ? const Center(child: Text('Пусто'),)
              : Column(
                children: [
                  ...cards
                      .expand(
                        (card) => [
                          const SizedBox(height: 10),
                          CardItemMobile(card: card),
                        ],
                      )
                      .skip(1),
                ],
              ),
        ],
      ),
    );
  }

  // void _handleCreateCard(BuildContext context, WidgetRef ref, String columnId) {
  //   final currentUser = ref.watch(currentUserProvider);
  //   if (currentUser == null) return;
  //   final newCard = CardItem.create(
  //     userId: currentUser.id,
  //     priority: Priority.low,
  //     columnId: columnId,
  //   );
  //   ref.read(cardNotifierProvider.notifier).addCard(newCard);
  //   showDialog(
  //     context: context,
  //     builder:
  //         (context) => Dialog(
  //           insetPadding: const EdgeInsets.all(16),
  //           backgroundColor: Colors.transparent,
  //           child: CartDetailsWidget(card: newCard),
  //         ),
  //   );
  // }
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

class ColumnStatusButton extends ConsumerWidget {
  final String columnId;
  final String title;
  const ColumnStatusButton({
    super.key,
    required this.columnId,
    required this.title,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final selectedColumnStatus = ref.watch(selectedColumnStatusProvider);
    return GestureDetector(
      onTap: () {ref.read(selectedColumnStatusProvider.notifier).state = columnId;},
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
        decoration: BoxDecoration(
          border: Border(bottom: BorderSide(color: selectedColumnStatus == columnId ? AppColors.white : AppColors.transparent)),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Text(
          title,
          style: 
            selectedColumnStatus == columnId ? AppTextStyles.style10 : AppTextStyles.style18,
        ),
      ),
    );
  }
}
